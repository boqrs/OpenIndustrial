package tag

import "testing"

func TestNewTag(t *testing.T) {
	md := Metadata{
		Name:        "Configuration Flag",
		Description: "A boolean flag for system configuration",
		DataType:    "bool",
	}
	tg := NewTag("tag-1", "feature_enabled", true, md)

	if tg.GetID() != "tag-1" {
		t.Errorf("Expected ID 'tag-1', got '%s'", tg.GetID())
	}
	if tg.GetKind() != KindTag {
		t.Errorf("Expected Kind 'Tag', got '%s'", tg.GetKind())
	}
	if tg.Key != "feature_enabled" {
		t.Errorf("Expected Key 'feature_enabled', got '%s'", tg.Key)
	}
	if val, ok := tg.GetValue().(bool); !ok || val != true {
		t.Errorf("Expected Value 'true', got '%v'", tg.GetValue())
	}
}

func TestTagValue(t *testing.T) {
	md := Metadata{
		Name:     "Counter",
		DataType: "int",
	}
	tg := NewTag("tag-2", "event_count", 0, md)

	if val, ok := tg.GetValue().(int); !ok || val != 0 {
		t.Errorf("Expected initial Value '0', got '%v'", tg.GetValue())
	}

	tg.SetValue(100)
	if val, ok := tg.GetValue().(int); !ok || val != 100 {
		t.Errorf("Expected updated Value '100', got '%v'", tg.GetValue())
	}

	tg.SetValue("hello") // Test changing type, though DataType in Metadata suggests it shouldn't
	if val, ok := tg.GetValue().(string); !ok || val != "hello" {
		t.Errorf("Expected updated Value 'hello', got '%v'", tg.GetValue())
	}
}