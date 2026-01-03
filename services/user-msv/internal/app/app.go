package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/hahaclassic/orpheon/api/user-msv/v1/proto"
	"github.com/hahaclassic/orpheon/pkg/http/router"
	"github.com/hahaclassic/orpheon/pkg/infrastructure/postgres"
	"github.com/hahaclassic/orpheon/services/user-msv/config"
	grpc_ctrl "github.com/hahaclassic/orpheon/services/user-msv/internal/controller/grpc"
	http_ctrl "github.com/hahaclassic/orpheon/services/user-msv/internal/controller/http"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/domain/service"
	"github.com/hahaclassic/orpheon/services/user-msv/internal/providers/http/middleware"
	user_postgres "github.com/hahaclassic/orpheon/services/user-msv/internal/repository/postgres"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pool := postgres.NewPostgresPool(ctx, postgres.PostgresConfig(cfg.Postgres))
	defer pool.Close()

	repo := user_postgres.NewUserRepository(pool)
	userService := service.New(repo)
	internalCtrl := grpc_ctrl.NewUserController(userService)
	externalCtrl := http_ctrl.NewUserController(userService, middleware.ClaimsRequired())

	// setup http server
	ginRouter := router.SetupRouter("/api/v1",
		[]router.RoutersRegistrator{
			externalCtrl,
		}, nil)

	srv := &http.Server{
		Addr:    net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler: ginRouter,
	}

	listen, err := net.Listen("tcp", net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port))
	if err != nil {
		slog.Error("failed to listen: ", "err", err)

		return
	}

	// setup grpc server
	server := grpc.NewServer()
	proto.RegisterUserServiceServer(server, internalCtrl)

	go func() {
		slog.Info("[USER-MSV]: (GRPC) start listening at", "addr", listen.Addr())

		err := server.Serve(listen)
		if err != nil {
			slog.Error("[USER-MSV]: failed to serve", "err", err)
		}
	}()

	go func() {
		slog.Info("[USER-MSV]: (HTTP) start listening at", "addr", srv.Addr)

		if err := srv.ListenAndServe(); err != nil {
			slog.Error("[USER-MSV]: failed to serve", "err", err)
		}
	}()

	<-ctx.Done()

	server.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("[USER-MSV]: server forced to shutdown", "err", err)
	}
	slog.Info("the server has shut down.")
}
