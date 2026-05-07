package sink

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// creates a connection pool to the pgwatch sink.

func NewReadOnlyPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse sink connection string: %w", err)
	}

	// every transaction is read-only 
	config.ConnConfig.RuntimeParams["default_transaction_read_only"] = "on"

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create sink pool: %w", err)
	}

	// Test the connection immediately
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping sink: %w", err)
	}

	return pool, nil
}

// returns all sys_id values in the sink.
func GetSysIDs(ctx context.Context, pool *pgxpool.Pool) ([]string, error) {
	rows, err := pool.Query(ctx,
		"SELECT DISTINCT sys_id FROM metrics.pgwatch2_pgwatch2 ORDER BY sys_id",
	)
	if err != nil {
		return nil, fmt.Errorf("query sys_ids: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}