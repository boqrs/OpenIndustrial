package topology

import "testing"

func TestNewTopology(t *testing.T) {
	md := Metadata{
		Name:        "Plant Layout",
		Description: "Overall factory topology",
	}
	top := NewTopology("topo-1", md)

	if top.GetID() != "topo-1" {
		t.Errorf("Expected ID 'topo-1', got '%s'", top.GetID())
	}
	if top.GetKind() != KindTopology {
		t.Errorf("Expected Kind 'Topology', got '%s'", top.GetKind())
	}
	if top.Metadata.Name != "Plant Layout" {
		t.Errorf("Expected Metadata Name 'Plant Layout', got '%s'", top.Metadata.Name)
	}
}

func TestTopologyLabels(t *testing.T) {
	md := Metadata{Name: "Plant Layout"}
	top := NewTopology("topo-1", md)

	top.SetLabel("version", "1.0")
	top.SetLabel("status", "active")

	if v, ok := top.Label("version"); !ok || v != "1.0" {
		t.Errorf("Expected label 'version' to be '1.0', got '%s'", v)
	}
	if v, ok := top.Label("status"); !ok || v != "active" {
		t.Errorf("Expected label 'status' to be 'active', got '%s'", v)
	}

	if _, ok := top.Label("non-existent"); ok {
		t.Errorf("Expected non-existent label to not be found")
	}
}