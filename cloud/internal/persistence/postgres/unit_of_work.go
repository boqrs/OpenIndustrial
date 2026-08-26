package postgres

import (
	"github.com/boqrs/nexus/database"
	//"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/postgres"
	"context"

	"gorm.io/gorm"
)

// transactionKey is an unexported type used as a key for storing the transaction in the context.
type transactionKey struct{}

// withTransaction injects the GORM transaction into the context.
func withTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, transactionKey{}, tx)
}

// transactionFromContext extracts the GORM transaction from the context, if one exists.
func transactionFromContext(ctx context.Context) *gorm.DB {
	tx, _ := ctx.Value(transactionKey{}).(*gorm.DB)
	return tx
}

// dbFromContext is the core function for repositories. It intelligently decides
// whether to use the transaction from the context or the default database connection.
func dbFromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx := transactionFromContext(ctx); tx != nil {
		// A transaction is in progress, use it.
		return tx
	}
	// No transaction, use the default DB connection but ensure it carries the context.
	return defaultDB.WithContext(ctx)
}

// UnitOfWork defines an interface for atomic operations that span multiple repositories.
type UnitOfWork interface {
	// Execute runs the given function in a single transaction.
	// If the function returns an error, the transaction is rolled back.
	// Otherwise, the transaction is committed.
	Execute(ctx context.Context, fn func(ctx context.Context) error) error
}

// unitOfWork is the Postgres implementation of the UnitOfWork interface.
type unitOfWork struct {
	dbProv *database.DBProvider
}

// NewUnitOfWork creates a new postgres-based UnitOfWork.
func NewUnitOfWork(dbProv *database.DBProvider) UnitOfWork {
	return &unitOfWork{
		dbProv: dbProv,
	}
}

// Execute runs the given function within a GORM transaction.
func (u *unitOfWork) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	// Get the base DB connection and ensure it carries the initial context.
	db := u.dbProv.Get().WithContext(ctx)

	// Use GORM's built-in transaction functionality. It handles all the boilerplate.
	return db.Transaction(func(tx *gorm.DB) error {
		// Inject the transaction into a new context.
		txCtx := withTransaction(ctx, tx)
		// Execute the application logic with the new transactional context.
		return fn(txCtx)
	})
}