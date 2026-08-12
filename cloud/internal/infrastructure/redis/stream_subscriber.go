package redis

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/OpenIndustrial/cloud/internal/event"
	"github.com/go-redis/redis/v8"
)

// StreamSubscriber listens to a Redis Stream for new events and dispatches them to a bus.
type StreamSubscriber struct {
	client    *redis.Client
	bus       event.Bus
	stream    string
	group     string
	consumer  string
	startFrom string
}

// NewStreamSubscriber creates a new subscriber instance.
// - client: The Redis client.
// - bus: The event bus to dispatch events to.
// - stream: The name of the Redis Stream to listen to (e.g., "openindustrial:events").
// - group: The consumer group name for this subscriber.
// - consumer: A unique name for this specific consumer instance within the group.
func NewStreamSubscriber(client *redis.Client, bus event.Bus, stream, group, consumer string) *StreamSubscriber {
	return &StreamSubscriber{
		client:    client,
		bus:       bus,
		stream:    stream,
		group:     group,
		consumer:  consumer,
		startFrom: "0", // Default to start from the beginning of the stream for a new group.
	}
}

// Start begins the event listening loop. This is a blocking call.
// It should be run in a separate goroutine.
func (s *StreamSubscriber) Start(ctx context.Context) {
	s.createGroupIfNeeded(ctx)

	log.Printf("INFO: Starting Redis Stream consumer loop for group '%s' on stream '%s'", s.group, s.stream)

	for {
		select {
		case <-ctx.Done():
			log.Println("INFO: Shutting down stream subscriber...")
			return
		default:
			s.processMessages(ctx)
		}
	}
}

func (s *StreamSubscriber) createGroupIfNeeded(ctx context.Context) {
	// XGROUP CREATE creates the group if it doesn't exist.
	// If it already exists, it returns an error, which we can safely ignore.
	// The 'MKSTREAM' option creates the stream itself if it's missing.
	err := s.client.XGroupCreateMkStream(ctx, s.stream, s.group, s.startFrom).Err()
	if err != nil && !errors.Is(err, redis.Nil) && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		// We log a warning but don't fail, as the group might already exist.
		log.Printf("WARN: Could not create consumer group (might already exist): %v", err)
	}
}

func (s *StreamSubscriber) processMessages(ctx context.Context) {
	// Block for up to 2 seconds waiting for a new message.
	// '>' means we want messages that have never been delivered to any other consumer in this group.
	streams, err := s.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    s.group,
		Consumer: s.consumer,
		Streams:  []string{s.stream, ">"},
		Count:    1, // Process one message at a time for simplicity.
		Block:    2 * time.Second,
	}).Result()

	if err != nil {
		// Timeouts are normal, just continue the loop.
		if errors.Is(err, redis.Nil) {
			return
		}
		log.Printf("ERROR: Failed to read from stream: %v", err)
		time.Sleep(1 * time.Second) // Avoid fast-spinning on persistent errors.
		return
	}

	for _, stream := range streams {
		for _, msg := range stream.Messages {
			s.handleMessage(ctx, &msg)
		}
	}
}

func (s *StreamSubscriber) handleMessage(ctx context.Context, msg *redis.XMessage) {
	log.Printf("INFO: Received message %s from stream '%s'", msg.ID, s.stream)

	// The event envelope is stored in the "event" field.
	eventJSON, ok := msg.Values["event"].(string)
	if !ok {
		log.Printf("ERROR: Message %s has no 'event' field. Skipping.", msg.ID)
		// We should probably ACK this malformed message to avoid processing it again.
		s.client.XAck(ctx, s.stream, s.group, msg.ID)
		return
	}

	var env event.Envelope
	if err := json.Unmarshal([]byte(eventJSON), &env); err != nil {
		log.Printf("ERROR: Failed to unmarshal event from message %s: %v", msg.ID, err)
		s.client.XAck(ctx, s.stream, s.group, msg.ID) // Malformed JSON, no point in retrying.
		return
	}

	// Dispatch the event to the bus.
	if err := s.bus.Dispatch(ctx, &env); err != nil {
		log.Printf("ERROR: Failed to dispatch event %s (type: %s): %v. Message will be retried later.", env.ID, env.Type, err)
		// IMPORTANT: We DO NOT ACK the message here.
		// This means the message will remain pending and can be re-processed later
		// by this or another consumer in the same group.
		return
	}

	// If dispatch was successful, ACK the message.
	if err := s.client.XAck(ctx, s.stream, s.group, msg.ID).Err(); err != nil {
		log.Printf("ERROR: Failed to ACK message %s: %v", msg.ID, err)
	} else {
		log.Printf("INFO: Successfully processed and ACKed message %s", msg.ID)
	}
}