package config

import (
	"log"

	"github.com/hahaclassic/orpheon/backend/internal/infrastructure/minio"
	"github.com/hahaclassic/orpheon/backend/internal/infrastructure/postgres"
	"github.com/hahaclassic/orpheon/backend/internal/infrastructure/redis"
	access_cache_redis "github.com/hahaclassic/orpheon/backend/internal/repository/content/playlist/access-cache/redis"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = ".env"
)

type Config struct {
	Postgres    postgres.PostgresConfig
	MinIO       minio.MinioConfig
	Redis       redis.RedisConfig
	AccessCache AccessCacheConfig
}

type AccessCacheConfig struct {
	TTL     access_cache_redis.TTLConfig
	LRUSize int
}

func MustLoad() *Config {
	config := &Config{}
	err := cleanenv.ReadConfig(configPath, config)
	if err != nil {
		log.Fatalf("Error while loading config: %s", err)
	}

	return config
}
