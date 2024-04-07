package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"users/cmd/users/config"
	"users/internal/models"
	"users/internal/storage/postgres"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("can't init config")
	}

	log := setupLogger(cfg.Env)
	log.Info("service users start")
	db, err := postgres.New(log, cfg.Storage)
	if err != nil {
		log.Error("can't init storage", "error", err)
		os.Exit(1)
	}
	user, err := db.CreateUser(context.Background(), &models.UserRegister{
		"adaskdaskdklasmd1",
		"nikita",
		"asdkjahsjdjkhsajdas",
		true,
	})
	if err != nil {
		log.Error("can't create user", "error", err)
	}
	fmt.Println(user)
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
