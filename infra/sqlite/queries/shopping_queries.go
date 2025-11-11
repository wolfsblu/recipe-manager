package queries

import (
	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/wolfsblu/recipe-manager/infra/sqlite/gen/table"
)

// SelectShoppingListsByUser returns a paginated query for user's shopping lists
func SelectShoppingListsByUser(userID int64, lastID int64, lastName string, limit int64) SelectStatement {
	return SELECT(
		ShoppingLists.ID,
		ShoppingLists.UserID,
		ShoppingLists.Name,
	).FROM(ShoppingLists).
		WHERE(
			ShoppingLists.UserID.EQ(Int(userID)).
				AND(ROW(ShoppingLists.Name, ShoppingLists.ID).GT(ROW(String(lastName), Int(lastID)))),
		).ORDER_BY(ShoppingLists.Name.ASC(), ShoppingLists.ID.ASC()).
		LIMIT(limit)
}

// SelectShoppingListByID returns a query for a single shopping list
func SelectShoppingListByID(listID int64) SelectStatement {
	return SELECT(
		ShoppingLists.ID,
		ShoppingLists.UserID,
		ShoppingLists.Name,
	).FROM(ShoppingLists).
		WHERE(ShoppingLists.ID.EQ(Int(listID)))
}

// SelectShoppingListItems returns a query for all items in a shopping list
func SelectShoppingListItems(listID int64) SelectStatement {
	return SELECT(
		ShoppingListItems.ID,
		ShoppingListItems.ShoppingListID,
		ShoppingListItems.Ingredient,
		ShoppingListItems.Quantity,
		ShoppingListItems.Unit,
		ShoppingListItems.Done,
		ShoppingListItems.SortOrder,
	).FROM(ShoppingListItems).
		WHERE(ShoppingListItems.ShoppingListID.EQ(Int(listID))).
		ORDER_BY(ShoppingListItems.SortOrder.ASC())
}

// SelectShoppingListItemByID returns a query for a single shopping list item
func SelectShoppingListItemByID(itemID int64) SelectStatement {
	return SELECT(
		ShoppingListItems.ID,
		ShoppingListItems.ShoppingListID,
		ShoppingListItems.Ingredient,
		ShoppingListItems.Quantity,
		ShoppingListItems.Unit,
		ShoppingListItems.Done,
		ShoppingListItems.SortOrder,
	).FROM(ShoppingListItems).
		WHERE(ShoppingListItems.ID.EQ(Int(itemID)))
}

// InsertShoppingList returns an insert statement for a new shopping list
func InsertShoppingList(userID int64, name string) InsertStatement {
	return ShoppingLists.INSERT(ShoppingLists.UserID, ShoppingLists.Name).
		VALUES(userID, name)
}

// InsertShoppingListItem returns an insert statement for a new shopping list item
func InsertShoppingListItem(listID int64, ingredient string, quantity, unit *string, done bool, sortOrder int64) InsertStatement {
	return ShoppingListItems.INSERT(
		ShoppingListItems.ShoppingListID,
		ShoppingListItems.Ingredient,
		ShoppingListItems.Quantity,
		ShoppingListItems.Unit,
		ShoppingListItems.Done,
		ShoppingListItems.SortOrder,
	).VALUES(listID, ingredient, quantity, unit, done, sortOrder)
}

// UpdateShoppingList returns an update statement for a shopping list
func UpdateShoppingList(listID int64, name string) UpdateStatement {
	return ShoppingLists.UPDATE(ShoppingLists.Name).
		SET(name).
		WHERE(ShoppingLists.ID.EQ(Int(listID)))
}

// UpdateShoppingListItem returns an update statement for a shopping list item
func UpdateShoppingListItem(itemID int64, ingredient string, quantity, unit *string, done bool) UpdateStatement {
	return ShoppingListItems.UPDATE(
		ShoppingListItems.Ingredient,
		ShoppingListItems.Quantity,
		ShoppingListItems.Unit,
		ShoppingListItems.Done,
	).SET(ingredient, quantity, unit, done).
		WHERE(ShoppingListItems.ID.EQ(Int(itemID)))
}

// DeleteShoppingList returns a delete statement for a shopping list
func DeleteShoppingList(listID int64) DeleteStatement {
	return ShoppingLists.DELETE().WHERE(ShoppingLists.ID.EQ(Int(listID)))
}

// DeleteShoppingListItem returns a delete statement for a shopping list item
func DeleteShoppingListItem(itemID int64) DeleteStatement {
	return ShoppingListItems.DELETE().WHERE(ShoppingListItems.ID.EQ(Int(itemID)))
}
