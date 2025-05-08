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
	Port         int           `env:"POSTGRES_PORT" env-required:"true"`
	User         string        `env:"POSTGRES_USER" env-required:"true"`
	Password     string        `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName       string        `env:"POSTGRES_DBNAME" env-required:"true"`
	SSLMode      string        `env:"POSTGRES_SSLMODE" env-default:"disable"`
	StartTimeout time.Duration `env:"POSTGRES_START_TIMEOUT" env-default:"5s"`
}

func (cfg PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
}

func NewPostgresPool(cfg PostgresConfig) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.StartTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("unable to create pgx pool: %v", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	return pool
}
