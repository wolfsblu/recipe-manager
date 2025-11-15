package queries

import (
	. "github.com/go-jet/jet/v2/sqlite"
	. "github.com/wolfsblu/recipe-manager/infra/sqlite/gen/table"
)

// SelectShoppingListsByUser returns a paginated query for user's shopping lists
func SelectShoppingListsByUser(userID int64, lastID int64, lastName string, limit int64) SelectStatement {
	return SELECT(
		ShoppingList.ID,
		ShoppingList.UserID,
		ShoppingList.Name,
	).FROM(ShoppingList).
		WHERE(
			ShoppingList.UserID.EQ(Int(userID)).
				AND(ROW(ShoppingList.Name, ShoppingList.ID).GT(ROW(String(lastName), Int(lastID)))),
		).ORDER_BY(ShoppingList.Name.ASC(), ShoppingList.ID.ASC()).
		LIMIT(limit)
}

// SelectShoppingListByID returns a query for a single shopping list
func SelectShoppingListByID(listID int64) SelectStatement {
	return SELECT(
		ShoppingList.ID,
		ShoppingList.UserID,
		ShoppingList.Name,
	).FROM(ShoppingList).
		WHERE(ShoppingList.ID.EQ(Int(listID)))
}

// SelectShoppingListItems returns a query for all items in a shopping list
func SelectShoppingListItems(listID int64) SelectStatement {
	return SELECT(ShoppingListItem.AllColumns).
		FROM(ShoppingListItem).
		WHERE(ShoppingListItem.ShoppingListID.EQ(Int(listID))).
		ORDER_BY(ShoppingListItem.SortOrder.ASC())
}

// SelectShoppingListItemByID returns a query for a single shopping list item
func SelectShoppingListItemByID(itemID int64) SelectStatement {
	return SELECT(ShoppingListItem.AllColumns).
		FROM(ShoppingListItem).
		WHERE(ShoppingListItem.ID.EQ(Int(itemID)))
}

// InsertShoppingList returns an insert statement for a new shopping list
func InsertShoppingList(userID int64, name string) InsertStatement {
	return ShoppingList.INSERT(ShoppingList.UserID, ShoppingList.Name).
		VALUES(userID, name).
		RETURNING(ShoppingList.AllColumns)
}

// InsertShoppingListItem returns an insert statement for a new shopping list item
func InsertShoppingListItem(listID int64, ingredient string, quantity, unit *string, done bool, sortOrder int64) InsertStatement {
	return ShoppingListItem.INSERT(ShoppingListItem.AllColumns).
		VALUES(listID, ingredient, quantity, unit, done, sortOrder).
		RETURNING(ShoppingListItem.AllColumns)
}

// UpdateShoppingList returns an update statement for a shopping list
func UpdateShoppingList(listID int64, name string) UpdateStatement {
	return ShoppingList.UPDATE(ShoppingList.Name).
		SET(name).
		WHERE(ShoppingList.ID.EQ(Int(listID)))
}

// UpdateShoppingListItem returns an update statement for a shopping list item
func UpdateShoppingListItem(itemID int64, ingredient string, quantity, unit *string, done bool) UpdateStatement {
	return ShoppingListItem.UPDATE(
		ShoppingListItem.Ingredient,
		ShoppingListItem.Quantity,
		ShoppingListItem.Unit,
		ShoppingListItem.Done,
	).SET(ingredient, quantity, unit, done).
		WHERE(ShoppingListItem.ID.EQ(Int(itemID)))
}

// DeleteShoppingList returns a delete statement for a shopping list
func DeleteShoppingList(listID int64) DeleteStatement {
	return ShoppingList.DELETE().WHERE(ShoppingList.ID.EQ(Int(listID)))
}

// DeleteShoppingListItem returns a delete statement for a shopping list item
func DeleteShoppingListItem(itemID int64) DeleteStatement {
	return ShoppingListItem.DELETE().WHERE(ShoppingListItem.ID.EQ(Int(itemID)))
}
