package safety

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewReadOnlyPool creates a connection pool where every session is read-only.

func NewReadOnlyPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse connection string: %w", err)
	}

	// This setting is enforced at the Postgres engine level.
	config.ConnConfig.RuntimeParams["default_transaction_read_only"] = "on"

	// Keep the pool small, the Copilot is interactive, not high-throughput
	config.MaxConns = 5
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("connection test failed: %w", err)
	}

	return pool, nil
}