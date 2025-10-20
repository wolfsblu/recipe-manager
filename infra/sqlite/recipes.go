package sqlite

import (
	"context"
	"net/url"
	"sort"
	"time"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
	_ "modernc.org/sqlite"
)

type recipeRelations struct {
	tags        []database.GetTagsForRecipesRow
	images      []database.GetImagesForRecipesRow
	steps       []database.GetStepsForRecipesRow
	ingredients []database.GetIngredientsForRecipesRow
	nutrients   []database.GetNutrientsForRecipesRow
}

func (s *Store) CreateRecipe(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	var recipeId int64
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		var err error
		recipeId, err = tx.query().CreateRecipe(ctx, tx.mapper.FromRecipe(recipe))
		if err != nil {
			return err
		}
		if err = tx.createRecipeSteps(ctx, recipeId, recipe.Steps); err != nil {
			return err
		}
		if err = tx.createRecipeImages(ctx, recipeId, recipe.Images); err != nil {
			return err
		}
		if err = tx.createRecipeTags(ctx, recipeId, recipe.Tags); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.Recipe{}, err
	}
	return s.GetRecipeById(ctx, recipe.CreatedBy, recipeId)
}

func (s *Store) createRecipeSteps(ctx context.Context, recipeID int64, steps []domain.RecipeStep) error {
	for i, step := range steps {
		stepId, err := s.query().CreateRecipeStep(ctx, s.mapper.FromRecipeStep(recipeID, step, int64(i)))
		if err != nil {
			return err
		}

		if err = s.createStepIngredients(ctx, stepId, step.Ingredients); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createStepIngredients(ctx context.Context, stepID int64, ingredients []domain.StepIngredient) error {
	for i, ingredient := range ingredients {
		_, err := s.query().CreateStepIngredient(ctx, s.mapper.FromStepIngredient(stepID, ingredient, int64(i)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createRecipeImages(ctx context.Context, recipeID int64, images []domain.RecipeImage) error {
	for _, image := range images {
		_, err := s.query().CreateRecipeImages(ctx, s.mapper.FromRecipeImage(recipeID, image))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createRecipeTags(ctx context.Context, recipeID int64, tags []domain.Tag) error {
	for _, tag := range tags {
		err := s.query().CreateRecipeTag(ctx, s.mapper.FromRecipeTag(recipeID, tag))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) DeleteRecipe(ctx context.Context, id int64) error {
	return s.query().DeleteRecipe(ctx, id)
}

func (s *Store) CreateMealPlan(ctx context.Context, entry domain.MealPlanEntry) error {
	return s.query().CreateMealPlan(ctx, s.mapper.FromMealPlanEntry(entry))
}

func (s *Store) DeleteMealPlan(ctx context.Context, userID int64, recipeID int64, date time.Time) error {
	return s.query().DeleteMealPlan(ctx, database.DeleteMealPlanParams{
		UserID:   userID,
		RecipeID: recipeID,
		Date:     date.Format(time.DateOnly),
	})
}

func (s *Store) GetMealPlan(ctx context.Context, user *domain.User, from time.Time, until time.Time, page domain.Page) (domain.Result[domain.MealPlan], error) {
	cursor, err := domain.DecodeCursor[*domain.DateCursor](page.Cursor)
	if err != nil {
		cursor = &domain.DateCursor{}
	}
	lastDate := ""
	if !cursor.LastDate.IsZero() {
		lastDate = cursor.LastDate.Format(time.DateOnly)
	}

	result, err := s.query().GetMealPlan(ctx, database.GetMealPlanParams{
		UserID:    user.ID,
		FromDate:  from.Format(time.DateOnly),
		UntilDate: until.Format(time.DateOnly),
		LastDate:  lastDate,
		LastID:    cursor.LastID,
		Limit:     int64(page.Limit + 1),
	})
	if err != nil {
		return domain.Result[domain.MealPlan]{}, err
	}

	grouped := s.groupMealPlansByDate(result)
	uniqueRecipes := s.extractUniqueRecipes(grouped)
	populatedMap, err := s.populateRecipeMap(ctx, user, uniqueRecipes)
	if err != nil {
		return domain.Result[domain.MealPlan]{}, err
	}

	mealPlans := s.buildMealPlans(grouped, populatedMap)

	return domain.NewPagedResult(mealPlans, page.Limit, func(m domain.MealPlan) domain.DateCursor {
		lastEntry := m.Entries[len(m.Entries)-1]
		return domain.DateCursor{
			LastID:   lastEntry.ID,
			LastDate: m.Date,
		}
	}), nil
}

func (s *Store) groupMealPlansByDate(result []database.GetMealPlanRow) map[string]*domain.MealPlan {
	grouped := make(map[string]*domain.MealPlan)

	for _, item := range result {
		dateKey := item.MealPlan.Date

		if _, exists := grouped[dateKey]; !exists {
			date, _ := time.Parse(time.DateOnly, dateKey)
			grouped[dateKey] = &domain.MealPlan{
				Date:    date,
				Entries: make([]domain.MealPlanRecipe, 0),
			}
		}

		mealPlan := grouped[dateKey]
		recipe := s.mapper.ToRecipe(item.Recipe)
		mealPlan.Entries = append(mealPlan.Entries, domain.MealPlanRecipe{
			ID:     item.MealPlan.ID,
			Recipe: recipe,
		})
	}

	return grouped
}

func (s *Store) extractUniqueRecipes(grouped map[string]*domain.MealPlan) []domain.Recipe {
	uniqueRecipes := make([]domain.Recipe, 0)
	seen := make(map[int64]bool)

	for _, mealPlan := range grouped {
		for _, entry := range mealPlan.Entries {
			if !seen[entry.Recipe.ID] {
				uniqueRecipes = append(uniqueRecipes, entry.Recipe)
				seen[entry.Recipe.ID] = true
			}
		}
	}

	return uniqueRecipes
}

func (s *Store) populateRecipeMap(ctx context.Context, user *domain.User, recipes []domain.Recipe) (map[int64]domain.Recipe, error) {
	populatedRecipes, err := s.populateRecipeRelations(ctx, user, recipes)
	if err != nil {
		return nil, err
	}

	recipeMap := make(map[int64]domain.Recipe, len(populatedRecipes))
	for _, recipe := range populatedRecipes {
		recipeMap[recipe.ID] = recipe
	}

	return recipeMap, nil
}

func (s *Store) buildMealPlans(grouped map[string]*domain.MealPlan, populatedMap map[int64]domain.Recipe) []domain.MealPlan {
	mealPlans := make([]domain.MealPlan, 0, len(grouped))

	for _, mealPlan := range grouped {
		for i, entry := range mealPlan.Entries {
			mealPlan.Entries[i].Recipe = populatedMap[entry.Recipe.ID]
		}
		mealPlans = append(mealPlans, *mealPlan)
	}

	sort.Slice(mealPlans, func(i, j int) bool {
		return mealPlans[i].Date.Before(mealPlans[j].Date)
	})

	return mealPlans
}

func (s *Store) GetTags(ctx context.Context, req domain.Page) (domain.Result[domain.Tag], error) {
	cursor, err := domain.DecodeCursor[*domain.NameCursor](req.Cursor)
	if err != nil {
		cursor = &domain.NameCursor{}
	}
	result, err := s.query().GetTags(ctx, database.GetTagsParams{
		LastName: cursor.LastName,
		LastID:   cursor.LastID,
		Limit:    int64(req.Limit + 1),
	})
	if err != nil {
		return domain.Result[domain.Tag]{}, err
	}

	tags := make([]domain.Tag, 0, len(result))
	for _, tag := range result {
		tags = append(tags, s.mapper.ToTag(tag))
	}

	return domain.NewPagedResult(tags, req.Limit, func(t domain.Tag) domain.NameCursor {
		return domain.NameCursor{
			LastID:   t.ID,
			LastName: t.Name,
		}
	}), nil
}

func (s *Store) GetRecipeById(ctx context.Context, user *domain.User, id int64) (recipe domain.Recipe, _ error) {
	result, err := s.query().GetRecipe(ctx, id)
	if err != nil {
		return recipe, err
	}

	recipe = s.mapper.ToRecipe(result)
	populatedRecipes, err := s.populateRecipeRelations(ctx, user, []domain.Recipe{recipe})
	if err != nil {
		return recipe, err
	}

	return populatedRecipes[0], nil
}

func (s *Store) GetRecipesByUser(ctx context.Context, user *domain.User, req domain.Page) (domain.Result[domain.Recipe], error) {
	cursor, err := domain.DecodeCursor[*domain.DateCursor](req.Cursor)
	if err != nil {
		cursor = &domain.DateCursor{}
	}

	result, err := s.query().ListRecipes(ctx, database.ListRecipesParams{
		CreatedBy:     user.ID,
		LastCreatedAt: cursor.LastDate,
		LastID:        cursor.LastID,
		Limit:         int64(req.Limit + 1),
	})
	if err != nil {
		return domain.Result[domain.Recipe]{}, err
	}

	recipes := make([]domain.Recipe, 0, len(result))
	for _, recipe := range result {
		recipes = append(recipes, s.mapper.ToRecipe(recipe))
	}

	// Populate relations before creating paginated result
	populatedRecipes, err := s.populateRecipeRelations(ctx, user, recipes)
	if err != nil {
		return domain.Result[domain.Recipe]{}, err
	}

	return domain.NewPagedResult(populatedRecipes, req.Limit, func(r domain.Recipe) domain.DateCursor {
		return domain.DateCursor{
			LastID:   r.ID,
			LastDate: r.CreatedAt,
		}
	}), nil
}

func (s *Store) populateRecipeRelations(ctx context.Context, user *domain.User, recipes []domain.Recipe) ([]domain.Recipe, error) {
	if len(recipes) == 0 {
		return recipes, nil
	}

	recipeIds := make([]int64, len(recipes))
	for i, recipe := range recipes {
		recipeIds[i] = recipe.ID
	}

	relations, err := s.getRecipeRelations(ctx, user, recipeIds)
	if err != nil {
		return nil, err
	}

	tagsByRecipe := s.groupTagsByRecipe(relations.tags)
	ingredientMap := s.groupIngredientsWithNutrients(relations.ingredients, relations.nutrients)
	ingredientsByStep := s.groupIngredientsByStep(relations.ingredients, ingredientMap)
	stepsByRecipe := s.groupStepsByRecipe(relations.steps, ingredientsByStep)
	imagesByRecipe, err := s.groupImagesByRecipe(relations.images)
	if err != nil {
		return nil, err
	}

	populatedRecipes := make([]domain.Recipe, len(recipes))
	for i, recipe := range recipes {
		recipe.Tags = tagsByRecipe[recipe.ID]
		recipe.Images = imagesByRecipe[recipe.ID]
		recipe.Steps = stepsByRecipe[recipe.ID]
		populatedRecipes[i] = recipe
	}
	return populatedRecipes, nil
}

func (s *Store) getRecipeRelations(ctx context.Context, user *domain.User, recipeIds []int64) (*recipeRelations, error) {
	tags, err := s.query().GetTagsForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}

	images, err := s.query().GetImagesForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}

	steps, err := s.query().GetStepsForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}

	ingredients, err := s.query().GetIngredientsForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}

	nutrients, err := s.query().GetNutrientsForRecipes(ctx, recipeIds)
	if err != nil {
		return nil, err
	}

	return &recipeRelations{
		tags:        tags,
		images:      images,
		steps:       steps,
		ingredients: ingredients,
		nutrients:   nutrients,
	}, nil
}

func (s *Store) groupTagsByRecipe(tags []database.GetTagsForRecipesRow) map[int64][]domain.Tag {
	tagsByRecipe := make(map[int64][]domain.Tag)
	for _, tag := range tags {
		tagsByRecipe[tag.RecipeID] = append(tagsByRecipe[tag.RecipeID], s.mapper.ToTag(tag.Tag))
	}
	return tagsByRecipe
}

func (s *Store) groupImagesByRecipe(images []database.GetImagesForRecipesRow) (map[int64][]domain.RecipeImage, error) {
	imagesByRecipe := make(map[int64][]domain.RecipeImage)
	for _, image := range images {
		imageUrl, err := url.ParseRequestURI(image.Path)
		if err != nil {
			return nil, err
		}
		imagesByRecipe[image.RecipeID] = append(imagesByRecipe[image.RecipeID], domain.RecipeImage{
			ID:  image.ID,
			URL: imageUrl,
		})
	}
	return imagesByRecipe, nil
}

func (s *Store) groupIngredientsWithNutrients(ingredients []database.GetIngredientsForRecipesRow, nutrients []database.GetNutrientsForRecipesRow) map[int64]domain.Ingredient {
	nutrientsByIngredient := make(map[int64][]domain.IngredientNutrient)
	for _, nutrient := range nutrients {
		nutrientsByIngredient[nutrient.IngredientID] = append(
			nutrientsByIngredient[nutrient.IngredientID],
			domain.IngredientNutrient{
				Nutrient: s.mapper.ToNutrient(nutrient.Nutrient),
				Amount:   nutrient.Amount,
			},
		)
	}

	ingredientMap := make(map[int64]domain.Ingredient)
	for _, ing := range ingredients {
		if _, exists := ingredientMap[ing.IngredientID]; !exists {
			ingredient := s.mapper.ToIngredientFromRecipeRow(ing)
			ingredient.Nutrients = nutrientsByIngredient[ing.IngredientID]
			if ingredient.Nutrients == nil {
				ingredient.Nutrients = []domain.IngredientNutrient{}
			}
			ingredientMap[ing.IngredientID] = ingredient
		}
	}
	return ingredientMap
}

func (s *Store) groupIngredientsByStep(ingredients []database.GetIngredientsForRecipesRow, ingredientMap map[int64]domain.Ingredient) map[int64][]domain.StepIngredient {
	ingredientsByStep := make(map[int64][]domain.StepIngredient)
	for _, ing := range ingredients {
		stepIngredient := s.mapper.ToStepIngredient(ing)
		stepIngredient.Ingredient = ingredientMap[ing.IngredientID]
		ingredientsByStep[ing.StepID] = append(ingredientsByStep[ing.StepID], stepIngredient)
	}
	return ingredientsByStep
}

func (s *Store) groupStepsByRecipe(steps []database.GetStepsForRecipesRow, ingredientsByStep map[int64][]domain.StepIngredient) map[int64][]domain.RecipeStep {
	stepsByRecipe := make(map[int64][]database.GetStepsForRecipesRow)
	for _, step := range steps {
		stepsByRecipe[step.RecipeID] = append(stepsByRecipe[step.RecipeID], step)
	}

	result := make(map[int64][]domain.RecipeStep)
	for recipeID, recipeSteps := range stepsByRecipe {
		stepsForRecipe := make([]domain.RecipeStep, len(recipeSteps))
		for i, step := range recipeSteps {
			recipeStep := s.mapper.ToRecipeStep(step)
			recipeStep.Ingredients = ingredientsByStep[step.ID]
			stepsForRecipe[i] = recipeStep
		}
		result[recipeID] = stepsForRecipe
	}
	return result
}

func (s *Store) UpdateRecipe(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		err := tx.query().UpdateRecipe(ctx, tx.mapper.FromRecipeForUpdate(recipe))
		if err != nil {
			return err
		}

		if err = tx.query().DeleteRecipeIngredients(ctx, recipe.ID); err != nil {
			return err
		}
		if err = tx.query().DeleteRecipeSteps(ctx, recipe.ID); err != nil {
			return err
		}
		if err = tx.query().DeleteRecipeImages(ctx, recipe.ID); err != nil {
			return err
		}
		if err = tx.query().DeleteRecipeTags(ctx, recipe.ID); err != nil {
			return err
		}

		if err = tx.createRecipeSteps(ctx, recipe.ID, recipe.Steps); err != nil {
			return err
		}
		if err = tx.createRecipeImages(ctx, recipe.ID, recipe.Images); err != nil {
			return err
		}
		if err = tx.createRecipeTags(ctx, recipe.ID, recipe.Tags); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Recipe{}, err
	}
	return s.GetRecipeById(ctx, recipe.CreatedBy, recipe.ID)
}
