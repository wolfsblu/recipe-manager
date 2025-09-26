package mapper

import (
	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
)

func (m *APIMapper) ToWriteRecipe(recipe domain.Recipe) (*api.WriteRecipe, error) {
	images, err := m.ToRecipeImageURLs(recipe.Images)
	if err != nil {
		return nil, err
	}

	steps := make([]api.WriteRecipeStep, len(recipe.Steps))
	for i, step := range recipe.Steps {
		writeStep, err := m.toWriteRecipeStep(step)
		if err != nil {
			return nil, err
		}
		steps[i] = writeStep
	}

	tagIDs := make([]int64, len(recipe.Tags))
	for i, tag := range recipe.Tags {
		tagIDs[i] = tag.ID
	}

	return &api.WriteRecipe{
		Name:        recipe.Name,
		Description: recipe.Description,
		Servings:    recipe.Servings,
		Minutes:     recipe.Minutes,
		Images:      images,
		Steps:       steps,
		Tags:        tagIDs,
	}, nil
}

func (m *APIMapper) toWriteRecipeStep(step domain.RecipeStep) (api.WriteRecipeStep, error) {
	ingredients := make([]api.WriteStepIngredient, len(step.Ingredients))
	for i, ingredient := range step.Ingredients {
		ingredients[i] = m.toWriteStepIngredient(ingredient)
	}

	return api.WriteRecipeStep{
		Instructions: step.Instructions,
		Ingredients:  ingredients,
	}, nil
}

func (m *APIMapper) toWriteStepIngredient(ingredient domain.StepIngredient) api.WriteStepIngredient {
	return api.WriteStepIngredient{
		IngredientId: ingredient.Ingredient.ID,
		UnitId:       ingredient.Unit.ID,
		Amount:       ingredient.Amount,
	}
}
