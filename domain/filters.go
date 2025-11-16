package domain

// IngredientFilters contains filtering and sorting parameters for fetching ingredients
type IngredientFilters struct {
	Page      Page
	SortBy    string // "name" or "id"
	SortOrder string // "asc" or "desc"
	Search    string // optional search filter
}

// SearchFilters contains basic filtering parameters for simple paginated searches
type SearchFilters struct {
	Page   Page
	Search string // optional search filter
}

// Type aliases for clarity and future extensibility
type UnitFilters = SearchFilters
type TagFilters = SearchFilters
