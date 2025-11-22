package app

import (
	"PR-appointer/config"
	"PR-appointer/internal/router"
	"PR-appointer/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func StartApplication(ctx context.Context) error {
	cfg := config.GetConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("starting application")

	cfg.Client = storage.NewConnection(ctx, cfg)

	r := router.SetupRouter(ctx, cfg.Client)

	addr := fmt.Sprintf("%s:%d", cfg.Env.IPAddress, cfg.Env.APIPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		slog.Info("starting server")
		slog.Info("Swagger UI available at", "url", fmt.Sprintf("http://localhost:8080/swagger/index.html"))

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "err", err)
			panic(err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", "err", err.Error())
		panic(err)
	}
	return nil
}
