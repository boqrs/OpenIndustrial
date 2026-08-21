package notification

import (
	"context"
	"encoding/json"
	"log"

	"github.com/boqrs/OpenIndustrial/cloud/internal/services/event"
)

// UserCreatedHandler is responsible for handling the "user created" event.
// For example, by sending a welcome email.
type UserCreatedHandler struct {
	// In a real application, this might have dependencies like an email client.
	// emailSender *some_email_service.Client
}

// NewUserCreatedHandler creates a new instance of UserCreatedHandler.
func NewUserCreatedHandler() *UserCreatedHandler {
	return &UserCreatedHandler{}
}

// Handle is the method that gets called by the event bus when a relevant event occurs.
func (h *UserCreatedHandler) Handle(ctx context.Context, evt *event.Envelope) error {
	// 1. Log that we've received the event.
	log.Printf("INFO: [Notification] Handling event '%s' (ID: %s)", evt.Type, evt.ID)

	// 2. Unmarshal the generic payload into our specific, expected struct.
	var payload event.UserCreatedPayload
	if err := json.Unmarshal(evt.Payload, &payload); err != nil {
		// If unmarshalling fails, it's a permanent error for this message.
		// We log it and return nil to ACK the message, preventing endless retries.
		log.Printf("ERROR: [Notification] failed to unmarshal UserCreatedPayload for event %s: %v", evt.ID, err)
		return nil
	}

	// 3. Execute the business logic.
	// Here, we simulate sending a welcome email.
	log.Printf("  📧 Sending welcome email to: %s (UserID: %s, TenantID: %s)\n", payload.Email, payload.UserID)

	// In a real-world scenario, you would call an email service here.
	// err := h.emailSender.SendWelcomeEmail(ctx, payload.Email, payload.Name)
	// if err != nil {
	//   log.Printf("ERROR: Failed to send welcome email: %v", err)
	//   // Returning an error will cause the event bus to NOT acknowledge the message,
	//   // allowing it to be retried later.
	//   return err
	// }

	// 4. Return nil to signify that the event has been successfully handled.
	log.Printf("INFO: [Notification] Successfully handled event %s.", evt.ID)
	return nil
}