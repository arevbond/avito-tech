package main

import (
	"banners/cmd/avito-tech/config"
	"banners/internal/models"
	"banners/internal/storage/postgres"
	"context"
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

	db, err := postgres.New(log, cfg.Storage)
	if err != nil {
		log.Error("can't init storage", "error", err)
		os.Exit(1)
	}
	//db.CreateTag(context.Background(), &models.Tag{2, "2"})
	//db.CreateTag(context.Background(), &models.Tag{3, "3"})

	//id, err := db.CreateBanner(context.Background(), &models.CreateBanner{
	//	TagIDS:    []int{1, 2},
	//	FeatureID: 1,
	//	Content: models.Content{
	//		Title: "title",
	//		Text:  "text",
	//		Url:   "url",
	//	},
	//	IsActive: false,
	//})
	//if err != nil {
	//	log.Error("can't create banner", "error", err)
	//	os.Exit(1)
	//}
	//log.Info("success create banner", slog.Int("id", id))

	//err = db.UpdateBanner(context.Background(), 1, &models.CreateBanner{
	//	TagIDS:    []int{1},
	//	FeatureID: 1,
	//	Content: models.Content{
	//		Title: "title",
	//		Text:  "text",
	//		Url:   "url",
	//	},
	//	IsActive: false,
	//})
	//if err != nil {
	//	log.Error("can't update banner", "error", err)
	//	os.Exit(1)
	//}
	//log.Info("success update banner", slog.Int("id", 1))
	//
	//err = db.DeleteBanner(context.Background(), 1)
	//if err != nil {
	//	log.Error("can't delete banner", "error", err)
	//	os.Exit(1)
	//}
	//log.Info("success delete banner")
	banners, err := db.GetBanners(context.Background(), &models.BannerParams{
		FeatureID: 0,
		TagID:     1,
		Limit:     10,
		Offset:    0,
	})
	if err != nil {
		log.Error("can't get banners", "error", err)
		os.Exit(1)
	}
	fmt.Println(banners)
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
