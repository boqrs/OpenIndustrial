package asset

import (
	"sync"

	"github.com/OpenGongChang/OpenIndustrial/runtime/object"
)

// KindAsset is the object kind for Asset.
const (
	KindAsset object.Kind = "Asset"
)

// Metadata describes static asset information.
type Metadata struct {
	Name        string
	Description string
	Location    string
}

// Asset is a logical grouping of industrial objects.
type Asset struct {
	ID       string
	Metadata Metadata

	mu     sync.RWMutex
	labels map[string]string
}

// NewAsset creates a new Asset instance.
func NewAsset(id string, md Metadata) *Asset {
	return &Asset{
		ID:       id,
		Metadata: md,
		labels:   make(map[string]string),
	}
}

// GetID returns the unique identifier of the asset.
func (a *Asset) GetID() string {
	return a.ID
}

// GetKind returns the kind of the object, which is "Asset".
func (a *Asset) GetKind() object.Kind {
	return KindAsset
}

// SetLabel sets a label for the asset.
func (a *Asset) SetLabel(k, v string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.labels[k] = v
}

// Label returns the value of a label for the given key.
func (a *Asset) Label(k string) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	v, ok := a.labels[k]
	return v, ok
}