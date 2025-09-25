package domain

import (
	"context"
)

func (s *ShoppingService) validateShoppingListOwnership(ctx context.Context, user *User, list *ShoppingList) error {
	if list.UserID != user.ID {
		return ErrAuthorization
	}
	return nil
}

func (s *ShoppingService) validateShoppingListOwnershipByID(ctx context.Context, user *User, listID int64) error {
	list, err := s.store.GetShoppingListByID(ctx, listID)
	if err != nil {
		return err
	}

	return s.validateShoppingListOwnership(ctx, user, &list)
}
