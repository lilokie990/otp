package utils

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lilokie/otp-auth/config"
)

// SetupDatabase sets up the database connection
func SetupDatabase(config *config.Config) (*sqlx.DB, error) {
	// Get connection string from config
	dsn := config.GetDSN()

	// Connect to database
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}
