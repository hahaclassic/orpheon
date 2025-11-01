package config

import (
	"flag"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Load environment variables in the following order:
// 1. If the -config flag is provided, load vars from the specified file.
// 2. If the -config flag is not set, load vars from the default ".env" file.
// 3. For any variables still unset, assign default values.

const defaultEnvPath = ".env"

func GetEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func GetEnvAsInt(name string, fallback int) int {
	vStr := GetEnv(name, "")
	if v, err := strconv.Atoi(vStr); err == nil {
		return v
	}
	return fallback
}

func GetEnvAsDuration(name string, fallback time.Duration) time.Duration {
	vStr := GetEnv(name, "")
	if v, err := time.ParseDuration(vStr); err == nil {
		return v
	}
	return fallback
}

func GetEnvAsBool(name string, fallback bool) bool {
	vStr := GetEnv(name, "")
	v, err := strconv.ParseBool(vStr)
	if err != nil {
		return fallback
	}

	return v
}

type InternalConfigLoader interface {
	Load()
}

type MainConfig interface {
	Loaders() []InternalConfigLoader
}

func Load(cfg MainConfig) {
	envFile := flag.String("config", "", "path to env file")
	flag.Parse()

	if *envFile == "" {
		*envFile = defaultEnvPath
	}
	slog.Info("load config from file...\n", "path", *envFile)
	if err := godotenv.Load(*envFile); err != nil {
		slog.Warn("could not load env file", "path", *envFile, "err", err)
	}

	for _, internalCfg := range cfg.Loaders() {
		internalCfg.Load()
	}
}
