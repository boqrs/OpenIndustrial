package identity

import (
	"context"
	"time"
		"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserType defines the type of a user.
type UserType string

const (
	UserTypeAdmin  UserType = "admin"
	UserTypeMember UserType = "member"
)



// HashPassword creates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plain-text password with a bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Claims defines the structure of the JWT claims.
type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT for a user.
func GenerateToken(userID, tenantID uuid.UUID, jwtSecret string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// Transaction defines an interface for database transactions.
// This allows the service layer to be independent of the specific db driver.
type Transaction interface {
	Commit() error
	Rollback() error
}

// TransactionalRepository defines an interface for repositories that can begin transactions.
type TransactionalRepository interface {
	BeginTx(ctx context.Context) (Transaction, error)
}

type ListUsersRepoParams struct {
	Limit  int
	Offset int
}
