package eventbus

import (
	"sync"
	"testing"
	"time"
)

func TestNewEventBus(t *testing.T) {
	eb := NewEventBus()
	if eb == nil {
		t.Error("NewEventBus returned nil")
	}
}

func TestPublishAndSubscribe(t *testing.T) {
	eb := NewEventBus()
	topic := "test.topic"
	var wg sync.WaitGroup
	receivedCount := 0
	expectedData := "hello event"

	// Handler 1
	wg.Add(1)
	handler1 := func(event Event) {
		defer wg.Done()
		if event.Topic != topic {
			t.Errorf("Handler 1: Expected topic %s, got %s", topic, event.Topic)
		}
		if event.Data != expectedData {
			t.Errorf("Handler 1: Expected data %s, got %v", expectedData, event.Data)
		}
		receivedCount++
	}
	unsubscribe1 := eb.Subscribe(topic, handler1)

	// Handler 2
	wg.Add(1)
	handler2 := func(event Event) {
		defer wg.Done()
		if event.Topic != topic {
			t.Errorf("Handler 2: Expected topic %s, got %s", topic, event.Topic)
		}
		if event.Data != expectedData {
			t.Errorf("Handler 2: Expected data %s, got %v", expectedData, event.Data)
		}
		receivedCount++
	}
	eb.Subscribe(topic, handler2)

	// Publish event
	eb.Publish(Event{Topic: topic, Data: expectedData})

	// Wait for handlers to complete
	wg.Wait()

	if receivedCount != 2 {
		t.Errorf("Expected 2 events received, got %d", receivedCount)
	}

	// Test unsubscribe
	unsubscribe1()
	receivedCount = 0 // Reset count for next publish
	wg.Add(1)         // Only handler2 should receive now
	eb.Publish(Event{Topic: topic, Data: expectedData})
	wg.Wait() // Wait for handler2

	if receivedCount != 1 {
		t.Errorf("Expected 1 event received after unsubscribe, got %d", receivedCount)
	}
}

func TestNoSubscribers(t *testing.T) {
	eb := NewEventBus()
	topic := "no.subscribers"
	// Publishing to a topic with no subscribers should not cause a panic
	eb.Publish(Event{Topic: topic, Data: "some data"})
}

func TestMultipleTopics(t *testing.T) {
	eb := NewEventBus()
	topic1 := "topic.one"
	topic2 := "topic.two"
	var wg sync.WaitGroup
	received1 := 0
	received2 := 0

	wg.Add(1)
	eb.Subscribe(topic1, func(event Event) {
		defer wg.Done()
		received1++
	})

	wg.Add(1)
	eb.Subscribe(topic2, func(event Event) {
		defer wg.Done()
		received2++
	})

	eb.Publish(Event{Topic: topic1, Data: "data1"})
	eb.Publish(Event{Topic: topic2, Data: "data2"})

	wg.Wait()

	if received1 != 1 {
		t.Errorf("Expected 1 event for topic1, got %d", received1)
	}
	if received2 != 1 {
		t.Errorf("Expected 1 event for topic2, got %d", received2)
	}
}

func TestConcurrency(t *testing.T) {
	eb := NewEventBus()
	topic := "concurrent.topic"
	var publishWg sync.WaitGroup
	var subscribeWg sync.WaitGroup
	var receivedCount int32
	numPublishers := 10
	numSubscribers := 10
	numEventsPerPublisher := 100

	// Concurrently subscribe
	for i := 0; i < numSubscribers; i++ {
		subscribeWg.Add(1)
		go func() {
			defer subscribeWg.Done()
			eb.Subscribe(topic, func(event Event) {
				// Simulate some work
				time.Sleep(time.Millisecond)
				receivedCount++
			})
		}()
	}
	subscribeWg.Wait() // Wait for all subscribers to be ready

	// Concurrently publish
	for i := 0; i < numPublishers; i++ {
		publishWg.Add(1)
		go func() {
			defer publishWg.Done()
			for j := 0; j < numEventsPerPublisher; j++ {
				eb.Publish(Event{Topic: topic, Data: "concurrent data"})
			}
		}()
	}
	publishWg.Wait()

	// Give some time for all goroutines to process events
	time.Sleep(time.Second * 2)

	expectedTotalReceived := int32(numPublishers * numEventsPerPublisher * numSubscribers)
	if receivedCount != expectedTotalReceived {
		t.Errorf("Expected %d total events received, got %d", expectedTotalReceived, receivedCount)
	}
}