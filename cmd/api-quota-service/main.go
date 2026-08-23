package main

import (
	"context"
	"errors"
	"github.com/example/api-quota-service/internal/application"
	configpkg "github.com/example/api-quota-service/internal/config"
	httpapi "github.com/example/api-quota-service/internal/http"
	"github.com/example/api-quota-service/internal/infrastructure"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := configpkg.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	pol := infrastructure.NewMemoryPolicies()
	svcs := infrastructure.NewMemoryServices()
	events := infrastructure.NewMemoryEvents()
	counter := infrastructure.NewMemoryCounter()
	pub := infrastructure.NewMemoryPublisher()
	app := application.NewService(pol, svcs, events, counter, pub)
	h := httpapi.NewHandler(app, logger)
	srv := &http.Server{Addr: cfg.Address, Handler: h.Routes()}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		logger.Info("api quota service started", "address", cfg.Address)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", "error", err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
	logger.Info("service stopped")
}
