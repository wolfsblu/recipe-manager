package handler

import (
	"time"

	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
)

type responseMapper struct {
	urlBuilder *urlBuilder
}

func newResponseMapper(urlBuilder *urlBuilder) *responseMapper {
	return &responseMapper{
		urlBuilder: urlBuilder,
	}
}

func (m *responseMapper) toRecipeResponse(recipe domain.Recipe) (*api.ReadRecipe, error) {
	images, err := m.urlBuilder.buildRecipeImageURLs(recipe.Images)
	if err != nil {
		return nil, err
	}

	return &api.ReadRecipe{
		ID:          recipe.ID,
		Name:        recipe.Name,
		Description: recipe.Description,
		Servings:    recipe.Servings,
		Minutes:     recipe.Minutes,
		Images:      images,
		Tags:        recipe.Tags,
	}, nil
}

func (m *responseMapper) toRecipeListResponse(recipes []domain.Recipe) ([]api.ReadRecipe, error) {
	result := make([]api.ReadRecipe, len(recipes))
	for i, recipe := range recipes {
		mapped, err := m.toRecipeResponse(recipe)
		if err != nil {
			return nil, err
		}
		result[i] = *mapped
	}
	return result, nil
}

func (m *responseMapper) toMealPlanResponse(mealPlan domain.MealPlan) (api.ReadMealPlan, error) {
	recipes := make([]api.ReadRecipe, len(mealPlan.Recipes))
	for i, recipe := range mealPlan.Recipes {
		response, err := m.toRecipeResponse(recipe)
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

func (m *responseMapper) toReadMealPlanList(mealPlans []domain.MealPlan) ([]api.ReadMealPlan, error) {
	result := make([]api.ReadMealPlan, len(mealPlans))
	for i, mealPlan := range mealPlans {
		response, err := m.toMealPlanResponse(mealPlan)
		if err != nil {
			return nil, err
		}
		result[i] = response
	}
	return result, nil
}
