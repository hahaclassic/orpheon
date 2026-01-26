package config

import (
	"time"

	"github.com/hahaclassic/orpheon/pkg/config"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/minio"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/postgres"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/redis"
)

type Config struct {
	Server ServerConfig

	UserInfo UserInfoClientCofig

	Postgres     PostgresConfig
	Redis        RedisConfig
	MinIO        MinIOConfig
	MinIOBuckets MinIOBucketsConfig

	RedisAccessCache RedisAccessMetaConfig
	LocalAccessCache LocalAccessMetaConfig
}

func (cfg *Config) Loaders() []config.InternalConfigLoader {
	return []config.InternalConfigLoader{
		&cfg.Postgres,
		&cfg.Redis,
		&cfg.Server,
	}
}

type ServerConfig struct {
	Host     string `env:"HOST" env-required:"true"`
	Port     string `env:"PORT" env-required:"true"`
	ReadOnly bool   `env:"READ_ONLY" env-required:"true"`
}

func (cfg *ServerConfig) Load() {
	cfg.Host = config.GetEnv("HOST", "localhost")
	cfg.Port = config.GetEnv("PORT", "8080")
	cfg.ReadOnly = config.GetEnvAsBool("READ_ONLY", false)
}

type UserInfoClientCofig struct {
	Host string `env:"USER_CREATOR_HOST" env-required:"true"`
	Port string `env:"USER_CREATOR_PORT" env-required:"true"`
	URL  string `env:"USER_CREATOR_URL" env-required:"true"`
}

func (cfg *UserInfoClientCofig) Load() {
	cfg.Host = config.GetEnv("USER_CREATOR_HOST", "localhost")
	cfg.Port = config.GetEnv("USER_CREATOR_PORT", "50051")
	cfg.URL = config.GetEnv("USER_CREATOR_URL", "/users")
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

type RedisConfig redis.RedisConfig

func (cfg *RedisConfig) Load() {
	cfg.Addr = config.GetEnv("REDIS_ADDR", "localhost:6379")
	cfg.Password = config.GetEnv("REDIS_PASSWORD", "")
	cfg.DB = config.GetEnvAsInt("REDIS_DB", 0)
}

type MinIOConfig minio.MinIOConfig

func (cfg *MinIOConfig) Load() {
	cfg.Endpoint = config.GetEnv("MINIO_ENDPOINT", "localhost:9000")
	cfg.AccessKey = config.GetEnv("MINIO_ROOT_USER", "admin")
	cfg.SecretKey = config.GetEnv("MINIO_ROOT_PASSWORD", "adminstrongpassword")
	cfg.Secure = config.GetEnvAsBool("MINIO_SECURE", false)
}

type MinIOBucketsConfig struct {
	BucketPlaylist     string `env:"MINIO_BUCKET_PLAYLIST_COVERS"`
	BucketAlbum        string `env:"MINIO_BUCKET_ALBUM_COVERS"`
	BucketArtistAvatar string `env:"MINIO_BUCKET_ARTIST_AVATARS"`
	BucketAudio        string `env:"MINIO_BUCKET_AUDIO_FILES"`
}

func (cfg *MinIOBucketsConfig) Load() {
	cfg.BucketPlaylist = config.GetEnv("MINIO_BUCKET_PLAYLIST_COVERS", "playlist-covers")
	cfg.BucketAlbum = config.GetEnv("MINIO_BUCKET_ALBUM_COVERS", "album-covers")
	cfg.BucketArtistAvatar = config.GetEnv("MINIO_BUCKET_ARTIST_AVATARS", "artist-avatars")
	cfg.BucketAudio = config.GetEnv("MINIO_BUCKET_AUDIO_FILES", "audio-files")
}

type RedisAccessMetaConfig struct {
	TTL    time.Duration `env:"ACCESS_META_TTL"`
	Jitter time.Duration `env:"ACCESS_META_JITTER"`
}

func (cfg *RedisAccessMetaConfig) Load() {
	cfg.TTL = config.GetEnvAsDuration("ACCESS_META_TTL", 1*time.Hour)
	cfg.Jitter = config.GetEnvAsDuration("ACCESS_META_JITTER", 5*time.Minute)
}

type LocalAccessMetaConfig struct {
	Size int `env:"LOCAL_CACHE_SIZE"`
}

func (cfg *LocalAccessMetaConfig) Load() {
	cfg.Size = config.GetEnvAsInt("LOCAL_CACHE_SIZE", 128)
}
