package domain

type ShoppingList struct {
	ID     int64
	UserID int64
	Name   string
	Items  []ShoppingListItem
}

type ShoppingListItem struct {
	ID         int64
	Ingredient string
	Quantity   *string
	Unit       *string
	Done       bool
	SortOrder  int64
}
