package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

const (
	DefaultLimit = 30  // DefaultLimit is the default page size if none is specified
	MaxLimit     = 100 // MaxLimit is the maximum allowed page size
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
	LastID   *int64  `json:"id"`
	LastName *string `json:"name"`
}

// EncodeCursor creates a base64-encoded cursor from the last item ID
func EncodeCursor(cursor Cursor) string {
	jsonBytes, err := json.Marshal(cursor)
	if err != nil {
		// This should never happen with simple struct
		return ""
	}
	return base64.StdEncoding.EncodeToString(jsonBytes)
}

// DecodeCursor decodes a base64-encoded cursor to extract the last ID
// Returns nil if cursor is empty or invalid
func DecodeCursor(cursor string) Cursor {
	if cursor == "" {
		return Cursor{}
	}

	jsonBytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return Cursor{}
	}

	var data Cursor
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return Cursor{}
	}

	return data
}

// NormalizeLimit ensures the limit is within acceptable bounds
func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}

// NewResult creates a paginated result from a slice of items
// It expects items to contain limit+1 items if there are more pages
func NewResult[T any](items []T, limit int, getCursor func(T) Cursor) Result[T] {
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
		cursor := DecodeCursor(cursorBase64)
		if cursor.LastID == nil {
			return Page{}, fmt.Errorf("invalid cursor format")
		}
	}

	return Page{
		Cursor: cursorBase64,
		Limit:  normalizedLimit,
	}, nil
}
