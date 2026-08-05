package registry

import (
	"fmt"
	"sync"
	"testing"

	"github.com/OpenGongChang/OpenIndustrial/runtime/object"
)

// MockObject implements the object.Object interface for testing purposes.
type MockObject struct {
	IDValue   string
	KindValue object.Kind
}

func (m *MockObject) GetID() string {
	return m.IDValue
}

func (m *MockObject) GetKind() object.Kind {
	return m.KindValue
}

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Error("NewRegistry returned nil")
	}
}

func TestAddAndGet(t *testing.T) {
	reg := NewRegistry()
	obj1 := &MockObject{IDValue: "obj1", KindValue: "TestKind"}
	obj2 := &MockObject{IDValue: "obj2", KindValue: "TestKind"}

	// Test Add
	err := reg.Add(obj1)
	if err != nil {
		t.Fatalf("Add obj1 failed: %v", err)
	}
	err = reg.Add(obj2)
	if err != nil {
		t.Fatalf("Add obj2 failed: %v", err)
	}

	// Test Get existing
	retrievedObj, found := reg.Get("obj1")
	if !found || retrievedObj.GetID() != "obj1" {
		t.Errorf("Get obj1 failed: found=%t, ID=%s", found, retrievedObj.GetID())
	}

	// Test Get non-existing
	_, found = reg.Get("obj3")
	if found {
		t.Error("Get obj3 unexpectedly found")
	}

	// Test Add duplicate
	err = reg.Add(obj1)
	if err == nil {
		t.Error("Add duplicate obj1 unexpectedly succeeded")
	}
	expectedErr := "object with ID 'obj1' already exists in registry"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestRemove(t *testing.T) {
	reg := NewRegistry()
	obj1 := &MockObject{IDValue: "obj1", KindValue: "TestKind"}
	reg.Add(obj1)

	// Test Remove existing
	err := reg.Remove("obj1")
	if err != nil {
		t.Fatalf("Remove obj1 failed: %v", err)
	}
	_, found := reg.Get("obj1")
	if found {
		t.Error("obj1 found after removal")
	}

	// Test Remove non-existing
	err = reg.Remove("obj2")
	if err == nil {
		t.Error("Remove non-existing obj2 unexpectedly succeeded")
	}
	expectedErr := "object with ID 'obj2' not found in registry"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestList(t *testing.T) {
	reg := NewRegistry()
	obj1 := &MockObject{IDValue: "obj1", KindValue: "TestKind"}
	obj2 := &MockObject{IDValue: "obj2", KindValue: "TestKind"}

	reg.Add(obj1)
	reg.Add(obj2)

	list := reg.List()
	if len(list) != 2 {
		t.Errorf("Expected list length 2, got %d", len(list))
	}

	foundObj1 := false
	foundObj2 := false
	for _, obj := range list {
		if obj.GetID() == "obj1" {
			foundObj1 = true
		}
		if obj.GetID() == "obj2" {
			foundObj2 = true
		}
	}
	if !foundObj1 || !foundObj2 {
		t.Error("List did not contain all added objects")
	}

	// Test List after removal
	reg.Remove("obj1")
	list = reg.List()
	if len(list) != 1 {
		t.Errorf("Expected list length 1 after removal, got %d", len(list))
	}
	if list[0].GetID() != "obj2" {
		t.Errorf("Expected obj2 in list, got %s", list[0].GetID())
	}
}

func TestConcurrency(t *testing.T) {
	reg := NewRegistry()
	var wg sync.WaitGroup
	numGoroutines := 100
	numAddsPerGoroutine := 10

	// Concurrently add objects
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for j := 0; j < numAddsPerGoroutine; j++ {
				id := fmt.Sprintf("obj-%d-%d", g, j)
				obj := &MockObject{IDValue: id, KindValue: "ConcurrentKind"}
				err := reg.Add(obj)
				if err != nil {
					// Only expect error if ID already exists, which shouldn't happen with unique IDs
					t.Errorf("Goroutine %d: Add %s failed: %v", g, id, err)
				}
			}
		}(i)
	}
	wg.Wait()

	expectedTotal := numGoroutines * numAddsPerGoroutine
	if len(reg.List()) != expectedTotal {
		t.Errorf("Expected %d objects after concurrent adds, got %d", expectedTotal, len(reg.List()))
	}

	// Concurrently get and remove objects
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for j := 0; j < numAddsPerGoroutine; j++ {
				id := fmt.Sprintf("obj-%d-%d", g, j)
				if j%2 == 0 { // Remove half
					reg.Remove(id)
				} else { // Get the other half
					_, found := reg.Get(id)
					if !found {
						t.Errorf("Goroutine %d: Get %s failed, object not found", g, id)
					}
				}
			}
		}(i)
	}
	wg.Wait()

	// Check remaining objects
	remaining := reg.List()
	if len(remaining) != expectedTotal/2 {
		t.Errorf("Expected %d objects after concurrent removes, got %d", expectedTotal/2, len(remaining))
	}
}