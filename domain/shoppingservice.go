package domain

import (
	"context"
)

type ShoppingService struct {
	store ShoppingStore
}

func (s *ShoppingService) GetByUser(ctx context.Context, user *User) ([]ShoppingList, error) {
	return s.store.GetShoppingListsByUser(ctx, user.ID)
}

func (s *ShoppingService) GetByID(ctx context.Context, user *User, listID int64) (ShoppingList, error) {
	list, err := s.store.GetShoppingListByID(ctx, listID)
	if err != nil {
		return ShoppingList{}, err
	}

	if err := s.validateShoppingListOwnership(ctx, user, &list); err != nil {
		return ShoppingList{}, err
	}

	return list, nil
}

func (s *ShoppingService) Create(ctx context.Context, user *User, name string) (ShoppingList, error) {
	return s.store.CreateShoppingList(ctx, user.ID, name)
}

func (s *ShoppingService) Update(ctx context.Context, user *User, listID int64, name string) (ShoppingList, error) {
	if err := s.validateShoppingListOwnershipByID(ctx, user, listID); err != nil {
		return ShoppingList{}, err
	}

	return s.store.UpdateShoppingList(ctx, listID, name)
}

func (s *ShoppingService) Delete(ctx context.Context, user *User, listID int64) error {
	if err := s.validateShoppingListOwnershipByID(ctx, user, listID); err != nil {
		return err
	}

	return s.store.DeleteShoppingList(ctx, listID)
}

func (s *ShoppingService) AddItem(ctx context.Context, user *User, listID int64, item ShoppingListItem) (ShoppingListItem, error) {
	list, err := s.GetByID(ctx, user, listID)
	if err != nil {
		return ShoppingListItem{}, err
	}

	item.SortOrder = s.calculateNextSortOrder(list.Items)
	return s.store.CreateShoppingListItem(ctx, listID, item)
}

func (s *ShoppingService) UpdateItem(ctx context.Context, user *User, listID int64, itemID int64, item ShoppingListItem) (ShoppingListItem, error) {
	if err := s.validateShoppingListOwnershipByID(ctx, user, listID); err != nil {
		return ShoppingListItem{}, err
	}

	return s.store.UpdateShoppingListItem(ctx, itemID, item)
}

func (s *ShoppingService) DeleteItem(ctx context.Context, user *User, listID int64, itemID int64) error {
	if err := s.validateShoppingListOwnershipByID(ctx, user, listID); err != nil {
		return err
	}

	return s.store.DeleteShoppingListItem(ctx, itemID)
}

func (s *ShoppingService) calculateNextSortOrder(items []ShoppingListItem) int64 {
	maxSortOrder := int64(-1)
	for _, item := range items {
		if item.SortOrder > maxSortOrder {
			maxSortOrder = item.SortOrder
		}
	}
	return maxSortOrder + 1
}
