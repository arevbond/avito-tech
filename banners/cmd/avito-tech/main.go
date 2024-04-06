package main

import (
	"avito-tech/cmd/avito-tech/application"
	"avito-tech/cmd/avito-tech/config"
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg, err := loadConfig()
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

func loadConfig() (*config.Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	cfg := config.NewDefaultConfig()

	configPath := flag.String("config", "", "config path")
	flag.Parse()

	if err := cleanenv.ReadConfig(*configPath, cfg); err != nil {
		return nil, fmt.Errorf("can't load config: %w", err)
	}

	return cfg, nil
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
