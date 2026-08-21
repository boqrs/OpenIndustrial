package event

// =============================================================================
// Event Registry
//
// This file acts as a centralized catalog for all events within the system.
// It defines event type constants and their corresponding payload structs.
// This approach provides type safety, code completion, and a single source of
// truth for all domain events.
//
// Naming Conventions:
//
// Event Type Constant:
//   <Domain><Aggregate><Action>
//   Example: IdentityUserCreated
//
// Event Type String:
//   "domain.aggregate.action"
//   Example: "identity.user.created"
//
// Payload Struct:
//   <Action>Payload
//   Example: UserCreatedPayload
//
// =============================================================================

const (
	// IdentityUserCreated is published when a new user is successfully created
	// in the identity domain.
	IdentityUserCreated = "identity.user.created"

	// IdentityTenantCreated is published when a new tenant (organization) is
	// successfully registered, along with its first admin user.
	IdentityTenantCreated = "identity.tenant.created"
)

// --- Event Payloads ---
// For each event type, a corresponding strongly-typed payload struct should be
// defined here. This enforces the contract between event producers and consumers.

// UserCreatedPayload is the data structure for the IdentityUserCreated event.
type UserCreatedPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// TenantCreatedPayload is the data structure for the IdentityTenantCreated event.
type TenantCreatedPayload struct {
	TenantID    string `json:"tenant_id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	AdminUserID string `json:"admin_user_id"`
}