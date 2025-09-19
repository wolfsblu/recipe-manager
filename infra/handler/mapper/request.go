package mapper

import (
	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
)

func (m *APIMapper) FromWriteRecipe(req *api.WriteRecipe) domain.Recipe {
	steps := make([]domain.RecipeStep, len(req.Steps))
	for i, step := range req.Steps {
		steps[i] = m.fromWriteRecipeStep(step)
	}
	images := make([]domain.RecipeImage, len(req.Images))
	for i, image := range req.Images {
		images[i] = domain.RecipeImage{
			URL: &image,
		}
	}
	return domain.Recipe{
		Images: images,
		Tags:   m.fromTagIDs(req.Tags),
		Steps:  steps,
		RecipeDetails: domain.RecipeDetails{
			Name:        req.Name,
			Description: req.Description,
			Servings:    req.Servings,
			Minutes:     req.Minutes,
		},
	}
}

func (m *APIMapper) fromTagIDs(ids []int64) []domain.Tag {
	tags := make([]domain.Tag, len(ids))
	for i, id := range ids {
		tags[i] = domain.Tag{ID: id}
	}
	return tags
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

func (m *APIMapper) FromWriteIngredient(req *api.WriteIngredient) domain.Ingredient {
	return domain.Ingredient{
		Name: req.Name,
	}
}

func (m *APIMapper) FromWriteUnit(req *api.WriteUnit) domain.Unit {
	return domain.Unit{
		Name: req.Name,
		Code: FromOptNilString(req.Code),
	}
}

func FromOptNilString(s api.OptNilString) *string {
	if v, ok := s.Get(); ok {
		return &v
	}
	return nil
}
