package registry

import (
	"fmt"
	"sync"

	"github.com/OpenGongChang/OpenIndustrial/runtime/object"
)

// Registry defines the interface for managing runtime objects.
type Registry interface {
	// Add adds an object to the registry. Returns an error if an object with the same ID already exists.
	Add(obj object.Object) error
	// Get retrieves an object by its ID. Returns the object and true if found, otherwise nil and false.
	Get(id string) (object.Object, bool)
	// Remove removes an object from the registry by its ID. Returns an error if the object does not exist.
	Remove(id string) error
	// List returns all objects currently in the registry.
	List() []object.Object
}

// NewRegistry creates a new in-memory Registry.
func NewRegistry() Registry {
	return &inMemoryRegistry{
		objects: make(map[string]object.Object),
	}
}

// inMemoryRegistry is a simple, concurrency-safe in-memory implementation of the Registry interface.
type inMemoryRegistry struct {
	mu      sync.RWMutex
	objects map[string]object.Object
}

// Add adds an object to the registry. Returns an error if an object with the same ID already exists.
func (r *inMemoryRegistry) Add(obj object.Object) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := obj.GetID()
	if _, exists := r.objects[id]; exists {
		return fmt.Errorf("object with ID '%s' already exists in registry", id)
	}
	r.objects[id] = obj
	return nil
}

// Get retrieves an object by its ID. Returns the object and true if found, otherwise nil and false.
func (r *inMemoryRegistry) Get(id string) (object.Object, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	obj, ok := r.objects[id]
	return obj, ok
}

// Remove removes an object from the registry by its ID. Returns an error if the object does not exist.
func (r *inMemoryRegistry) Remove(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.objects[id]; !exists {
		return fmt.Errorf("object with ID '%s' not found in registry", id)
	}
	delete(r.objects, id)
	return nil
}

// List returns all objects currently in the registry.
func (r *inMemoryRegistry) List() []object.Object {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]object.Object, 0, len(r.objects))
	for _, obj := range r.objects {
		list = append(list, obj)
	}
	return list
}