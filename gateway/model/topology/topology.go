package topology

import (
	"sync"

	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/object"
)

// KindTopology is the object kind for Topology.
const (
	KindTopology object.Kind = "Topology"
)

// Metadata describes static topology information.
type Metadata struct {
	Name        string
	Description string
}

// Topology represents the relationships between industrial objects.
type Topology struct {
	ID       string
	Metadata Metadata

	mu     sync.RWMutex
	labels map[string]string

	// TODO: Add fields for relationships (e.g., parent-child, peer-to-peer)
}

// NewTopology creates a new Topology instance.
func NewTopology(id string, md Metadata) *Topology {
	return &Topology{
		ID:       id,
		Metadata: md,
		labels:   make(map[string]string),
	}
}

// GetID returns the unique identifier of the topology.
func (t *Topology) GetID() string {
	return t.ID
}

// GetKind returns the kind of the object, which is "Topology".
func (t *Topology) GetKind() object.Kind {
	return KindTopology
}

// SetLabel sets a label for the topology.
func (t *Topology) SetLabel(k, v string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.labels[k] = v
}

// Label returns the value of a label for the given key.
func (t *Topology) Label(k string) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	v, ok := t.labels[k]
	return v, ok
}