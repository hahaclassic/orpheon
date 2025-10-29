package config

import (
	"os"
	"strconv"
	"time"
)

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
