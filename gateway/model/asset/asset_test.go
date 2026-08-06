package asset

import "testing"

func TestNewAsset(t *testing.T) {
	md := Metadata{
		Name:        "Factory Floor 1",
		Description: "Main production area",
		Location:    "Building A",
	}
	ast := NewAsset("asset-1", md)

	if ast.GetID() != "asset-1" {
		t.Errorf("Expected ID 'asset-1', got '%s'", ast.GetID())
	}
	if ast.GetKind() != KindAsset {
		t.Errorf("Expected Kind 'Asset', got '%s'", ast.GetKind())
	}
	if ast.Metadata.Name != "Factory Floor 1" {
		t.Errorf("Expected Metadata Name 'Factory Floor 1', got '%s'", ast.Metadata.Name)
	}
}

func TestAssetLabels(t *testing.T) {
	md := Metadata{Name: "Factory Floor 1"}
	ast := NewAsset("asset-1", md)

	ast.SetLabel("owner", "Operations")
	ast.SetLabel("priority", "high")

	if v, ok := ast.Label("owner"); !ok || v != "Operations" {
		t.Errorf("Expected label 'owner' to be 'Operations', got '%s'", v)
	}
	if v, ok := ast.Label("priority"); !ok || v != "high" {
		t.Errorf("Expected label 'priority' to be 'high', got '%s'", v)
	}

	if _, ok := ast.Label("non-existent"); ok {
		t.Errorf("Expected non-existent label to not be found")
	}
}