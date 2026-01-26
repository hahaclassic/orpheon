package config

import (
	"time"

	"github.com/hahaclassic/orpheon/pkg/config"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/postgres"
)

type Config struct {
	Postgres PostgresConfig
	GRPC     GRPCServerConfig
	HTTP     HTTPServerConfig
}

func (c *Config) Loaders() []config.InternalConfigLoader {
	return []config.InternalConfigLoader{
		&c.GRPC,
		&c.HTTP,
		&c.Postgres,
	}
}

type GRPCServerConfig struct {
	Host string `env:"GRPC_HOST" env-required:"true"`
	Port string `env:"GRPC_PORT" env-required:"true"`
}

func (cfg *GRPCServerConfig) Load() {
	cfg.Host = config.GetEnv("GRPC_HOST", "localhost")
	cfg.Port = config.GetEnv("GRPC_PORT", "50051")
}

type HTTPServerConfig struct {
	Host     string `env:"HTTP_HOST" env-required:"true"`
	Port     string `env:"HTTP_PORT" env-required:"true"`
	ReadOnly bool   `env:"READ_ONLY" env-required:"true"`
}

func (cfg *HTTPServerConfig) Load() {
	cfg.Host = config.GetEnv("HTTP_HOST", "localhost")
	cfg.Port = config.GetEnv("HTTP_PORT", "8081")
	cfg.ReadOnly = config.GetEnvAsBool("READ_ONLY", false)
}

type PostgresConfig postgres.PostgresConfig

func (cfg *PostgresConfig) Load() {
	cfg.Host = config.GetEnv("POSTGRES_HOST", "localhost")
	cfg.Port = config.GetEnv("POSTGRES_PORT", "5432")
	cfg.User = config.GetEnv("POSTGRES_USER", "user")
	cfg.Password = config.GetEnv("POSTGRES_PASSWORD", "password")
	cfg.DBName = config.GetEnv("POSTGRES_DB", "orpheon")
	cfg.SSLMode = config.GetEnv("POSTGRES_SSL_MODE", "disable")
	cfg.StartTimeout = config.GetEnvAsDuration("POSTGRES_START_TIMEOUT", 5*time.Second)
}
