package storage

import (
	"banners/internal/models"
	"context"
)

type Storage interface {
	GetBanner(ctx context.Context, params *models.BannerParams) (*models.Banner, error)
	CreateBanner(ctx context.Context, params *models.CreateBanner) (*models.Banner, error)
}
