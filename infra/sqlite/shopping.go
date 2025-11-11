package sqlite

import (
	"context"

	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/gen/model"
	"github.com/wolfsblu/recipe-manager/infra/sqlite/queries"
	_ "modernc.org/sqlite"
)

func (s *Store) GetShoppingListsByUser(ctx context.Context, userID int64, req domain.Page) (domain.Result[domain.ShoppingList], error) {
	cursor, err := domain.DecodeCursor[*domain.NameCursor](req.Cursor)
	if err != nil {
		cursor = &domain.NameCursor{}
	}

	var result []model.ShoppingList
	err = queries.SelectShoppingListsByUser(userID, cursor.LastID, cursor.LastName, int64(req.Limit+1)).Query(s.DB(), &result)
	if err != nil {
		return domain.Result[domain.ShoppingList]{}, err
	}

	var lists []domain.ShoppingList
	for _, row := range result {
		list := domain.ShoppingList{
			ID:     row.ID,
			UserID: row.UserID,
			Name:   row.Name,
		}

		var items []model.ShoppingListItem
		err = queries.SelectShoppingListItems(row.ID).Query(s.DB(), &items)
		if err != nil {
			return domain.Result[domain.ShoppingList]{}, err
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

	return domain.NewPagedResult(lists, req.Limit, func(l domain.ShoppingList) domain.NameCursor {
		return domain.NameCursor{
			LastID:   l.ID,
			LastName: l.Name,
		}
	}), nil
}

func (s *Store) GetShoppingListByID(ctx context.Context, listID int64) (domain.ShoppingList, error) {
	var row model.ShoppingList
	err := queries.SelectShoppingListByID(listID).Query(s.DB(), &row)
	if err != nil {
		return domain.ShoppingList{}, err
	}

	list := domain.ShoppingList{
		ID:     row.ID,
		UserID: row.UserID,
		Name:   row.Name,
	}

	var items []model.ShoppingListItem
	err = queries.SelectShoppingListItems(listID).Query(s.DB(), &items)
	if err != nil {
		return domain.ShoppingList{}, err
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

	return list, nil
}

func (s *Store) CreateShoppingList(ctx context.Context, userID int64, name string) (domain.ShoppingList, error) {
	var row model.ShoppingList
	err := queries.InsertShoppingList(userID, name).Query(s.DB(), &row)
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
	_, err := queries.UpdateShoppingList(listID, name).Exec(s.DB())
	if err != nil {
		return domain.ShoppingList{}, err
	}

	return s.GetShoppingListByID(ctx, listID)
}

func (s *Store) DeleteShoppingList(ctx context.Context, listID int64) error {
	_, err := queries.DeleteShoppingList(listID).Exec(s.DB())
	return err
}

func (s *Store) CreateShoppingListItem(ctx context.Context, listID int64, item domain.ShoppingListItem) (domain.ShoppingListItem, error) {
	var row model.ShoppingListItem
	err := queries.InsertShoppingListItem(listID, item.Ingredient, item.Quantity, item.Unit, item.Done, item.SortOrder).Query(s.DB(), &row)
	if err != nil {
		return domain.ShoppingListItem{}, err
	}

	return domain.ShoppingListItem{
		ID:         row.ID,
		Ingredient: row.Ingredient,
		Quantity:   row.Quantity,
		Unit:       row.Unit,
		Done:       row.Done,
		SortOrder:  row.SortOrder,
	}, nil
}

func (s *Store) UpdateShoppingListItem(ctx context.Context, itemID int64, item domain.ShoppingListItem) (domain.ShoppingListItem, error) {
	_, err := queries.UpdateShoppingListItem(itemID, item.Ingredient, item.Quantity, item.Unit, item.Done).Exec(s.DB())
	if err != nil {
		return domain.ShoppingListItem{}, err
	}

	var row model.ShoppingListItem
	err = queries.SelectShoppingListItemByID(itemID).Query(s.DB(), &row)
	if err != nil {
		return domain.ShoppingListItem{}, err
	}

	return domain.ShoppingListItem{
		ID:         row.ID,
		Ingredient: row.Ingredient,
		Quantity:   row.Quantity,
		Unit:       row.Unit,
		Done:       row.Done,
		SortOrder:  row.SortOrder,
	}, nil
}

func (s *Store) DeleteShoppingListItem(ctx context.Context, itemID int64) error {
	_, err := queries.DeleteShoppingListItem(itemID).Exec(s.DB())
	return err
}
