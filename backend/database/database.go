package database

import (
	"context"
	"fmt"

	"github.com/berkkaradalan/stackflow/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, env *config.Env) (*pgxpool.Pool, error) {
	var dsn string

	if env.DbURL != "" {
		dsn = env.DbURL
	} else {
		fmt.Printf("postgres://%s:%s@%s:%s/%s?sslmode=disable",env.DbUser, env.DbPassword, env.DbHost, env.DbPort, env.DbName)
		// Local/manual config
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			env.DbUser, env.DbPassword, env.DbHost, env.DbPort, env.DbName,
		)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}