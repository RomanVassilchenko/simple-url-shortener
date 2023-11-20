package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"simple-url-shortener/internal/app/config"
)

// MigrateDB runs database migrations.
func MigrateDB(cfg *config.DatabaseCredentials) error {
	// Register the pq driver
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to register pq driver: %w", err)
	}

	// Explicitly register the pq driver for the sql package
	if _, err := sql.Open("postgres", ""); err != nil {
		return fmt.Errorf("failed to open sql connection: %w", err)
	}

	migrateDB, err := sql.Open("postgres", generateDsn(cfg))
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer migrateDB.Close()

	migrationsDir := "internal/app/db/migrations"

	if err := goose.Up(migrateDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
