package sqlite

import (
	"context"
	"net/url"
	"sort"
	"time"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/gen/model"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/queries"
	_ "modernc.org/sqlite"
)

type recipeRelations struct {
	tags        []recipeTag
	images      []recipeImage
	steps       []recipeStep
	ingredients []recipeIngredient
	nutrients   []ingredientNutrient
}

type recipeTag struct {
	RecipeID int64
	Tag      model.Tag
}

type recipeImage struct {
	model.RecipeImage
}

type recipeStep struct {
	model.RecipeStep
}

type recipeIngredient struct {
	RecipeIngredient model.RecipeIngredient
	Ingredient       model.Ingredient
	Unit             model.Unit
	StepID           int64
}

type ingredientNutrient struct {
	IngredientNutrient model.IngredientNutrient
	Nutrient           model.Nutrient
}

func (s *Store) CreateRecipe(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	var recipeId int64
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		var result model.Recipe
		err := queries.InsertRecipe(recipe.Name, recipe.Description, recipe.Servings, recipe.Minutes, recipe.CreatedBy.ID).Query(tx.DB(), &result)
		if err != nil {
			return err
		}
		recipeId = result.ID

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
		var result model.RecipeStep
		err := queries.InsertRecipeStep(recipeID, step.Instructions, int64(i)).Query(s.DB(), &result)
		if err != nil {
			return err
		}

		if err = s.createStepIngredients(ctx, result.ID, step.Ingredients); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createStepIngredients(ctx context.Context, stepID int64, ingredients []domain.StepIngredient) error {
	for i, ingredient := range ingredients {
		_, err := queries.InsertStepIngredient(stepID, ingredient.Ingredient.ID, ingredient.Unit.ID, ingredient.Amount, int64(i)).Exec(s.DB())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createRecipeImages(ctx context.Context, recipeID int64, images []domain.RecipeImage) error {
	for i, image := range images {
		_, err := queries.InsertRecipeImage(recipeID, image.URL.String(), int64(i)).Exec(s.DB())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createRecipeTags(ctx context.Context, recipeID int64, tags []domain.Tag) error {
	for _, tag := range tags {
		_, err := queries.InsertRecipeTag(recipeID, tag.ID).Exec(s.DB())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) DeleteRecipe(ctx context.Context, id int64) error {
	_, err := queries.DeleteRecipe(id).Exec(s.DB())
	return err
}

func (s *Store) CreateMealPlan(ctx context.Context, entry domain.MealPlanEntry) error {
	_, err := queries.InsertMealPlan(entry.Date.Format(time.DateOnly), entry.UserID, entry.RecipeID, entry.SortOrder).Exec(s.DB())
	return err
}

func (s *Store) DeleteMealPlan(ctx context.Context, userID int64, recipeID int64, date time.Time) error {
	_, err := queries.DeleteMealPlan(userID, recipeID, date.Format(time.DateOnly)).Exec(s.DB())
	return err
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

	type ResultRow struct {
		MealPlan model.MealPlan
		Recipe   model.Recipe
	}
	var results []ResultRow
	err = queries.SelectMealPlan(user.ID, from.Format(time.DateOnly), until.Format(time.DateOnly), lastDate, cursor.LastID, int64(page.Limit+1)).Query(s.DB(), &results)
	if err != nil {
		return domain.Result[domain.MealPlan]{}, err
	}

	// Convert to domain model
	grouped := make(map[string]*domain.MealPlan)
	for _, item := range results {
		dateKey := item.MealPlan.Date

		if _, exists := grouped[dateKey]; !exists {
			date, _ := time.Parse(time.DateOnly, dateKey)
			grouped[dateKey] = &domain.MealPlan{
				Date:    date,
				Entries: make([]domain.MealPlanRecipe, 0),
			}
		}

		mealPlan := grouped[dateKey]
		recipe := s.mapper.ToRecipeFromModel(item.Recipe)
		mealPlan.Entries = append(mealPlan.Entries, domain.MealPlanRecipe{
			ID:     item.MealPlan.ID,
			Recipe: recipe,
		})
	}

	uniqueRecipes := s.extractUniqueRecipes(grouped)
	populatedMap, err := s.populateRecipeMap(ctx, uniqueRecipes)
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

func (s *Store) groupMealPlansByDate(grouped map[string]*domain.MealPlan) map[string]*domain.MealPlan {
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

func (s *Store) populateRecipeMap(ctx context.Context, recipes []domain.Recipe) (map[int64]domain.Recipe, error) {
	populatedRecipes, err := s.populateRecipeRelations(ctx, recipes)
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

	var result []model.Tag
	err = queries.SelectTags(cursor.LastID, cursor.LastName, int64(req.Limit+1)).Query(s.DB(), &result)
	if err != nil {
		return domain.Result[domain.Tag]{}, err
	}

	tags := make([]domain.Tag, 0, len(result))
	for _, tag := range result {
		tags = append(tags, s.mapper.ToTagFromModel(tag))
	}

	return domain.NewPagedResult(tags, req.Limit, func(t domain.Tag) domain.NameCursor {
		return domain.NameCursor{
			LastID:   t.ID,
			LastName: t.Name,
		}
	}), nil
}

func (s *Store) GetRecipeById(ctx context.Context, user *domain.User, id int64) (recipe domain.Recipe, _ error) {
	var result model.Recipe
	err := queries.SelectRecipeByID(id).Query(s.DB(), &result)
	if err != nil {
		return recipe, err
	}

	recipe = s.mapper.ToRecipeFromModel(result)
	populatedRecipes, err := s.populateRecipeRelations(ctx, []domain.Recipe{recipe})
	if err != nil {
		return recipe, err
	}

	return populatedRecipes[0], nil
}

func (s *Store) GetRecipesByUser(ctx context.Context, user *domain.User, req domain.Page, sortConfig domain.RecipeSort) (domain.Result[domain.Recipe], error) {
	cursor, err := domain.DecodeCursor[*domain.RecipeCursor](req.Cursor)
	if err != nil {
		cursor = s.getDefaultCursor(sortConfig)
	}

	result, err := s.listRecipesBySortOrder(ctx, user.ID, cursor, int64(req.Limit+1), sortConfig)
	if err != nil {
		return domain.Result[domain.Recipe]{}, err
	}

	recipes := make([]domain.Recipe, 0, len(result))
	for _, recipe := range result {
		recipes = append(recipes, s.mapper.ToRecipeFromModel(recipe))
	}

	populatedRecipes, err := s.populateRecipeRelations(ctx, recipes)
	if err != nil {
		return domain.Result[domain.Recipe]{}, err
	}

	return domain.NewPagedResult(populatedRecipes, req.Limit, func(r domain.Recipe) domain.RecipeCursor {
		return domain.RecipeCursor{
			LastID:        r.ID,
			LastName:      r.Name,
			LastCreatedAt: r.CreatedAt,
			LastServings:  r.Servings,
		}
	}), nil
}

func (s *Store) getDefaultCursor(sortConfig domain.RecipeSort) *domain.RecipeCursor {
	if sortConfig.Order == domain.SortOrderAsc {
		return &domain.RecipeCursor{}
	}

	switch sortConfig.Field {
	case domain.SortFieldName:
		return domain.NewDescendingRecipeCursorByName()
	case domain.SortFieldServings:
		return domain.NewDescendingRecipeCursorByServings()
	case domain.SortFieldCreatedAt:
		return domain.NewDescendingRecipeCursorByCreatedAt()
	default:
		return domain.NewDescendingRecipeCursorByCreatedAt()
	}
}

func (s *Store) listRecipesBySortOrder(ctx context.Context, userID int64, cursor *domain.RecipeCursor, limit int64, sortConfig domain.RecipeSort) ([]model.Recipe, error) {
	recipeCursor := queries.RecipeCursor{
		LastID:        cursor.LastID,
		LastName:      cursor.LastName,
		LastCreatedAt: cursor.LastCreatedAt,
		LastServings:  cursor.LastServings,
	}

	var results []model.Recipe
	err := queries.SelectRecipesByUserPaginated(userID, recipeCursor, limit, string(sortConfig.Field), string(sortConfig.Order)).Query(s.DB(), &results)
	return results, err
}

func (s *Store) populateRecipeRelations(ctx context.Context, recipes []domain.Recipe) ([]domain.Recipe, error) {
	if len(recipes) == 0 {
		return recipes, nil
	}

	recipeIds := make([]int64, len(recipes))
	for i, recipe := range recipes {
		recipeIds[i] = recipe.ID
	}

	relations, err := s.getRecipeRelations(ctx, recipeIds)
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

func (s *Store) getRecipeRelations(ctx context.Context, recipeIds []int64) (*recipeRelations, error) {
	var tags []recipeTag
	err := queries.SelectTagsForRecipes(recipeIds).Query(s.DB(), &tags)
	if err != nil {
		return nil, err
	}

	var images []recipeImage
	err = queries.SelectImagesForRecipes(recipeIds).Query(s.DB(), &images)
	if err != nil {
		return nil, err
	}

	var steps []recipeStep
	err = queries.SelectStepsForRecipes(recipeIds).Query(s.DB(), &steps)
	if err != nil {
		return nil, err
	}

	var ingredients []recipeIngredient
	err = queries.SelectIngredientsForRecipes(recipeIds).Query(s.DB(), &ingredients)
	if err != nil {
		return nil, err
	}

	var nutrients []ingredientNutrient
	err = queries.SelectNutrientsForRecipes(recipeIds).Query(s.DB(), &nutrients)
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

func (s *Store) groupTagsByRecipe(tags []recipeTag) map[int64][]domain.Tag {
	tagsByRecipe := make(map[int64][]domain.Tag)
	for _, tag := range tags {
		tagsByRecipe[tag.RecipeID] = append(tagsByRecipe[tag.RecipeID], s.mapper.ToTagFromModel(tag.Tag))
	}
	return tagsByRecipe
}

func (s *Store) groupImagesByRecipe(images []recipeImage) (map[int64][]domain.RecipeImage, error) {
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

func (s *Store) groupIngredientsWithNutrients(ingredients []recipeIngredient, nutrients []ingredientNutrient) map[int64]domain.Ingredient {
	nutrientsByIngredient := make(map[int64][]domain.IngredientNutrient)
	for _, nutrient := range nutrients {
		nutrientsByIngredient[nutrient.IngredientNutrient.IngredientID] = append(
			nutrientsByIngredient[nutrient.IngredientNutrient.IngredientID],
			domain.IngredientNutrient{
				Nutrient: s.mapper.ToNutrientFromModel(nutrient.Nutrient),
				Amount:   nutrient.IngredientNutrient.Amount,
			},
		)
	}

	ingredientMap := make(map[int64]domain.Ingredient)
	for _, ing := range ingredients {
		if _, exists := ingredientMap[ing.Ingredient.ID]; !exists {
			ingredient := domain.Ingredient{
				ID:   ing.Ingredient.ID,
				Name: ing.Ingredient.Name,
			}
			ingredient.Nutrients = nutrientsByIngredient[ing.Ingredient.ID]
			if ingredient.Nutrients == nil {
				ingredient.Nutrients = []domain.IngredientNutrient{}
			}
			ingredientMap[ing.Ingredient.ID] = ingredient
		}
	}
	return ingredientMap
}

func (s *Store) groupIngredientsByStep(ingredients []recipeIngredient, ingredientMap map[int64]domain.Ingredient) map[int64][]domain.StepIngredient {
	ingredientsByStep := make(map[int64][]domain.StepIngredient)
	for _, ing := range ingredients {
		stepIngredient := domain.StepIngredient{
			Unit: domain.Unit{
				ID:     ing.Unit.ID,
				Name:   ing.Unit.Name,
				Symbol: ing.Unit.Symbol,
			},
			Amount:     ing.RecipeIngredient.Amount,
			Ingredient: ingredientMap[ing.Ingredient.ID],
		}
		ingredientsByStep[ing.StepID] = append(ingredientsByStep[ing.StepID], stepIngredient)
	}
	return ingredientsByStep
}

func (s *Store) groupStepsByRecipe(steps []recipeStep, ingredientsByStep map[int64][]domain.StepIngredient) map[int64][]domain.RecipeStep {
	stepsByRecipe := make(map[int64][]recipeStep)
	for _, step := range steps {
		stepsByRecipe[step.RecipeID] = append(stepsByRecipe[step.RecipeID], step)
	}

	result := make(map[int64][]domain.RecipeStep)
	for recipeID, recipeSteps := range stepsByRecipe {
		stepsForRecipe := make([]domain.RecipeStep, len(recipeSteps))
		for i, step := range recipeSteps {
			recipeStep := domain.RecipeStep{
				ID:           step.ID,
				Instructions: step.Instructions,
				Ingredients:  ingredientsByStep[step.ID],
			}
			stepsForRecipe[i] = recipeStep
		}
		result[recipeID] = stepsForRecipe
	}
	return result
}

func (s *Store) UpdateRecipe(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		_, err := queries.UpdateRecipe(recipe.ID, recipe.Name, recipe.Description, recipe.Servings, recipe.Minutes).Exec(tx.DB())
		if err != nil {
			return err
		}

		if err = tx.deleteRecipeIngredients(ctx, recipe.ID); err != nil {
			return err
		}
		if err = tx.deleteRecipeSteps(ctx, recipe.ID); err != nil {
			return err
		}
		if err = tx.deleteRecipeImages(ctx, recipe.ID); err != nil {
			return err
		}
		if err = tx.deleteRecipeTags(ctx, recipe.ID); err != nil {
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

func (s *Store) deleteRecipeIngredients(ctx context.Context, recipeID int64) error {
	_, err := queries.DeleteRecipeIngredients(recipeID).Exec(s.DB())
	return err
}

func (s *Store) deleteRecipeSteps(ctx context.Context, recipeID int64) error {
	_, err := queries.DeleteRecipeSteps(recipeID).Exec(s.DB())
	return err
}

func (s *Store) deleteRecipeImages(ctx context.Context, recipeID int64) error {
	_, err := queries.DeleteRecipeImages(recipeID).Exec(s.DB())
	return err
}

func (s *Store) deleteRecipeTags(ctx context.Context, recipeID int64) error {
	_, err := queries.DeleteRecipeTags(recipeID).Exec(s.DB())
	return err
}
