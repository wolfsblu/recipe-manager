package sqlite

import (
	"context"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/database"
	_ "modernc.org/sqlite"
)

func (s *Store) GetShoppingListsByUser(ctx context.Context, userID int64) ([]domain.ShoppingList, error) {
	result, err := s.query().GetShoppingListsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var lists []domain.ShoppingList
	for _, row := range result {
		list := domain.ShoppingList{
			ID:     row.ID,
			UserID: row.UserID,
			Name:   row.Name,
		}

		items, err := s.query().GetShoppingListItemsByListID(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		for _, item := range items {
			list.Items = append(list.Items, domain.ShoppingListItem{
				ID:         item.ID,
				Ingredient: item.Ingredient,
				Quantity:   item.Quantity,
				Unit:       item.Unit,
				Done:       item.Done,
				SortOrder:  item.SortOrder,
			})
		}

		lists = append(lists, list)
	}

	return lists, nil
}

func (s *Store) GetShoppingListByID(ctx context.Context, listID int64) (domain.ShoppingList, error) {
	row, err := s.query().GetShoppingListByID(ctx, listID)
	if err != nil {
		return domain.ShoppingList{}, err
	}

	list := domain.ShoppingList{
		ID:     row.ID,
		UserID: row.UserID,
		Name:   row.Name,
	}

	items, err := s.query().GetShoppingListItemsByListID(ctx, listID)
	if err != nil {
		return domain.ShoppingList{}, err
	}

	for _, item := range items {
		list.Items = append(list.Items, domain.ShoppingListItem{
			ID:         item.ID,
			Ingredient: item.Ingredient,
			Quantity:   (item.Quantity),
			Unit:       (item.Unit),
			Done:       item.Done,
			SortOrder:  item.SortOrder,
		})
	}

	return list, nil
}

func (s *Store) CreateShoppingList(ctx context.Context, userID int64, name string) (domain.ShoppingList, error) {
	row, err := s.query().CreateShoppingList(ctx, database.CreateShoppingListParams{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		return domain.ShoppingList{}, err
	}

	return domain.ShoppingList{
		ID:     row.ID,
		UserID: row.UserID,
		Name:   row.Name,
		Items:  []domain.ShoppingListItem{},
	}, nil
}

func (s *Store) UpdateShoppingList(ctx context.Context, listID int64, name string) (domain.ShoppingList, error) {
	row, err := s.query().UpdateShoppingList(ctx, database.UpdateShoppingListParams{
		Name: name,
		ID:   listID,
	})
	if err != nil {
		return domain.ShoppingList{}, err
	}

	return s.GetShoppingListByID(ctx, row.ID)
}

func (s *Store) DeleteShoppingList(ctx context.Context, listID int64) error {
	return s.query().DeleteShoppingList(ctx, listID)
}

func (s *Store) CreateShoppingListItem(ctx context.Context, listID int64, item domain.ShoppingListItem) (domain.ShoppingListItem, error) {
	row, err := s.query().CreateShoppingListItem(ctx, database.CreateShoppingListItemParams{
		ShoppingListID: listID,
		Ingredient:     item.Ingredient,
		Quantity:       item.Quantity,
		Unit:           item.Unit,
		Done:           item.Done,
		SortOrder:      item.SortOrder,
	})
	if err != nil {
		return domain.ShoppingListItem{}, err
	}

	return domain.ShoppingListItem{
		ID:         row.ID,
		Ingredient: row.Ingredient,
		Quantity:   (row.Quantity),
		Unit:       (row.Unit),
		Done:       row.Done,
		SortOrder:  row.SortOrder,
	}, nil
}

func (s *Store) UpdateShoppingListItem(ctx context.Context, itemID int64, item domain.ShoppingListItem) (domain.ShoppingListItem, error) {
	row, err := s.query().UpdateShoppingListItem(ctx, database.UpdateShoppingListItemParams{
		Ingredient: item.Ingredient,
		Quantity:   item.Quantity,
		Unit:       item.Unit,
		Done:       item.Done,
		ID:         itemID,
	})
	if err != nil {
		return domain.ShoppingListItem{}, err
	}

	return domain.ShoppingListItem{
		ID:         row.ID,
		Ingredient: row.Ingredient,
		Quantity:   (row.Quantity),
		Unit:       (row.Unit),
		Done:       row.Done,
		SortOrder:  row.SortOrder,
	}, nil
}

func (s *Store) DeleteShoppingListItem(ctx context.Context, itemID int64) error {
	return s.query().DeleteShoppingListItem(ctx, itemID)
}
