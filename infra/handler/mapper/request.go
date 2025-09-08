package mapper

import (
	"github.com/wolfsblu/go-chef/api"
	"github.com/wolfsblu/go-chef/domain"
)

func (m *APIMapper) FromWriteRecipe(req *api.WriteRecipe) domain.Recipe {
	steps := make([]domain.RecipeStep, len(req.Steps))
	for i, step := range req.Steps {
		steps[i] = m.fromWriteRecipeStep(step)
	}
	/*
		images := make([]domain.RecipeImage, len(req.Images))
		for i, image := range req.Images {
			images[i] = domain.RecipeImage{
				Path:      image,
				SortOrder: int64(i),
			}
		}
	*/
	return domain.Recipe{
		//Images: images,
		Tags:  req.Tags,
		Steps: steps,
		RecipeDetails: domain.RecipeDetails{
			Name:        req.Name,
			Description: req.Description,
			Servings:    req.Servings,
			Minutes:     req.Minutes,
		},
	}
}

func (m *APIMapper) fromWriteRecipeStep(step api.WriteRecipeStep) domain.RecipeStep {
	ingredients := make([]domain.StepIngredient, len(step.Ingredients))
	for i, ingredient := range step.Ingredients {
		ingredients[i] = m.fromWriteStepIngredient(ingredient)
	}

	return domain.RecipeStep{
		Instructions: step.Instructions,
		Ingredients:  ingredients,
	}
}

func (m *APIMapper) fromWriteStepIngredient(ingredient api.WriteStepIngredient) domain.StepIngredient {
	return domain.StepIngredient{
		Unit: domain.Unit{
			ID: ingredient.UnitId,
		},
		Amount: ingredient.Amount,
		Ingredient: domain.Ingredient{
			ID: ingredient.IngredientId,
		},
	}
}
