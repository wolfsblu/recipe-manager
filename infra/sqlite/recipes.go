package sqlite

import (
	"context"
	"net/url"
	"time"

	"github.com/wolfsblu/recipe-manager/domain"
	_ "modernc.org/sqlite"
)

func (s *Store) BrowseRecipes(ctx context.Context) (recipes []domain.Recipe, err error) {
	result, err := s.query().BrowseRecipes(ctx)
	if err != nil {
		return nil, err
	}
	for _, recipe := range result {
		recipes = append(recipes, recipe.AsDomainModel())
	}
	return recipes, err
}

func (s *Store) CreateRecipe(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	if err := s.Begin(ctx); err != nil {
		return domain.Recipe{}, err
	}
	defer s.Rollback()

	recipeId, err := s.query().CreateRecipe(ctx, CreateRecipeParams{
		Name:        recipe.Name,
		Servings:    recipe.Servings,
		Minutes:     recipe.Minutes,
		Description: recipe.Description,
		CreatedBy:   recipe.CreatedBy.ID,
	})
	if err != nil {
		return domain.Recipe{}, err
	}
	if err = s.createRecipeSteps(ctx, recipeId, recipe.Steps); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.createRecipeImages(ctx, recipeId, recipe.Images); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.Commit(); err != nil {
		return domain.Recipe{}, err
	}
	return s.GetRecipeById(ctx, recipeId)
}

func (s *Store) createRecipeSteps(ctx context.Context, recipeID int64, steps []domain.RecipeStep) error {
	for _, step := range steps {
		stepId, err := s.query().CreateRecipeStep(ctx, CreateRecipeStepParams{
			RecipeID:     recipeID,
			Instructions: step.Instructions,
		})
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
		_, err := s.query().CreateStepIngredient(ctx, CreateStepIngredientParams{
			StepID:       stepID,
			IngredientID: ingredient.Ingredient.ID,
			UnitID:       ingredient.Unit.ID,
			Amount:       ingredient.Amount,
			SortOrder:    int64(i),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) createRecipeImages(ctx context.Context, recipeID int64, images []domain.RecipeImage) error {
	for _, image := range images {
		_, err := s.query().CreateRecipeImages(ctx, CreateRecipeImagesParams{
			RecipeID: recipeID,
			Path:     image.URL.String(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) DeleteRecipe(ctx context.Context, id int64) error {
	return s.query().DeleteRecipe(ctx, id)
}

func (s *Store) GetMealPlan(ctx context.Context, user *domain.User, from time.Time, until time.Time) ([]domain.MealPlan, error) {
	result, err := s.query().GetMealPlan(ctx, GetMealPlanParams{
		UserID:    user.ID,
		FromDate:  from.Format(time.DateOnly),
		UntilDate: until.Format(time.DateOnly),
	})
	if err != nil {
		return []domain.MealPlan{}, err
	}

	grouped := make(map[string][]domain.Recipe)
	for _, item := range result {
		grouped[item.MealPlan.Date] = append(grouped[item.MealPlan.Date], item.Recipe.AsDomainModel())
	}

	i := 0
	mealPlan := make([]domain.MealPlan, len(grouped))
	for key, recipes := range grouped {
		date, err := time.Parse(time.DateOnly, key)
		if err != nil {
			return []domain.MealPlan{}, err
		}
		mealPlan[i] = domain.MealPlan{
			Date:    date,
			Recipes: recipes,
		}
		i++
	}
	return mealPlan, nil
}

func (s *Store) GetIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	result, err := s.query().GetIngredients(ctx)
	if err != nil {
		return nil, err
	}

	ingredients := make([]domain.Ingredient, len(result))
	for i, ingredient := range result {
		ingredients[i] = ingredient.AsDomainModel()
	}

	return ingredients, nil
}

func (s *Store) GetUnits(ctx context.Context) ([]domain.Unit, error) {
	result, err := s.query().GetUnits(ctx)
	if err != nil {
		return nil, err
	}

	units := make([]domain.Unit, len(result))
	for i, unit := range result {
		units[i] = unit.AsDomainModel()
	}

	return units, nil
}

func (s *Store) GetRecipeById(ctx context.Context, id int64) (recipe domain.Recipe, _ error) {
	result, err := s.query().GetRecipe(ctx, id)
	if err != nil {
		return recipe, err
	}

	recipe = result.AsDomainModel()
	recipeIds := []int64{id}
	steps, err := s.query().GetStepsForRecipes(ctx, recipeIds)
	if err != nil {
		return recipe, err
	}
	ingredients, err := s.query().GetIngredientsForRecipes(ctx, recipeIds)
	if err != nil {
		return recipe, err
	}
	tags, err := s.query().GetTagsForRecipes(ctx, recipeIds)
	if err != nil {
		return recipe, err
	}
	images, err := s.query().GetImagesForRecipes(ctx, recipeIds)
	if err != nil {
		return recipe, err
	}

	ingredientsByStep := make(map[int64][]domain.StepIngredient)
	for _, ing := range ingredients {
		ingredientsByStep[ing.StepID] = append(ingredientsByStep[ing.StepID], ing.AsDomainModel())
	}

	recipeSteps := make([]domain.RecipeStep, len(steps))
	for i, step := range steps {
		recipeStep := step.AsDomainModel()
		recipeStep.Ingredients = ingredientsByStep[step.ID]
		recipeSteps[i] = recipeStep
	}

	recipeTags := make([]string, len(tags))
	for i, tag := range tags {
		recipeTags[i] = tag.AsDomainModel()
	}

	recipeImages := make([]domain.RecipeImage, len(images))
	for i, image := range images {
		recipeImages[i] = image.AsDomainModel()
		imageUrl, err := url.ParseRequestURI(image.Path)
		if err != nil {
			return recipe, err
		}
		recipeImages[i].URL = imageUrl
	}

	recipe.Steps = recipeSteps
	recipe.Tags = recipeTags
	recipe.Images = recipeImages
	return recipe, nil
}

func (s *Store) GetRecipesByUser(ctx context.Context, user *domain.User) (recipes []domain.Recipe, _ error) {
	result, err := s.query().ListRecipes(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, len(result))
	for i, item := range result {
		ids[i] = item.ID
		recipes = append(recipes, item.AsDomainModel())
	}

	tags, err := s.query().GetTagsForRecipes(ctx, ids)
	if err != nil {
		return nil, err
	}

	images, err := s.query().GetImagesForRecipes(ctx, ids)
	if err != nil {
		return nil, err
	}

	tagsByRecipe := make(map[int64][]string)
	for _, tag := range tags {
		tagsByRecipe[tag.RecipeID] = append(tagsByRecipe[tag.RecipeID], tag.Name)
	}

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

	recipes = make([]domain.Recipe, len(result))
	for i, recipe := range result {
		r := recipe.AsDomainModel()
		r.Images = imagesByRecipe[r.ID]
		r.Tags = tagsByRecipe[r.ID]
		recipes[i] = r
	}
	return recipes, nil
}

func (s *Store) UpdateRecipe(ctx context.Context, recipe domain.Recipe) (domain.Recipe, error) {
	if err := s.Begin(ctx); err != nil {
		return domain.Recipe{}, err
	}
	defer s.Rollback()

	err := s.query().UpdateRecipe(ctx, UpdateRecipeParams{
		Name:        recipe.Name,
		Servings:    recipe.Servings,
		Minutes:     recipe.Minutes,
		Description: recipe.Description,
		ID:          recipe.ID,
	})
	if err != nil {
		return domain.Recipe{}, err
	}

	if err = s.query().DeleteRecipeSteps(ctx, recipe.ID); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.query().DeleteRecipeImages(ctx, recipe.ID); err != nil {
		return domain.Recipe{}, err
	}

	if err = s.createRecipeSteps(ctx, recipe.ID, recipe.Steps); err != nil {
		return domain.Recipe{}, err
	}
	if err = s.createRecipeImages(ctx, recipe.ID, recipe.Images); err != nil {
		return domain.Recipe{}, err
	}

	if err = s.Commit(); err != nil {
		return domain.Recipe{}, err
	}
	return s.GetRecipeById(ctx, recipe.ID)
}
