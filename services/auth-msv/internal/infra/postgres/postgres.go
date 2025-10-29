package postgres

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/hahaclassic/orpheon/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig = config.PostgresConfig

func NewPostgresPool(cfg PostgresConfig) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.StartTimeout)
	defer cancel()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.DB,
		cfg.SSLMode,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to create pgx pool: %v", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	return pool
}
