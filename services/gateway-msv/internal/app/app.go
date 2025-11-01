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

	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/services/gateway-msv/config"
	grpc_auth "github.com/hahaclassic/orpheon/services/gateway-msv/internal/clients/auth"
	grpc_user "github.com/hahaclassic/orpheon/services/gateway-msv/internal/clients/user"
	auth_ctrl "github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/api/auth"
	user_ctrl "github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/api/user"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/middleware"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/router"
	user_me_router "github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/router/router-registrators/user-me"
	"github.com/hahaclassic/orpheon/services/gateway-msv/internal/controller/http/utils/cookie"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run(conf *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	userServiceAddress := net.JoinHostPort(conf.UserGRPC.Host, conf.UserGRPC.Port)
	userGRPCClient, err := grpc.NewClient(userServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	userService := grpc_user.NewUserService(userGRPCClient)

	authServiceAddress := net.JoinHostPort(conf.AuthGRPC.Host, conf.AuthGRPC.Port)
	authGRPCClient, err := grpc.NewClient(authServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	authService := grpc_auth.NewAuthService(authGRPCClient)

	cookieTokensSetter := cookie.NewCookieTokensSetter(&conf.Cookie)

	authMiddleware := middleware.NewAuthMiddleware(authService, cookieTokensSetter)
	authMiddlewareRequired := authMiddleware.Required()
	// authMiddlewareOptional := authMiddleware.Optional()

	authController := auth_ctrl.NewAuthController(authService, cookieTokensSetter, authMiddlewareRequired)
	userController := user_ctrl.NewUserController(userService)

	meRouter := user_me_router.NewMeRouter(nil, userController,
		nil, authMiddlewareRequired)

	// loggerMiddleware, err := middleware.SetupLoggerMiddleware(conf.Logger.Path, conf.Logger.Level)
	// if err != nil {
	// 	slog.Error("failed to create logger middleware", "err", err)
	// 	return
	// }

	// Initialize router
	router := router.SetupRouter(
		[]router.RoutersRegistrator{
			authController,
			meRouter,
		},
		[]gin.HandlerFunc{})

	// addr := net.JoinHostPort(conf.HTTP.Host, conf.HTTP.Port)
	// if err := router.Run(addr); err != nil {
	// 	slog.Error("failed to start HTTP server", "err", err)
	// }
	addr := net.JoinHostPort(conf.HTTP.Host, conf.HTTP.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		slog.Info("starting server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	} else {
		slog.Info("server exited properly")
	}
}
