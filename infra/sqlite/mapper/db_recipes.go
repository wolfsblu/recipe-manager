package mapper

import (
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/gen/model"
)

func (m *DBMapper) ToIngredient(r model.Ingredients) domain.Ingredient {
	return domain.Ingredient{
		ID:        r.ID,
		Name:      r.Name,
		Nutrients: []domain.IngredientNutrient{},
	}
}

func (m *DBMapper) ToNutrient(r model.Nutrient) domain.Nutrient {
	return domain.Nutrient{
		ID:   r.ID,
		Name: r.Name,
		Unit: r.Unit,
	}
}

func (m *DBMapper) ToRecipeFromModel(r model.Recipe) domain.Recipe {
	return domain.Recipe{
		ID:        r.ID,
		CreatedAt: r.CreatedAt,
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

func (m *DBMapper) ToTagFromModel(t model.Tag) domain.Tag {
	return domain.Tag{
		ID:   t.ID,
		Name: t.Name,
	}
}

func (m *DBMapper) ToUnit(u model.Unit) domain.Unit {
	return domain.Unit{
		ID:     u.ID,
		Name:   u.Name,
		Symbol: u.Symbol,
	}
}

func (m *DBMapper) ToNutrientFromModel(n model.Nutrient) domain.Nutrient {
	return domain.Nutrient{
		ID:   n.ID,
		Name: n.Name,
		Unit: n.Unit,
	}
}
