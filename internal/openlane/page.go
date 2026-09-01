package openlane

const (
	DefaultPageSize = 20
	MaxPageSize     = 50
)

// Page is the compact pagination envelope returned by list tools.
type Page[T any] struct {
	Items      []T     `json:"items"`
	NextCursor *string `json:"next_cursor"`
	HasMore    bool    `json:"has_more"`
	TotalCount int64   `json:"total_count"`
}

// ClampLimit bounds a requested page size.
func ClampLimit(n int) int64 {
	if n <= 0 {
		return DefaultPageSize
	}
	if n > MaxPageSize {
		return MaxPageSize
	}
	return int64(n)
}

// CursorPtr returns nil for an empty cursor.
func CursorPtr(cursor string) *string {
	if cursor == "" {
		return nil
	}
	return &cursor
}
