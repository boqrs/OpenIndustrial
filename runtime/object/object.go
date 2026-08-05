package object

// Kind represents the type of an Object.
type Kind string

// Object is the common interface for all entities managed by the Runtime.
// All domain models (Device, Point, Asset, Topology, Tag, etc.) must implement this interface.
type Object interface {
	// GetID returns the unique identifier of the object.
	GetID() string
	// GetKind returns the kind of the object (e.g., "Device", "Point", "Asset").
	GetKind() Kind
}