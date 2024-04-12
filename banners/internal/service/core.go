package service

import (
	"banners/internal/clients"
	"banners/internal/models"
	"banners/internal/storage"
	"banners/internal/storage/cache"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrUnauthorized = errors.New("user unauthorized")
	ErrNotFound     = errors.New("not found")
	ErrForbidden    = errors.New("forbidden")

	ErrIncorrectData = errors.New("incorrect data")
)

type Service interface {
	UserBanner(ctx context.Context, token string, params *UserBannerParams) (*models.Content, error)
	Banners(ctx context.Context, token string, params *models.BannersParams) ([]*models.Banner, error)
	CreateBanner(ctx context.Context, token string, banner *models.CreateBanner) (int, error)
	UpdateBanner(ctx context.Context, token string, id int, updateBanner *models.CreateBanner) error
	DeleteBanner(ctx context.Context, token string, id int) error
}

type BannerService struct {
	Storage     storage.Storage
	Cache       *cache.Cache
	UserService *clients.Users
	Log         *slog.Logger
}

type UserBannerParams struct {
	TagID           int
	FeatureID       int
	UseLastRevision bool
}

func (s *BannerService) UserBanner(ctx context.Context, token string, params *UserBannerParams) (*models.Content, error) {
	isValidToken, err := s.UserService.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("can't verify token: %w", err)
	}
	if !isValidToken {
		return nil, ErrUnauthorized
	}

	if !params.UseLastRevision {
		bannerFromCache, err := s.Cache.GetBannerByTagAndFeature(ctx, params.TagID, params.FeatureID)
		if err != nil {
			s.Log.Error("can't get banner from cache", "error", err)
		} else if bannerFromCache != nil {
			if bannerFromCache.IsActive {
				return &bannerFromCache.Content, nil
			}
			isAdmin, err := s.UserService.IsAdmin(token)
			if err != nil {
				return nil, fmt.Errorf("can't verify admin token: %w", err)
			}
			if isAdmin {
				return &bannerFromCache.Content, nil
			}
			return nil, ErrForbidden
		}
	}

	userBanner, err := s.Storage.GetBanner(ctx, &models.BannerParams{
		FeatureID: params.FeatureID,
		TagID:     params.TagID,
	})
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("storage error: %w", err)
	}

	if !params.UseLastRevision {
		err = s.Cache.AddBanner(ctx, userBanner)
		if err != nil {
			s.Log.Error("can't add banner to cache", "error", err)
		}
	}

	if userBanner.IsActive {
		return &userBanner.Content, nil
	}

	isAdmin, err := s.UserService.IsAdmin(token)
	if err != nil {
		return nil, fmt.Errorf("can't verify admin token: %w", err)
	}
	if isAdmin {
		return &userBanner.Content, nil
	}
	return nil, ErrForbidden
}

func (s *BannerService) Banners(ctx context.Context, token string, params *models.BannersParams) ([]*models.Banner, error) {
	isValidToken, err := s.UserService.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("can't verify token: %w", err)
	}
	if !isValidToken {
		return nil, ErrUnauthorized
	}

	isAdmin, err := s.UserService.IsAdmin(token)
	if err != nil {
		return nil, fmt.Errorf("can't verify admin token: %w", err)
	}
	if !isAdmin {
		return nil, ErrForbidden
	}

	banners, err := s.Storage.GetBanners(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("storage error: %w", err)
	}

	return banners, nil
}

func (s *BannerService) CreateBanner(ctx context.Context, token string, createBanner *models.CreateBanner) (int, error) {
	isValidToken, err := s.UserService.VerifyToken(token)
	if err != nil {
		return -1, fmt.Errorf("can't verify token: %w", err)
	}
	if !isValidToken {
		return -1, ErrUnauthorized
	}
	isAdmin, err := s.UserService.IsAdmin(token)
	if err != nil {
		return -1, fmt.Errorf("can't verify admin token: %w", err)
	}
	if !isAdmin {
		return -1, ErrForbidden
	}

	banner, err := s.Storage.CreateBanner(ctx, createBanner)
	if err != nil {
		s.Log.Error("can't create banner", "error", err)
		return -1, ErrIncorrectData
	}

	err = s.Cache.AddBanner(ctx, banner)
	if err != nil {
		s.Log.Error("can't insert banner into cache", "error", err)
	}
	return banner.ID, nil
}

func (s *BannerService) UpdateBanner(ctx context.Context, token string, id int, updateBanner *models.CreateBanner) error {
	isValidToken, err := s.UserService.VerifyToken(token)
	if err != nil {
		return fmt.Errorf("can't verify token: %w", err)
	}
	if !isValidToken {
		return ErrUnauthorized
	}
	isAdmin, err := s.UserService.IsAdmin(token)
	if err != nil {
		return fmt.Errorf("can't verify admin token: %w", err)
	}
	if !isAdmin {
		return ErrForbidden
	}

	err = s.Storage.UpdateBanner(ctx, id, updateBanner)
	if err != nil {
		return fmt.Errorf("can't update banner: %w", err)
	}

	err = s.Cache.UpdateBanner(ctx, id, updateBanner)
	if err != nil {
		s.Log.Error("can''t update banner in cache", "error", err)
	}

	return nil
}

func (s *BannerService) DeleteBanner(ctx context.Context, token string, id int) error {
	isValidToken, err := s.UserService.VerifyToken(token)
	if err != nil {
		return fmt.Errorf("can't verify token: %w", err)
	}
	if !isValidToken {
		return ErrUnauthorized
	}
	isAdmin, err := s.UserService.IsAdmin(token)
	if err != nil {
		return fmt.Errorf("can't verify admin token: %w", err)
	}
	if !isAdmin {
		return ErrForbidden
	}

	err = s.Cache.DeleteBanner(ctx, id)
	if err != nil {
		s.Log.Error("can''t delete banner from cache", "error", err)
	}

	err = s.Storage.DeleteBanner(ctx, id)
	if err != nil {
		return fmt.Errorf("can't delete banner: %w", err)
	}

	return nil
}
