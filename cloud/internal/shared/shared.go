package shared

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// BaseEntity provides common fields for all domain entities.
// It can be embedded in other entities to provide ID, CreatedAt, and UpdatedAt fields.
type BaseEntity struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PageRequest defines the standard parameters for a paginated query.
type PageRequest struct {
	Page int `form:"page" json:"page"`
	Size int `form:"size" json:"size"`
}

// GetOffset calculates the offset for a database query based on the page number and size.
// It ensures that the page number is at least 1.
func (p *PageRequest) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return (p.Page - 1) * p.GetSize()
}

// GetSize returns the size for a database query.
// It provides a default size if the requested size is not specified.
func (p *PageRequest) GetSize() int {
	if p.Size <= 0 {
		p.Size = 10
	}
	return p.Size
}

// PageResult defines the standard structure for a paginated response.
// It is a generic type, allowing it to be used for any kind of entity list.
type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

// Common domain errors used across different modules.
var (
	ErrNotFound         = errors.New("resource not found")
	ErrPermissionDenied = errors.New("permission denied")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrConflict         = errors.New("resource conflict or already exists")
)