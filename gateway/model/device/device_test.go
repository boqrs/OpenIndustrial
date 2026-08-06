package device

import (
	"testing"
)

func TestNewDevice(t *testing.T) {
	md := Metadata{
		Name:         "Test Device",
		Model:        "Model A",
		Manufacturer: "Acme Corp",
		Description:  "A test industrial device",
	}
	dev := NewDevice("dev-123", md)

	if dev.GetID() != "dev-123" {
		t.Errorf("Expected ID 'dev-123', got '%s'", dev.GetID())
	}
	if dev.GetKind() != KindDevice {
		t.Errorf("Expected Kind 'Device', got '%s'", dev.GetKind())
	}
	if dev.Metadata.Name != "Test Device" {
		t.Errorf("Expected Metadata Name 'Test Device', got '%s'", dev.Metadata.Name)
	}
}

func TestDeviceLabels(t *testing.T) {
	md := Metadata{Name: "Test Device"}
	dev := NewDevice("dev-123", md)

	dev.SetLabel("location", "factory-floor-1")
	dev.SetLabel("area", "production-line-A")

	if v, ok := dev.Label("location"); !ok || v != "factory-floor-1" {
		t.Errorf("Expected label 'location' to be 'factory-floor-1', got '%s'", v)
	}
	if v, ok := dev.Label("area"); !ok || v != "production-line-A" {
		t.Errorf("Expected label 'area' to be 'production-line-A', got '%s'", v)
	}

	if _, ok := dev.Label("non-existent"); ok {
		t.Errorf("Expected non-existent label to not be found")
	}
}