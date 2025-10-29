package postgres

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresConfig struct {
	Host         string        `env:"POSTGRES_HOST" env-required:"true"`
	Port         string        `env:"POSTGRES_PORT" env-required:"true"`
	User         string        `env:"POSTGRES_USER" env-required:"true"`
	Password     string        `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName       string        `env:"POSTGRES_DBNAME" env-required:"true"`
	SSLMode      string        `env:"POSTGRES_SSLMODE" env-default:"disable"`
	StartTimeout time.Duration `env:"POSTGRES_START_TIMEOUT" env-default:"5s"`
}

func NewPostgresPool(ctx context.Context, cfg PostgresConfig) *pgxpool.Pool {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, cfg.StartTimeout)
	defer cancel()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	pool, err := pgxpool.New(ctxWithTimeout, dsn)
	if err != nil {
		log.Fatalf("unable to create pgx pool: %v", err)
	}

	if err = pool.Ping(ctxWithTimeout); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	return pool
}
