package postgres

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The PostgreSQL driver, blank import
)

// NewDB creates a new database connection pool.
func NewDB(databaseURL string) (*sqlx.DB, error) {
	// sqlx.Connect is a wrapper around sql.Open and db.Ping.
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// You can set connection pool settings here if needed.
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxLifetime(time.Hour)

	return db, nil
}