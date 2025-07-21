package mapper

import (
	"time"

	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/urlbuilder"
)

type RecipeMapper struct {
	URLBuilder *urlbuilder.Builder
}

func NewRecipeMapper(urlBuilder *urlbuilder.Builder) *RecipeMapper {
	return &RecipeMapper{
		URLBuilder: urlBuilder,
	}
}

func (m *RecipeMapper) ToRecipeResponse(recipe domain.Recipe) (*api.ReadRecipe, error) {
	images, err := m.URLBuilder.BuildRecipeImageURLs(recipe.Images)
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

func (m *RecipeMapper) ToRecipeListResponse(recipes []domain.Recipe) ([]api.ReadRecipe, error) {
	result := make([]api.ReadRecipe, len(recipes))
	for i, recipe := range recipes {
		mapped, err := m.ToRecipeResponse(recipe)
		if err != nil {
			return nil, err
		}
		result[i] = *mapped
	}
	return result, nil
}

func (m *RecipeMapper) ToMealPlanResponse(mealPlan domain.MealPlan) (api.ReadMealPlan, error) {
	recipes := make([]api.ReadRecipe, len(mealPlan.Recipes))
	for i, recipe := range mealPlan.Recipes {
		response, err := m.ToRecipeResponse(recipe)
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

func (m *RecipeMapper) ToReadMealPlanList(mealPlans []domain.MealPlan) ([]api.ReadMealPlan, error) {
	result := make([]api.ReadMealPlan, len(mealPlans))
	for i, mealPlan := range mealPlans {
		response, err := m.ToMealPlanResponse(mealPlan)
		if err != nil {
			return nil, err
		}
		result[i] = response
	}
	return result, nil
}
