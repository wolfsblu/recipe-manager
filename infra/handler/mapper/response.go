package mapper

import (
	"net/url"
	"time"

	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
)

func (m *APIMapper) ToIngredient(ingredient domain.Ingredient) (api.Ingredient, error) {
	return api.Ingredient{
		ID:   ingredient.ID,
		Name: ingredient.Name,
	}, nil
}

func (m *APIMapper) ToReadStepIngredient(ingredient domain.StepIngredient) (api.ReadStepIngredient, error) {
	unit, err := m.ToUnit(ingredient.Unit)
	if err != nil {
		return api.ReadStepIngredient{}, err
	}

	apiIngredient, err := m.ToIngredient(ingredient.Ingredient)
	if err != nil {
		return api.ReadStepIngredient{}, err
	}

	return api.ReadStepIngredient{
		Ingredient: apiIngredient,
		Unit:       unit,
		Amount:     ingredient.Amount,
	}, nil
}

func (m *APIMapper) ToReadRecipeStep(step domain.RecipeStep) (api.ReadRecipeStep, error) {
	ingredients := make([]api.ReadStepIngredient, len(step.Ingredients))
	for i, ingredient := range step.Ingredients {
		apiIngredient, err := m.ToReadStepIngredient(ingredient)
		if err != nil {
			return api.ReadRecipeStep{}, err
		}
		ingredients[i] = apiIngredient
	}

	return api.ReadRecipeStep{
		ID:           step.ID,
		Ingredients:  ingredients,
		Instructions: step.Instructions,
	}, nil
}

func (m *APIMapper) ToIngredients(ingredients []domain.Ingredient) ([]api.Ingredient, error) {
	result := make([]api.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		mapped, err := m.ToIngredient(ingredient)
		if err != nil {
			return nil, err
		}
		result[i] = mapped
	}
	return result, nil
}

func (m *APIMapper) ToMealPlan(mealPlan domain.MealPlan) (api.ReadMealPlan, error) {
	recipes := make([]api.ReadRecipe, len(mealPlan.Recipes))
	for i, recipe := range mealPlan.Recipes {
		response, err := m.ToReadRecipe(recipe)
		if err != nil {
			return api.ReadMealPlan{}, err
		}
		recipes[i] = *response
	}

	return api.ReadMealPlan{
		Date:    mealPlan.Date.Format(time.DateOnly),
		Recipes: recipes,
	}, nil
}

func (m *APIMapper) ToMealPlans(mealPlans []domain.MealPlan) ([]api.ReadMealPlan, error) {
	result := make([]api.ReadMealPlan, len(mealPlans))
	for i, mealPlan := range mealPlans {
		response, err := m.ToMealPlan(mealPlan)
		if err != nil {
			return nil, err
		}
		result[i] = response
	}
	return result, nil
}

func (m *APIMapper) ToReadRecipe(recipe domain.Recipe) (*api.ReadRecipe, error) {
	images, err := m.ToRecipeImageURLs(recipe.Images)
	if err != nil {
		return nil, err
	}

	steps := make([]api.ReadRecipeStep, len(recipe.Steps))
	for i, step := range recipe.Steps {
		readStep, err := m.ToReadRecipeStep(step)
		if err != nil {
			return nil, err
		}
		steps[i] = readStep
	}

	return &api.ReadRecipe{
		ID:          recipe.ID,
		Name:        recipe.Name,
		Description: recipe.Description,
		Servings:    recipe.Servings,
		Minutes:     recipe.Minutes,
		Images:      images,
		Tags:        recipe.Tags,
		Steps:       steps,
	}, nil
}

func (m *APIMapper) ToRecipes(recipes []domain.Recipe) ([]api.ReadRecipe, error) {
	result := make([]api.ReadRecipe, len(recipes))
	for i, recipe := range recipes {
		mapped, err := m.ToReadRecipe(recipe)
		if err != nil {
			return nil, err
		}
		result[i] = *mapped
	}
	return result, nil
}

func (m *APIMapper) ToUnit(unit domain.Unit) (api.ReadUnit, error) {
	return api.ReadUnit{
		ID:   unit.ID,
		Name: unit.Name,
		Code: ToNilString(unit.Code),
	}, nil
}

func (m *APIMapper) ToUnits(units []domain.Unit) ([]api.ReadUnit, error) {
	result := make([]api.ReadUnit, len(units))
	for i, unit := range units {
		mapped, err := m.ToUnit(unit)
		if err != nil {
			return nil, err
		}
		result[i] = mapped
	}
	return result, nil
}

func (m *APIMapper) ToRecipeImageURLs(images []domain.RecipeImage) ([]url.URL, error) {
	urls := make([]url.URL, len(images))
	for i, image := range images {
		urls[i] = *image.URL
	}
	return urls, nil
}

func ToNilString(s *string) api.NilString {
	if s != nil {
		return api.NewNilString(*s)
	}
	return api.NilString{
		Value: "",
		Null:  true,
	}
}
