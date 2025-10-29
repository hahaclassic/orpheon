package app

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"syscall"

	"github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/postgres"
	"github.com/hahaclassic/orpheon/services/user-msv/config"
	grpc_ctrl "github.com/hahaclassic/orpheon/services/user-msv/internal/controller/grpc"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/service"
	user_postgres "github.com/hahaclassic/orpheon/services/user-msv/internal/repository/postgres"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool := postgres.NewPostgresPool(ctx, cfg.Postgres)
	repo := user_postgres.NewUserRepository(pool)
	userService := service.New(repo)
	ctrl := grpc_ctrl.NewUserController(userService)

	address := net.JoinHostPort(cfg.GRPCServer.Host, cfg.GRPCServer.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("failed to listen: ", "err", err)

		return
	}

	server := grpc.NewServer()
	proto.RegisterUserServiceServer(server, ctrl)

	go func() {
		slog.Info("[USER-MSV]: start listening at ", "addr", listen.Addr())

		err := server.Serve(listen)
		if err != nil {
			slog.Error("[USER-MSV]: failed to serve: ", "err", err)
		}
	}()

	<-ctx.Done()
	server.GracefulStop()
	slog.Info("the server has shut down.")
}
