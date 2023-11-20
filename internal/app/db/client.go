package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"simple-url-shortener/internal/app/config"
)

func NewDB(ctx context.Context, cfg *config.DatabaseCredentials) (*Database, error) {
	if err := MigrateDB(cfg); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	pool, err := pgxpool.Connect(ctx, generateDsn(cfg))
	if err != nil {
		return nil, err
	}
	return newDatabase(pool), nil
}

func generateDsn(cfg *config.DatabaseCredentials) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
}
