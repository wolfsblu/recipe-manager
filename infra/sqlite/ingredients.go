package sqlite

import (
	"context"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/gen/model"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/queries"
)

func (s *Store) GetIngredients(ctx context.Context, filters domain.IngredientFilters) (domain.Result[domain.Ingredient], error) {
	// Decode cursor - use NameCursor which has both LastID and LastName
	cursor, err := domain.DecodeCursor[*domain.NameCursor](filters.Page.Cursor)
	if err != nil {
		cursor = &domain.NameCursor{}
	}

	var result []model.Ingredient
	err = queries.SelectIngredientsPaginated(
		cursor.LastID,
		cursor.LastName,
		int64(filters.Page.Limit+1),
		filters.SortBy,
		filters.SortOrder,
		filters.Search,
	).Query(s.DB(), &result)
	if err != nil {
		return domain.Result[domain.Ingredient]{}, err
	}

	ingredients := make([]domain.Ingredient, len(result))
	for i, ingredient := range result {
		ingredients[i] = s.mapper.ToIngredient(ingredient)
	}

	populatedIngredients, err := s.populateIngredientNutrients(ctx, ingredients)
	if err != nil {
		return domain.Result[domain.Ingredient]{}, err
	}

	// Always use NameCursor since it has both LastID and LastName fields needed for any sort
	return domain.NewPagedResult(populatedIngredients, filters.Page.Limit, func(i domain.Ingredient) domain.NameCursor {
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

	type NutrientRow struct {
		IngredientID int64
		Nutrient     model.Nutrient
		Amount       float64
	}
	var nutrients []NutrientRow
	err := queries.SelectNutrientsForIngredients(ingredientIDs).Query(s.DB(), &nutrients)
	if err != nil {
		return nil, err
	}

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
		var result model.Ingredient
		err := queries.InsertIngredient(ingredient.Name).Query(tx.DB(), &result)
		if err != nil {
			return err
		}
		id = result.ID

		for _, nutrient := range ingredient.Nutrients {
			_, err = queries.InsertIngredientNutrient(id, nutrient.Nutrient.ID, nutrient.Amount).Exec(tx.DB())
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
		_, err := queries.UpdateIngredient(ingredient.ID, ingredient.Name).Exec(tx.DB())
		if err != nil {
			return err
		}

		_, err = queries.DeleteIngredientNutrients(ingredient.ID).Exec(tx.DB())
		if err != nil {
			return err
		}

		for _, nutrient := range ingredient.Nutrients {
			_, err = queries.InsertIngredientNutrient(ingredient.ID, nutrient.Nutrient.ID, nutrient.Amount).Exec(tx.DB())
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
	_, err := queries.DeleteIngredient(id).Exec(s.DB())
	return err
}
