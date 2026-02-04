package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"github.com/berkkaradalan/stackflow/config"
)

func Migrate(ctx context.Context, pool *pgxpool.Pool, cfg *config.Config) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			avatar_url VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,		
	}

	for i, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("migration failed at step %d: %w", i+1, err)
		}
	}

	if err := createDefaultAdmin(ctx, pool, cfg); err != nil {
		return fmt.Errorf("failed to create default admin: %w", err)
	}

	return nil
}

func createDefaultAdmin(ctx context.Context, pool *pgxpool.Pool, cfg *config.Config) error {
	var exists bool
	err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE role = 'admin')`).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.Env.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO users (username, email, password_hash, avatar_url, role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, cfg.Env.AdminUsername, cfg.Env.AdminEmail, string(hash), "https://github.com/shadcn.png", "admin", true)

	return err
}