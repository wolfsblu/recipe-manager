package sqlite

import (
	"context"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

func (s *Store) GetIngredients(ctx context.Context) ([]domain.Ingredient, error) {
	result, err := s.query().GetIngredients(ctx)
	if err != nil {
		return nil, err
	}

	ingredients := make([]domain.Ingredient, len(result))
	for i, ingredient := range result {
		ingredients[i] = s.mapper.ToIngredient(ingredient)
	}

	return s.populateIngredientNutrients(ctx, ingredients)
}

func (s *Store) populateIngredientNutrients(ctx context.Context, ingredients []domain.Ingredient) ([]domain.Ingredient, error) {
	if len(ingredients) == 0 {
		return ingredients, nil
	}

	ingredientIDs := make([]int64, len(ingredients))
	for i, ing := range ingredients {
		ingredientIDs[i] = ing.ID
	}

	nutrients, err := s.query().GetNutrientsForIngredients(ctx, ingredientIDs)
	if err != nil {
		return nil, err
	}

	// Group nutrients by ingredient ID
	nutrientsByIngredient := make(map[int64][]domain.IngredientNutrient)
	for _, nutrient := range nutrients {
		nutrientsByIngredient[nutrient.IngredientID] = append(
			nutrientsByIngredient[nutrient.IngredientID],
			domain.IngredientNutrient{
				Nutrient: s.mapper.ToNutrient(nutrient.Nutrient),
				Amount:   nutrient.Amount,
			},
		)
	}

	// Assign nutrients to ingredients
	populatedIngredients := make([]domain.Ingredient, len(ingredients))
	for i, ingredient := range ingredients {
		ingredient.Nutrients = nutrientsByIngredient[ingredient.ID]
		if ingredient.Nutrients == nil {
			ingredient.Nutrients = []domain.IngredientNutrient{}
		}
		populatedIngredients[i] = ingredient
	}

	return populatedIngredients, nil
}

func (s *Store) CreateIngredient(ctx context.Context, ingredient domain.Ingredient) (domain.Ingredient, error) {
	var id int64
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		var err error
		id, err = tx.query().CreateIngredient(ctx, ingredient.Name)
		if err != nil {
			return err
		}

		for _, nutrient := range ingredient.Nutrients {
			err = tx.query().AddIngredientNutrient(ctx, tx.mapper.FromIngredientNutrient(id, nutrient))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return domain.Ingredient{}, err
	}
	ingredient.ID = id
	populated, err := s.populateIngredientNutrients(ctx, []domain.Ingredient{ingredient})
	if err != nil {
		return domain.Ingredient{}, err
	}
	return populated[0], nil
}

func (s *Store) UpdateIngredient(ctx context.Context, ingredient domain.Ingredient) (domain.Ingredient, error) {
	err := s.WithTransaction(ctx, func(tx *TxStore) error {
		err := tx.query().UpdateIngredient(ctx, database.UpdateIngredientParams{
			Name: ingredient.Name,
			ID:   ingredient.ID,
		})
		if err != nil {
			return err
		}

		// Delete existing nutrients and add new ones
		err = tx.query().DeleteIngredientNutrients(ctx, ingredient.ID)
		if err != nil {
			return err
		}

		for _, nutrient := range ingredient.Nutrients {
			err = tx.query().AddIngredientNutrient(ctx, tx.mapper.FromIngredientNutrient(ingredient.ID, nutrient))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return domain.Ingredient{}, err
	}
	populated, err := s.populateIngredientNutrients(ctx, []domain.Ingredient{ingredient})
	if err != nil {
		return domain.Ingredient{}, err
	}
	return populated[0], nil
}

func (s *Store) DeleteIngredient(ctx context.Context, id int64) error {
	return s.query().DeleteIngredient(ctx, id)
}
