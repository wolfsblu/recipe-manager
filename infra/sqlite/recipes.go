package sqlite

import (
	"context"
	"net/url"
	"sort"
	"time"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/domain/pagination"
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

func (s *Store) GetMealPlan(ctx context.Context, user *domain.User, from time.Time, until time.Time, req pagination.Page) (pagination.Result[domain.MealPlan], error) {
	result, err := s.query().GetMealPlan(ctx, database.GetMealPlanParams{
		UserID:    user.ID,
		FromDate:  from.Format(time.DateOnly),
		UntilDate: until.Format(time.DateOnly),
		Cursor:    pagination.DecodeCursor(req.Cursor),
		Limit:     int64(req.Limit + 1),
	})
	if err != nil {
		return pagination.Result[domain.MealPlan]{}, err
	}

	recipeMap := make(map[int64]domain.Recipe)
	for _, item := range result {
		recipe := s.mapper.ToRecipe(item.Recipe)
		recipeMap[recipe.ID] = recipe
	}

	recipes := make([]domain.Recipe, 0, len(recipeMap))
	for _, recipe := range recipeMap {
		recipes = append(recipes, recipe)
	}

	populatedRecipes, err := s.populateRecipeRelations(ctx, user, recipes)
	if err != nil {
		return pagination.Result[domain.MealPlan]{}, err
	}

	populatedRecipeMap := make(map[int64]domain.Recipe)
	for _, recipe := range populatedRecipes {
		populatedRecipeMap[recipe.ID] = recipe
	}

	grouped := make(map[string][]domain.Recipe)
	for _, item := range result {
		recipe := populatedRecipeMap[item.Recipe.ID]
		grouped[item.MealPlan.Date] = append(grouped[item.MealPlan.Date], recipe)
	}

	i := 0
	mealPlan := make([]domain.MealPlan, 0, len(grouped))
	for key, recipes := range grouped {
		date, err := time.Parse(time.DateOnly, key)
		if err != nil {
			return pagination.Result[domain.MealPlan]{}, err
		}
		mealPlan = append(mealPlan, domain.MealPlan{
			Date:    date,
			Recipes: recipes,
		})
		i++
	}

	sort.Slice(mealPlan, func(i, j int) bool {
		return mealPlan[i].Date.Before(mealPlan[j].Date)
	})

	// Note: For meal plan pagination, we use the first meal plan's date hash as the ID
	return pagination.NewResult(mealPlan, req.Limit, func(m domain.MealPlan) int64 {
		// Use date as a unique identifier (convert to unix timestamp)
		return m.Date.Unix()
	}), nil
}

func (s *Store) GetTags(ctx context.Context, req pagination.Page) (pagination.Result[domain.Tag], error) {
	result, err := s.query().GetTags(ctx, database.GetTagsParams{
		Cursor: pagination.DecodeCursor(req.Cursor),
		Limit:  int64(req.Limit + 1),
	})
	if err != nil {
		return pagination.Result[domain.Tag]{}, err
	}

	tags := make([]domain.Tag, 0, len(result))
	for _, tag := range result {
		tags = append(tags, s.mapper.ToTag(tag))
	}

	return pagination.NewResult(tags, req.Limit, func(t domain.Tag) int64 {
		return t.ID
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

func (s *Store) GetRecipesByUser(ctx context.Context, user *domain.User, req pagination.Page) (pagination.Result[domain.Recipe], error) {
	result, err := s.query().ListRecipes(ctx, database.ListRecipesParams{
		CreatedBy: user.ID,
		Cursor:    pagination.DecodeCursor(req.Cursor),
		Limit:     int64(req.Limit + 1),
	})
	if err != nil {
		return pagination.Result[domain.Recipe]{}, err
	}

	recipes := make([]domain.Recipe, 0, len(result))
	for _, recipe := range result {
		recipes = append(recipes, s.mapper.ToRecipe(recipe))
	}

	// Populate relations before creating paginated result
	populatedRecipes, err := s.populateRecipeRelations(ctx, user, recipes)
	if err != nil {
		return pagination.Result[domain.Recipe]{}, err
	}

	return pagination.NewResult(populatedRecipes, req.Limit, func(r domain.Recipe) int64 {
		return r.ID
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
