package config

import (
	"time"

	"github.com/hahaclassic/orpheon/pkg/config"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/postgres"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/redis"
)

type Config struct {
	Postgres           PostgresConfig
	Redis              RedisConfig
	Server             ServerConfig
	UserCreatorService UserCreatorClientCofig
	PasswordHasher     PasswordHasherConfig
	AccessToken        AccessTokenConfig
	RefreshToken       RefreshTokenConfig
	Cookie             CookieConfig
}

func (cfg *Config) Loaders() []config.InternalConfigLoader {
	return []config.InternalConfigLoader{
		&cfg.Postgres,
		&cfg.Redis,
		&cfg.Server,
		&cfg.UserCreatorService,
		&cfg.PasswordHasher,
		&cfg.AccessToken,
		&cfg.RefreshToken,
	}
}

type ServerConfig struct {
	Host string `env:"HOST" env-required:"true"`
	Port string `env:"PORT" env-required:"true"`
}

func (cfg *ServerConfig) Load() {
	cfg.Host = config.GetEnv("HOST", "localhost")
	cfg.Port = config.GetEnv("PORT", "8080")
}

type UserCreatorClientCofig struct {
	Host string `env:"USER_CREATOR_HOST" env-required:"true"`
	Port string `env:"USER_CREATOR_PORT" env-required:"true"`
	URL  string `env:"USER_CREATOR_URL" env-required:"true"`
}

func (cfg *UserCreatorClientCofig) Load() {
	cfg.Host = config.GetEnv("USER_CREATOR_HOST", "localhost")
	cfg.Port = config.GetEnv("USER_CREATOR_PORT", "50051")
	cfg.URL = config.GetEnv("USER_CREATOR_URL", "/users")
}

type PasswordHasherConfig struct {
	Cost int `env:"PASSWORD_HASHER_COST"`
}

func (cfg *PasswordHasherConfig) Load() {
	cfg.Cost = config.GetEnvAsInt("PASSWORD_HASHER_COST", 12)
}

type AccessTokenConfig struct {
	TTL       time.Duration `env:"ACCESS_TOKEN_TTL"`
	Jitter    time.Duration `env:"ACCESS_TOKEN_JITTER"`
	SecretKey []byte        `env:"SECRET_KEY"`
}

func (cfg *AccessTokenConfig) Load() {
	cfg.TTL = config.GetEnvAsDuration("ACCESS_TOKEN_TTL", 3*time.Hour)
	cfg.Jitter = config.GetEnvAsDuration("ACCESS_TOKEN_JITTER", 10*time.Minute)
	cfg.SecretKey = []byte(config.GetEnv("ACCESS_SECRET_KEY", "secret_key"))
}

type RefreshTokenConfig struct {
	TTL    time.Duration `env:"REFRESH_TOKEN_TTL"`
	Jitter time.Duration `env:"REFRESH_TOKEN_JITTER"`
}

func (cfg *RefreshTokenConfig) Load() {
	cfg.TTL = config.GetEnvAsDuration("REFRESH_TOKEN_TTL", 720*time.Hour)
	cfg.Jitter = config.GetEnvAsDuration("REFRESH_TOKEN_JITTER", 5*time.Hour)
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

type CookieConfig struct {
	Domain     string        `env:"COOKIE_DOMAIN"`
	Path       string        `env:"COOKIE_PATH"`
	Secure     bool          `env:"COOKIE_SECURE"`
	HttpOnly   bool          `env:"COOKIE_HTTP_ONLY"`
	RefreshTTL time.Duration `env:"COOKIE_REFRESH_TTL"`
	AccessTTL  time.Duration `env:"COOKIE_ACCESS_TTL"`
}

func (cfg *CookieConfig) Load() {
	cfg.Domain = config.GetEnv("COOKIE_DOMAIN", "")
	cfg.Path = config.GetEnv("COOKIE_PATH", "/")
	cfg.Secure = config.GetEnvAsBool("COOKIE_SECURE", false)
	cfg.HttpOnly = config.GetEnvAsBool("COOKIE_HTTP_ONLY", true)
	cfg.RefreshTTL = config.GetEnvAsDuration("COOKIE_REFRESH_TTL", 725*time.Hour)
	cfg.AccessTTL = config.GetEnvAsDuration("COOKIE_ACCESS_TTL", 16*time.Minute)
}
