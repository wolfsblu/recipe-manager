package mapper

import (
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

func (m *DBMapper) ToIngredient(r database.Ingredient) domain.Ingredient {
	return domain.Ingredient{
		ID:   r.ID,
		Name: r.Name,
	}
}

func (m *DBMapper) ToRecipe(r database.Recipe) domain.Recipe {
	return domain.Recipe{
		ID: r.ID,
		RecipeDetails: domain.RecipeDetails{
			Name:        r.Name,
			Description: r.Description,
			CreatedBy: &domain.User{
				ID: r.CreatedBy,
			},
			Servings: r.Servings,
			Minutes:  r.Minutes,
		},
	}
}

func (m *DBMapper) ToTag(t database.Tag) domain.Tag {
	return domain.Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}

func (m *DBMapper) ToUnit(u database.Unit) domain.Unit {
	return domain.Unit{
		ID:   u.ID,
		Name: u.Name,
		Code: u.Code,
	}
}

func (m *DBMapper) ToRecipeStep(r database.GetStepsForRecipesRow) domain.RecipeStep {
	return domain.RecipeStep{
		ID:           r.ID,
		Instructions: r.Instructions,
		Ingredients:  []domain.StepIngredient{}, // Will be populated separately
	}
}

func (m *DBMapper) ToStepIngredient(r database.GetIngredientsForRecipesRow) domain.StepIngredient {
	return domain.StepIngredient{
		Unit: domain.Unit{
			ID:   r.UnitID,
			Name: r.UnitName,
			Code: r.UnitCode,
		},
		Ingredient: domain.Ingredient{
			ID:   r.IngredientID,
			Name: r.IngredientName,
		},
		Amount: r.Amount,
	}
}

func (m *DBMapper) ToRecipeImage(r database.GetImagesForRecipesRow) domain.RecipeImage {
	return domain.RecipeImage{
		ID: r.ID,
	}
}
