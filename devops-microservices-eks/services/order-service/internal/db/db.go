package db

import (
	"context"
	"time"

	"github.com/devopsdemo/order-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect builds a connection pool and pings Postgres immediately — this is
// what makes startup fail fast if the DB is unreachable, instead of only
// discovering it on the first incoming request.
func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
