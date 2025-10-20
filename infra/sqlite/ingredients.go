package sqlite

import (
	"context"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
)

func (s *Store) GetIngredients(ctx context.Context, page domain.Page) (domain.Result[domain.Ingredient], error) {
	cursor, err := domain.DecodeCursor[*domain.NameCursor](page.Cursor)
	if err != nil {
		// Empty cursor is valid for first page
		cursor = &domain.NameCursor{}
	}
	result, err := s.query().GetIngredients(ctx, database.GetIngredientsParams{
		LastName: cursor.LastName,
		LastID:   cursor.LastID,
		Limit:    int64(page.Limit + 1),
	})
	if err != nil {
		return domain.Result[domain.Ingredient]{}, err
	}

	ingredients := make([]domain.Ingredient, len(result))
	for i, ingredient := range result {
		ingredients[i] = s.mapper.ToIngredient(ingredient)
	}

	// Populate nutrients before creating paginated result
	populatedIngredients, err := s.populateIngredientNutrients(ctx, ingredients)
	if err != nil {
		return domain.Result[domain.Ingredient]{}, err
	}

	return domain.NewPagedResult(populatedIngredients, page.Limit, func(i domain.Ingredient) domain.NameCursor {
		return domain.NameCursor{
			LastID:   i.ID,
			LastName: i.Name,
		}
	}), nil
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
