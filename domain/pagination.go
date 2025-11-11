package domain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

const (
	DefaultPageSize = 30  // DefaultPageSize is the default page size if none is specified
	MaxPageSize     = 100 // MaxPageSize is the maximum allowed page size
)

// Page contains pagination parameters
type Page struct {
	Cursor string
	Limit  int
}

// Result wraps paginated data with metadata
type Result[T any] struct {
	Data       []T
	NextCursor *string
	HasMore    bool
}

// Cursor represents the data encoded in a pagination cursor
type Cursor struct {
	LastID   int64     `json:"id"`
	LastName string    `json:"name"`
	LastDate time.Time `json:"date"`
}

type NameCursor struct {
	LastID   int64  `json:"id"`
	LastName string `json:"name"`
}

type DateCursor struct {
	LastID   int64     `json:"id"`
	LastDate time.Time `json:"date"`
}

type RecipeCursor struct {
	LastID        int64     `json:"id"`
	LastName      string    `json:"name,omitempty"`
	LastCreatedAt time.Time `json:"created_at,omitempty"`
	LastServings  int64     `json:"servings,omitempty"`
}

type SortField string

const (
	SortFieldName      SortField = "name"
	SortFieldCreatedAt SortField = "created_at"
	SortFieldServings  SortField = "servings"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

// RecipeSort combines sort field and order
type RecipeSort struct {
	Field SortField
	Order SortOrder
}

func NewDescendingDateCursor() *DateCursor {
	return &DateCursor{
		LastID:   math.MaxInt64,
		LastDate: time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
	}
}

func NewDescendingRecipeCursorByCreatedAt() *RecipeCursor {
	return &RecipeCursor{
		LastID:        math.MaxInt64,
		LastCreatedAt: time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC),
	}
}

func NewDescendingRecipeCursorByName() *RecipeCursor {
	return &RecipeCursor{
		LastID:   math.MaxInt64,
		LastName: string(rune(0x10FFFF)), // Max unicode character
	}
}

func NewDescendingRecipeCursorByServings() *RecipeCursor {
	return &RecipeCursor{
		LastID:       math.MaxInt64,
		LastServings: math.MaxInt64,
	}
}

// IsValid checks if a SortField value is valid
func (f SortField) IsValid() bool {
	switch f {
	case SortFieldName, SortFieldCreatedAt, SortFieldServings:
		return true
	default:
		return false
	}
}

// IsValid checks if a SortOrder value is valid
func (o SortOrder) IsValid() bool {
	switch o {
	case SortOrderAsc, SortOrderDesc:
		return true
	default:
		return false
	}
}

// NewRecipeSort creates a RecipeSort with defaults if invalid values provided
func NewRecipeSort(field SortField, order SortOrder) RecipeSort {
	// Default to created_at if invalid
	if !field.IsValid() {
		field = SortFieldCreatedAt
	}
	// Default to desc if invalid
	if !order.IsValid() {
		order = SortOrderDesc
	}
	return RecipeSort{Field: field, Order: order}
}

// EncodeCursor creates a base64-encoded cursor from the last item ID
func EncodeCursor(cursor any) string {
	jsonBytes, err := json.Marshal(cursor)
	if err != nil {
		// This should never happen with simple struct
		return ""
	}
	return base64.StdEncoding.EncodeToString(jsonBytes)
}

// DecodeCursor decodes a base64-encoded cursor to extract the last ID
// T should be a pointer type (*NameCursor, *DateCursor, etc.)
func DecodeCursor[T any](cursor string) (T, error) {
	var zero T
	if cursor == "" {
		return zero, fmt.Errorf("cursor is empty")
	}

	jsonBytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return zero, err
	}

	if err := json.Unmarshal(jsonBytes, &zero); err != nil {
		return zero, err
	}
	return zero, nil
}

// NormalizeLimit ensures the limit is within acceptable bounds
func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultPageSize
	}
	if limit > MaxPageSize {
		return MaxPageSize
	}
	return limit
}

// NewPagedResult creates a paginated result from a slice of items
// It expects items to contain limit+1 items if there are more pages
func NewPagedResult[T any, C any](items []T, limit int, getCursor func(T) C) Result[T] {
	hasMore := len(items) > limit
	data := items
	if hasMore {
		data = items[:limit]
	}

	var nextCursor *string
	if hasMore && len(data) > 0 {
		lastItem := data[len(data)-1]
		cursor := EncodeCursor(getCursor(lastItem))
		nextCursor = &cursor
	}

	return Result[T]{
		Data:       data,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}
}

// ValidatePage validates and normalizes pagination parameters
func ValidatePage(cursorBase64 string, limit int) (Page, error) {
	normalizedLimit := NormalizeLimit(limit)

	if cursorBase64 != "" {
		_, err := DecodeCursor[any](cursorBase64)
		if err != nil {
			return Page{}, WrapError(ErrPagination, err)
		}
	}

	return Page{
		Cursor: cursorBase64,
		Limit:  normalizedLimit,
	}, nil
}
