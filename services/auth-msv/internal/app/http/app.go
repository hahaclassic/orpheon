package app

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	userproto "github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
	"github.com/hahaclassic/orpheon/pkg/http/router"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/postgres"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/redis"
	"github.com/hahaclassic/orpheon/services/auth-msv/config"
	http_ctrl "github.com/hahaclassic/orpheon/services/auth-msv/internal/controller/http"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/domain/service"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/providers/http/cookie"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/providers/http/middleware"
	bcrypt_hasher "github.com/hahaclassic/orpheon/services/auth-msv/internal/providers/password-hasher/bcrypt-hasher"
	jwttokens "github.com/hahaclassic/orpheon/services/auth-msv/internal/providers/tokens/jwt"
	grpc_user_creator "github.com/hahaclassic/orpheon/services/auth-msv/internal/providers/user-creator/grpc"
	creds_postgres "github.com/hahaclassic/orpheon/services/auth-msv/internal/repository/credentials/postgres"
	refresh_redis "github.com/hahaclassic/orpheon/services/auth-msv/internal/repository/refresh-token/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool := postgres.NewPostgresPool(ctx, postgres.PostgresConfig(cfg.Postgres))
	defer pool.Close()

	redisClient, err := redis.NewRedisClient(redis.RedisConfig(cfg.Redis))
	if err != nil {
		slog.Error("[AUTH-MSV]: failed to connect to redis", "err", err)
		return
	}
	defer func() {
		err := redisClient.Close()
		if err != nil {
			slog.Error("[AUTH-MSV]: failed to close redis connection", "err", err)
		}
	}()

	passwordHasher := bcrypt_hasher.New(cfg.PasswordHasher.Cost)
	jwtTokenService := jwttokens.New(cfg.AccessToken)
	cookieTokensSetter := cookie.NewCookieTokensSetter(&cfg.Cookie)

	credentialsRepo := creds_postgres.NewCredentialsRepository(pool)
	refreshTokenRepo := refresh_redis.NewRefreshTokenRepository(redisClient, &cfg.RefreshToken)

	userCreatorAddress := net.JoinHostPort(cfg.UserCreatorService.Host, cfg.UserCreatorService.Port)
	userCreatorGRPCClient, err := grpc.NewClient(userCreatorAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	userproto.NewUserServiceClient(userCreatorGRPCClient)
	userCreatorService := grpc_user_creator.NewUserCreatorService(userCreatorGRPCClient)

	authService := service.NewAuthService(credentialsRepo,
		refreshTokenRepo, userCreatorService, passwordHasher, jwtTokenService)
	authCtrl := http_ctrl.NewAuthController(authService, cookieTokensSetter, middleware.ClaimsRequired())

	ginRouter := router.SetupRouter("/api/v1",
		[]router.RoutersRegistrator{
			authCtrl,
		}, nil)

	srv := &http.Server{
		Addr:    net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler: ginRouter,
	}

	go func() {
		slog.Info("[AUTH-MSV]: start listening at", "addr", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			slog.Error("[AUTH-MSV]: failed to serve", "err", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("[AUTH-MSV]: server forced to shutdown", "err", err)
	} else {
		slog.Info("[AUTH-MSV]: server exited properly")
	}
}
