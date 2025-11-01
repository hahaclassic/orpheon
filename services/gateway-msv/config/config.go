package config

import (
	"time"

	"github.com/hahaclassic/orpheon/pkg/config"
)

type Config struct {
	HTTP          HTTPServerConfig
	UserGRPC      UserMsvGRPCClient
	AuthGRPC      AuthMsvGRPCClient
	ContentGRPC   ContentMsvGRPCClient
	StreamingGRPC StreaminMsvGRPCClient
	Cookie        CookieConfig
	Logger        LoggerConfig
}

func (cfg *Config) Loaders() []config.InternalConfigLoader {
	return []config.InternalConfigLoader{
		&cfg.HTTP,
		&cfg.UserGRPC,
		&cfg.AuthGRPC,
		&cfg.ContentGRPC,
		&cfg.StreamingGRPC,
	}
}

type HTTPServerConfig struct {
	Host string
	Port string
}

func (cfg *HTTPServerConfig) Load() {
	cfg.Host = config.GetEnv("HOST", "localhost")
	cfg.Port = config.GetEnv("PORT", "8080")
}

type UserMsvGRPCClient struct {
	Host string `env:"USER_MSV_GRPC_HOST" env-required:"true"`
	Port string `env:"USER_MSV_GRPC_PORT" env-required:"true"`
}

func (cfg *UserMsvGRPCClient) Load() {
	cfg.Host = config.GetEnv("USER_MSV_GRPC_HOST", "user")
	cfg.Port = config.GetEnv("USER_MSV_GRPC_PORT", "50051")
}

type AuthMsvGRPCClient struct {
	Host string `env:"AUTH_MSV_GRPC_HOST" env-required:"true"`
	Port string `env:"AUTH_MSV_GRPC_PORT" env-required:"true"`
}

func (cfg *AuthMsvGRPCClient) Load() {
	cfg.Host = config.GetEnv("AUTH_MSV_GRPC_HOST", "auth")
	cfg.Port = config.GetEnv("AUTH_MSV_GRPC_PORT", "50052")
}

type ContentMsvGRPCClient struct {
	Host string `env:"CONTENT_MSV_GRPC_HOST" env-required:"true"`
	Port string `env:"CONTENT_MSV_GRPC_PORT" env-required:"true"`
}

func (cfg *ContentMsvGRPCClient) Load() {
	cfg.Host = config.GetEnv("CONTENT_MSV_GRPC_HOST", "content")
	cfg.Port = config.GetEnv("CONTENT_MSV_GRPC_PORT", "50053")
}

type StreaminMsvGRPCClient struct {
	Host string `env:"STREAMING_MSV_GRPC_HOST" env-required:"true"`
	Port string `env:"STREAMING_MSV_GRPC_PORT" env-required:"true"`
}

func (cfg *StreaminMsvGRPCClient) Load() {
	cfg.Host = config.GetEnv("STREAMING_MSV_GRPC_HOST", "streaming")
	cfg.Port = config.GetEnv("STREAMING_MSV_GRPC_PORT", "50054")
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
	cfg.AccessTTL = config.GetEnvAsDuration("COOKIE_ACCESS_TTL", 15*time.Minute)
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL"`
	Path  string `env:"LOG_PATH"`
}

func (cfg *LoggerConfig) Load() {
	cfg.Level = config.GetEnv("LOG_LEVEL", "info")
	cfg.Path = config.GetEnv("LOG_PATH", "../log/log.log")
}
