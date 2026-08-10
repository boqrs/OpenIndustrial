package identity

import "context"

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