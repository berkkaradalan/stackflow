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
		`CREATE TABLE IF NOT EXISTS invite_tokens (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) NOT NULL,
			username VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			token VARCHAR(255) UNIQUE NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			used_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_invite_tokens_token ON invite_tokens(token)`,
		`CREATE INDEX IF NOT EXISTS idx_invite_tokens_email ON invite_tokens(email)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			status VARCHAR(50) NOT NULL DEFAULT 'active',
			created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_created_by ON projects(created_by)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status)`,
		`CREATE TABLE IF NOT EXISTS agents (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			created_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(50) NOT NULL,
			level VARCHAR(20) NOT NULL,
			provider VARCHAR(50) NOT NULL,
			model VARCHAR(100) NOT NULL,
			api_key TEXT NOT NULL,
			config JSONB DEFAULT '{
				"temperature": 0.7,
				"max_tokens": 2000,
				"top_p": 1.0,
				"frequency_penalty": 0.0,
				"presence_penalty": 0.0
			}'::jsonb,
			status VARCHAR(50) NOT NULL DEFAULT 'idle',
			is_active BOOLEAN DEFAULT true,
			last_active_at TIMESTAMP,
			total_tokens_used BIGINT DEFAULT 0,
			total_cost DECIMAL(10, 4) DEFAULT 0.0,
			total_requests BIGINT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_agents_project_id ON agents(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_agents_created_by ON agents(created_by)`,
		`CREATE INDEX IF NOT EXISTS idx_agents_status ON agents(status)`,
		`CREATE INDEX IF NOT EXISTS idx_agents_role ON agents(role)`,
		`CREATE INDEX IF NOT EXISTS idx_agents_provider ON agents(provider)`,
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