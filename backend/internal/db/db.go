package db

import (
	"context"
	"fmt"

	"asset-flow/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(
	cfg *config.Config,
) (*DB, error) {

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pool, err := pgxpool.New(
		context.Background(),
		connStr,
	)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(
		context.Background(),
	); err != nil {
		pool.Close()
		return nil, err
	}

	return &DB{
		Pool: pool,
	}, nil
}
