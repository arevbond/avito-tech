package main

import (
	"banners/cmd/avito-tech/config"
	"banners/internal/clients"
	"fmt"
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

	//_, err := postgres.New(log, cfg.Storage)
	//if err != nil {
	//	log.Error("can't init storage", "error", err)
	//	os.Exit(1)
	//}

	usersClient := clients.New(cfg.UsersService)
	isVerified, err := usersClient.VerifyToken("eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhc2Rhc2QifQ.jPX4lm8T8fmI6-8-1G17FShhnqSyD5tijv127AvorcAI4uS_wB1uzOzZMNuYRR9_tWhMgmy3GT0oRu3z4NdJBQ")
	fmt.Println(isVerified)
	isAdmin, err := usersClient.IsAdmin("eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhc2Rhc2QifQ.jPX4lm8T8fmI6-8-1G17FShhnqSyD5tijv127AvorcAI4uS_wB1uzOzZMNuYRR9_tWhMgmy3GT0oRu3z4NdJBQ")
	if err != nil {
		log.Error("error", "error", err)
		os.Exit(1)
	}
	fmt.Println(isAdmin)
	//app := application.New(cfg, log)
	//if err = app.Run(); err != nil {
	//	log.Error("application stopped with error", "error", err)
	//	os.Exit(1)
	//} else {
	//	log.Info("application stopped")
	//}

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
