package sqlite

import (
	"context"
	"time"

	"github.com/wolfsblu/go-chef/domain"
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

	dbRecipe, err := s.query().CreateRecipe(ctx, CreateRecipeParams{
		Name:        recipe.Name,
		Servings:    recipe.Servings,
		Minutes:     recipe.Minutes,
		Description: recipe.Description,
		CreatedBy:   recipe.CreatedBy.ID,
	})
	if err != nil {
		return domain.Recipe{}, err
	}

	result := dbRecipe.AsDomainModel()
	if result.Steps, err = s.createRecipeSteps(ctx, dbRecipe.ID, recipe.Steps); err != nil {
		return domain.Recipe{}, err
	}

	if err = s.Commit(); err != nil {
		return domain.Recipe{}, err
	}

	return result, nil
}

func (s *Store) createRecipeSteps(ctx context.Context, recipeID int64, steps []domain.RecipeStep) ([]domain.RecipeStep, error) {
	for i, step := range steps {
		stepResult, err := s.query().CreateRecipeStep(ctx, CreateRecipeStepParams{
			RecipeID:     recipeID,
			Instructions: step.Instructions,
		})
		if err != nil {
			return nil, err
		}

		steps[i] = domain.RecipeStep{
			ID:           stepResult.ID,
			Instructions: stepResult.Instructions,
		}

		if steps[i].Ingredients, err = s.createStepIngredients(ctx, stepResult.ID, step.Ingredients); err != nil {
			return nil, err
		}
	}
	return steps, nil
}

func (s *Store) createStepIngredients(ctx context.Context, stepID int64, ingredients []domain.StepIngredient) ([]domain.StepIngredient, error) {
	for i, ingredient := range ingredients {
		ingredientResult, err := s.query().CreateStepIngredient(ctx, CreateStepIngredientParams{
			StepID:       stepID,
			IngredientID: ingredient.Ingredient.ID,
			UnitID:       ingredient.Unit.ID,
			Amount:       ingredient.Amount,
		})
		if err != nil {
			return nil, err
		}

		ingredients[i] = domain.StepIngredient{
			Unit: domain.Unit{
				ID: ingredientResult.UnitID,
			},
			Amount: ingredientResult.Amount,
			Ingredient: domain.Ingredient{
				ID: ingredientResult.IngredientID,
			},
		}
	}
	return ingredients, nil
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
	mealplan := make([]domain.MealPlan, len(grouped))
	for key, recipes := range grouped {
		date, err := time.Parse(time.DateOnly, key)
		if err != nil {
			return []domain.MealPlan{}, err
		}
		mealplan[i] = domain.MealPlan{
			Date:    date,
			Recipes: recipes,
		}
		i++
	}
	return mealplan, nil
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
	return result.AsDomainModel(), nil
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
		imagesByRecipe[image.RecipeID] = append(imagesByRecipe[image.RecipeID], domain.RecipeImage{
			ID:        image.ID,
			Path:      image.Path,
			SortOrder: image.SortOrder,
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
