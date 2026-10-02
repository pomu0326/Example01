// Package todo defines the Todo model and its persistence.
package todo

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// MaxTitleLength is the maximum number of characters allowed in a title.
const MaxTitleLength = 200

var (
	// ErrNotFound is returned when a todo with the given ID does not exist.
	ErrNotFound = errors.New("todo not found")
	// ErrInvalidTitle is returned when a title is empty or too long.
	ErrInvalidTitle = errors.New("title must be 1 to 200 characters")
)

// Todo is a single task.
type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Patch holds the fields to change in a partial update. Nil fields are left unchanged.
type Patch struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

// NormalizeTitle trims surrounding whitespace and validates the length.
func NormalizeTitle(title string) (string, error) {
	t := strings.TrimSpace(title)
	if n := utf8.RuneCountInString(t); n == 0 || n > MaxTitleLength {
		return "", ErrInvalidTitle
	}
	return t, nil
}
