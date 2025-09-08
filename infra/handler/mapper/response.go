package mapper

import (
	"net/url"
	"time"

	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/config"
)

func (m *APIMapper) ToIngredient(ingredient domain.Ingredient) (api.Ingredient, error) {
	return api.Ingredient{
		ID:   ingredient.ID,
		Name: ingredient.Name,
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
		response, err := m.ToRecipe(recipe)
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

func (m *APIMapper) ToRecipe(recipe domain.Recipe) (*api.ReadRecipe, error) {
	images, err := m.ToRecipeImageURLs(recipe.Images)
	if err != nil {
		return nil, err
	}

	steps := make([]api.WriteRecipeStep, len(recipe.Steps))
	for i, step := range recipe.Steps {
		ingredients := make([]api.WriteStepIngredient, len(step.Ingredients))
		for j, ingredient := range step.Ingredients {
			ingredients[j] = api.WriteStepIngredient{
				IngredientId: ingredient.Ingredient.ID,
				UnitId:       ingredient.Unit.ID,
				Amount:       ingredient.Amount,
			}
		}
		steps[i] = api.WriteRecipeStep{
			Instructions: step.Instructions,
		}
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
		mapped, err := m.ToRecipe(recipe)
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

func (m *APIMapper) ToImageURL(imagePath string) (*url.URL, error) {
	path, err := url.JoinPath(m.baseURL, config.ImagesPathPrefix, imagePath)
	if err != nil {
		return nil, err
	}
	return url.Parse(path)
}

func (m *APIMapper) ToRecipeImageURLs(images []domain.RecipeImage) ([]url.URL, error) {
	urls := make([]url.URL, len(images))
	for i, image := range images {
		imageURL, err := m.ToImageURL(image.Path)
		if err != nil {
			return nil, err
		}
		urls[i] = *imageURL
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
