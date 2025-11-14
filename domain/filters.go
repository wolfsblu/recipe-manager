package domain

// IngredientFilters contains filtering and sorting parameters for fetching ingredients
type IngredientFilters struct {
	Page      Page
	SortBy    string // "name" or "id"
	SortOrder string // "asc" or "desc"
	Search    string // optional search filter
}
