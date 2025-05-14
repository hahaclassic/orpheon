package config

import (
	"log"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = ".env"
)

type HTTPConfig struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

type PostgresConfig struct {
	Host         string        `env:"POSTGRES_HOST"`
	Port         string        `env:"POSTGRES_PORT"`
	User         string        `env:"POSTGRES_USER"`
	Password     string        `env:"POSTGRES_PASSWORD"`
	DB           string        `env:"POSTGRES_DB"`
	SSLMode      string        `env:"POSTGRES_SSL_MODE"`
	StartTimeout time.Duration `env:"POSTGRES_START_TIMEOUT"`
}

type MinIOConfig struct {
	Endpoint           string `env:"MINIO_ENDPOINT"`
	AccessKey          string `env:"MINIO_ROOT_USER"`
	SecretKey          string `env:"MINIO_ROOT_PASSWORD"`
	Secure             bool   `env:"MINIO_SECURE"`
	BucketPlaylist     string `env:"MINIO_BUCKET_PLAYLIST_COVERS"`
	BucketAlbum        string `env:"MINIO_BUCKET_ALBUM_COVERS"`
	BucketArtistAvatar string `env:"MINIO_BUCKET_ARTIST_AVATARS"`
	BucketAudio        string `env:"MINIO_BUCKET_AUDIO_FILES"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR"`
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB"`
}

type PasswordHasherConfig struct {
	Cost int `env:"PASSWORD_HASHER_COST"`
}

type AccessTokenConfig struct {
	TTL       time.Duration `env:"ACCESS_TOKEN_TTL"`
	Jitter    time.Duration `env:"ACCESS_TOKEN_JITTER"`
	SecretKey []byte        `env:"SECRET_KEY"`
}

type RefreshTokenConfig struct {
	TTL    time.Duration `env:"REFRESH_TOKEN_TTL"`
	Jitter time.Duration `env:"REFRESH_TOKEN_JITTER"`
}

type RedisAccessMetaConfig struct {
	TTL    time.Duration `env:"ACCESS_META_TTL"`
	Jitter time.Duration `env:"ACCESS_META_JITTER"`
}

type LocalAccessMetaConfig struct {
	Size int `env:"LOCAL_CACHE_SIZE"`
}

type Config struct {
	HTTP                 HTTPConfig
	Postgres             PostgresConfig
	MinIO                MinIOConfig
	Redis                RedisConfig
	PasswordHasher       PasswordHasherConfig
	AccessToken          AccessTokenConfig
	RefreshToken         RefreshTokenConfig
	RedisAccessMetaCache RedisAccessMetaConfig
	LocalAccessMetaCache LocalAccessMetaConfig
}

var (
	cfg  *Config
	once sync.Once
)

func MustLoad() *Config {
	once.Do(func() {
		log.Println("Loading config from environment variables...")
		cfg = &Config{}
		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	})
	return cfg
}

// type Config struct {
// 	Postgres    postgres.PostgresConfig
// 	MinIO       minio.MinioConfig
// 	Redis       redis.RedisConfig
// 	AccessCache AccessCacheConfig
// 	Refresh     refresh_redis.TTLConfig
// }

// type AccessCacheConfig struct {
// 	TTL     access_cache_redis.TTLConfig
// 	LRUSize int
// }

// func MustLoad() *Config {
// 	config := &Config{}
// 	err := cleanenv.ReadConfig(configPath, config)
// 	if err != nil {
// 		log.Fatalf("Error while loading config: %s", err)
// 	}

// 	return config
// }
