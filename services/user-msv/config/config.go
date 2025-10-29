package config

import (
	"flag"
	"log"
	"time"

	"github.com/hahaclassic/orpheon/pkg/config"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/postgres"
	"github.com/joho/godotenv"
)

// Load environment variables in the following order:
// 1. If the -config flag is provided, load vars from the specified file.
// 2. If the -config flag is not set, load vars from the default ".env" file.
// 3. For any variables still unset, assign default values.

const defaultEnvPath = ".env"

type Config struct {
	Postgres   postgres.PostgresConfig
	GRPCServer GRPCServerConfig
}

type GRPCServerConfig struct {
	Host string `env:"GRPC_HOST" env-required:"true"`
	Port string `env:"GRPC_PORT" env-required:"true"`
}

func loadGRPCServerCfg(cfg *GRPCServerConfig) {
	cfg.Host = config.GetEnv("GRPC_HOST", "user")
	cfg.Port = config.GetEnv("GRPC_PORT", "50051")
}

func loadPostgresCfg(cfg *postgres.PostgresConfig) {
	cfg.Host = config.GetEnv("POSTGRES_HOST", "localhost")
	cfg.Port = config.GetEnv("POSTGRES_PORT", "5432")
	cfg.User = config.GetEnv("POSTGRES_USER", "user")
	cfg.Password = config.GetEnv("POSTGRES_PASSWORD", "password")
	cfg.DBName = config.GetEnv("POSTGRES_DB", "orpheon")
	cfg.SSLMode = config.GetEnv("POSTGRES_SSL_MODE", "disable")
	cfg.StartTimeout = config.GetEnvAsDuration("POSTGRES_START_TIMEOUT", 5*time.Second)
}

func Load() *Config {
	envFile := flag.String("config", "", "path to env file")
	flag.Parse()

	if *envFile == "" {
		*envFile = defaultEnvPath
	}
	log.Printf("load config from %s...\n", *envFile)
	if err := godotenv.Load(*envFile); err != nil {
		log.Println("Warning: could not load env file:", *envFile, err)
	}

	cfg := &Config{}
	loadPostgresCfg(&cfg.Postgres)
	loadGRPCServerCfg(&cfg.GRPCServer)

	return cfg
}
