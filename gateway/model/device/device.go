package device

import (
	"sync"

	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/object"
)

// KindDevice is the object kind for Device.
const (
	KindDevice object.Kind = "Device"
)

// Metadata describes static device information.
type Metadata struct {
	Name         string
	Model        string
	Manufacturer string
	Description  string
}

// Device is the unified industrial device model.
type Device struct {
	ID       string
	Metadata Metadata

	mu     sync.RWMutex
	labels map[string]string
}

// NewDevice creates a device.
func NewDevice(id string, md Metadata) *Device {
	return &Device{
		ID:       id,
		Metadata: md,
		labels:   make(map[string]string),
	}
}

// GetID returns the unique identifier of the device.
func (d *Device) GetID() string {
	return d.ID
}

// GetKind returns the kind of the object, which is "Device".
func (d *Device) GetKind() object.Kind {
	return KindDevice
}

// SetLabel sets a label for the device.
func (d *Device) SetLabel(k, v string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.labels[k] = v
}

// Label returns the value of a label for the given key.
func (d *Device) Label(k string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	v, ok := d.labels[k]
	return v, ok
}