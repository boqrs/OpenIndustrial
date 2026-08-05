package tag

import "sync"

// Kind represents the type of an Object.
type Kind string

const (
	KindTag Kind = "Tag"
)

// Object is the common interface for all entities managed by the Runtime.
// This is a placeholder and will be replaced by the actual runtime/object.Object interface.
type Object interface {
	ID() string
	Kind() Kind
}

// Metadata describes static tag information.
type Metadata struct {
	Name        string
	Description string
	DataType    string // e.g., "string", "int", "bool"
}

// Tag represents a generic key-value pair with metadata.
type Tag struct {
	ID       string
	Key      string
	Value    any // The actual value of the tag
	Metadata Metadata

	mu sync.RWMutex
}

// NewTag creates a new Tag.
func NewTag(id, key string, value any, md Metadata) *Tag {
	return &Tag{
		ID:       id,
		Key:      key,
		Value:    value,
		Metadata: md,
	}
}

// ID returns the unique identifier of the tag.
func (t *Tag) GetID() string {
	return t.ID
}

// Kind returns the kind of the object, which is "Tag".
func (t *Tag) GetKind() Kind {
	return KindTag
}

// GetValue returns the current value of the tag.
func (t *Tag) GetValue() any {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.Value
}

// SetValue sets the value of the tag.
func (t *Tag) SetValue(value any) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Value = value
}