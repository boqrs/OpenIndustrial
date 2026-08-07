package auth

import "time"

// Session represents a user's login session on a specific client.
type Session struct {
	ID        string
	UserID    string
	TokenHash string // Hash of the refresh token
	Client    string // e.g., "web-browser", "mobile-app"
	ExpiresAt time.Time
	CreatedAt time.Time
}