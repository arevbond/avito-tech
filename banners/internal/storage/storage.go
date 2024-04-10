package storage

import (
	"banners/internal/models"
	"context"
)

type Storage interface {
	GetBanners(ctx context.Context, params *models.BannersParams) ([]*models.Banner, error)
	GetBanner(ctx context.Context, params *models.BannerParams) (*models.Banner, error)
	CreateBanner(ctx context.Context, params *models.CreateBanner) (*models.Banner, error)
	UpdateBanner(ctx context.Context, id int, params *models.CreateBanner) error
	DeleteBanner(ctx context.Context, bannerID int) error
}
