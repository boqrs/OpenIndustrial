package point

import (
	"testing"
	"time"
)

func TestNewPoint(t *testing.T) {
	def := PointDefinition{
		ID:          "1",
		Key:         "temperature",
		Name:        "Temperature Sensor",
		Description: "Measures ambient temperature",
		Unit:        "°C",
		DataType:    TypeFloat64,
		Readable:    true,
		Writable:    false,
	}

	p := NewPoint(def)

	if p.GetID() != "1" {
		t.Errorf("Expected ID '1', got '%s'", p.GetID())
	}
	if p.GetKind() != KindPoint {
		t.Errorf("Expected Kind 'Point', got '%s'", p.GetKind())
	}
	if p.Definition.Key != "temperature" {
		t.Errorf("Expected Key 'temperature', got '%s'", p.Definition.Key)
	}
	if p.State.Quality != QualityBad {
		t.Errorf("Expected initial Quality 'QualityBad', got '%v'", p.State.Quality)
	}
	if p.State.Value != nil {
		t.Errorf("Expected initial Value 'nil', got '%v'", p.State.Value)
	}
}

func TestUpdateState(t *testing.T) {
	def := PointDefinition{
		ID:       "1",
		Key:      "temperature",
		DataType: TypeFloat64,
	}
	p := NewPoint(def)

	testValue := 25.5
	testQuality := QualityGood
	testTimestamp := time.Now().Truncate(time.Millisecond) // Truncate to avoid nanosecond differences

	p.UpdateState(testValue, testQuality, testTimestamp)

	if val, ok := p.State.Value.(float64); !ok || val != testValue {
		t.Errorf("Expected Value '%v', got '%v'", testValue, p.State.Value)
	}
	if p.State.Quality != testQuality {
		t.Errorf("Expected Quality '%v', got '%v'", testQuality, p.State.Quality)
	}
	if !p.State.Timestamp.Equal(testTimestamp) {
		t.Errorf("Expected Timestamp '%v', got '%v'", testTimestamp, p.State.Timestamp)
	}

	// Test updating with different type
	p.Definition.DataType = TypeString
	testStringValue := "active"
	p.UpdateState(testStringValue, QualityGood, time.Now())
	if val, ok := p.State.Value.(string); !ok || val != testStringValue {
		t.Errorf("Expected string Value '%v', got '%v'", testStringValue, p.State.Value)
	}
}