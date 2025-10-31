package app

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	authproto "github.com/hahaclassic/orpheon/api/auth-msv/v1/proto"
	userproto "github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/postgres"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/redis"
	"github.com/hahaclassic/orpheon/services/auth-msv/config"
	grpc_ctrl "github.com/hahaclassic/orpheon/services/auth-msv/internal/controller/grpc"
	"github.com/hahaclassic/orpheon/services/auth-msv/internal/domain/service"
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
	ctrl := grpc_ctrl.NewAuthController(authService)

	address := net.JoinHostPort(cfg.GRPCServer.Host, cfg.GRPCServer.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("[AUTH-MSV]: failed to listen: ", "err", err)

		return
	}

	server := grpc.NewServer()
	authproto.RegisterAuthServiceServer(server, ctrl)

	go func() {
		slog.Info("[AUTH-MSV]: start listening at", "addr", listen.Addr())

		err := server.Serve(listen)
		if err != nil {
			slog.Error("[AUTH-MSV]: failed to serve", "err", err)
		}
	}()

	<-ctx.Done()
	server.GracefulStop()
	slog.Info("[AUTH-MSV]: the server has shut down.")
}
