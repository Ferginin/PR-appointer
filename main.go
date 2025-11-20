package main

import (
	"PR-appointer/cmd/app"
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	_ "PR-appointer/docs"
)

// @title PR Reviewer Assignment Service (Test Task, Fall 2025)
// @version 1.0
// @description API для управления пул реквестами
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@clinic.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:6060
// @BasePath

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {

	slog.Info("Starting main")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error, 1)

	go func() {
		if err := app.StartApplication(ctx); err != nil {
			errChan <- err
		} else {
			errChan <- nil
		}
	}()

	select {
	case <-ctx.Done():
		slog.Info("Stopping main")

	case err := <-errChan:
		if err != nil {
			slog.Error("Error during application start:", err)
		} else {
			slog.Info("Application stopped")
		}
	}

	slog.Info("Shutdown completed")
}
