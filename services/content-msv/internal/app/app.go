package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hahaclassic/orpheon/services/content-msv/config"
)

type App struct {
	infra *Infra
}

func (a *App) Build(ctx context.Context, cfg *config.Config) (*http.Server, error) {
	a.infra = &Infra{}
	err := a.infra.Init(ctx, cfg)
	if err != nil {
		return nil, err
	}

	repos, err := InitRepositories(ctx, cfg, a.infra)
	if err != nil {
		a.Close(ctx)
		return nil, err
	}

	services, err := InitServices(cfg, repos)
	if err != nil {
		a.Close(ctx)
		return nil, err
	}

	return InitServerHTTP(cfg, services), nil
}

func (a *App) Close(ctx context.Context) error {
	if a.infra != nil {
		return a.infra.Close(ctx)
	}
	return nil
}

func Run(cfg *config.Config) {
	buildCtx, buildCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer buildCancel()

	app := &App{}
	server, err := app.Build(buildCtx, cfg)
	if err != nil {
		slog.Error("[CONTENT-MSV]: failed to start service", "err", err)
		buildCancel()
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("[CONTENT-MSV]: starting server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("[CONTENT-MSV]: server error", "err", err)
		}
	}()

	<-ctx.Done()
	slog.Info("[CONTENT-MSV]: shutdown signal received")

	shutdownSrvCtx, srvCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer srvCancel()

	if err := server.Shutdown(shutdownSrvCtx); err != nil {
		slog.Error("[CONTENT-MSV]: server forced to shutdown", "err", err)
	}

	shutdownAppCtx, infraCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer infraCancel()
	if err = app.Close(shutdownAppCtx); err != nil {
		slog.Error("[CONTENT-MSV]: failed to close infra", "err", err)
	}

	slog.Info("[CONTENT-MSV]: the server has shut down.")
}
