package main

import (
	"banners/cmd/avito-tech/application"
	"banners/cmd/avito-tech/config"
	"log"
	"log/slog"
	"os"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	log := setupLogger(cfg.Env)

	app := application.New(cfg, log)
	if err = app.Run(); err != nil {
		log.Error("application stopped with error", "error", err)
		os.Exit(1)
	} else {
		log.Info("application stopped")
	}

}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
