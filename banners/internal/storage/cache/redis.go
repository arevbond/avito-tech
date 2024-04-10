package cache

import (
	"banners/cmd/avito-tech/config"
	"banners/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

type Cache struct {
	Client *redis.Client
}

func New(cfg config.RedisConfig) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
	})
	return &Cache{Client: client}
}

func (c *Cache) AddBanner(ctx context.Context, banner *models.Banner) error {
	bannerJSON, err := json.Marshal(banner)
	if err != nil {
		return err
	}

	err = c.Client.HSet(ctx, fmt.Sprintf("banner:%d", banner.ID), "data", bannerJSON).Err()
	if err != nil {
		return err
	}

	for _, tagID := range banner.TagIDs {
		err = c.Client.SAdd(ctx, fmt.Sprintf("tag:%d:banners", tagID), banner.ID).Err()
		if err != nil {
			return err
		}
	}

	err = c.Client.SAdd(ctx, fmt.Sprintf("feature:%d:banners", banner.FeatureID), banner.ID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) getBannersByTagAndFeatureID(ctx context.Context, tagID, featureID int) ([]*models.Banner, error) {
	tagBannerIDs, err := c.Client.SMembers(ctx, fmt.Sprintf("tag:%d:banners", tagID)).Result()
	if err != nil {
		return nil, err
	}

	featureBannerIDs, err := c.Client.SMembers(ctx, fmt.Sprintf("feature:%d:banners", featureID)).Result()
	if err != nil {
		return nil, err
	}

	intersection := make([]string, 0)
	for _, bannerID := range tagBannerIDs {
		if contains(featureBannerIDs, bannerID) {
			intersection = append(intersection, bannerID)
		}
	}

	banners := make([]*models.Banner, 0)
	for _, bannerID := range intersection {
		bannerData, err := c.Client.HGet(ctx, fmt.Sprintf("banner:%s", bannerID), "data").Result()
		if err != nil {
			return nil, err
		}

		var retrievedBanner *models.Banner
		err = json.Unmarshal([]byte(bannerData), &retrievedBanner)
		if err != nil {
			return nil, err
		}

		banners = append(banners, retrievedBanner)
	}

	return banners, nil
}

func (c *Cache) GetBannerByTagAndFeature(ctx context.Context, tagID, featureID int) (*models.Banner, error) {
	tagBannerIDs, err := c.Client.SMembers(ctx, fmt.Sprintf("tag:%d:banners", tagID)).Result()
	if err != nil {
		return nil, err
	}

	featureBannerIDs, err := c.Client.SMembers(ctx, fmt.Sprintf("feature:%d:banners", featureID)).Result()
	if err != nil {
		return nil, err
	}

	for _, bannerID := range tagBannerIDs {
		if contains(featureBannerIDs, bannerID) {
			bannerData, err := c.Client.HGet(ctx, fmt.Sprintf("banner:%s", bannerID), "data").Result()
			if err != nil {
				return nil, fmt.Errorf("can't get banner from cache: %w", err)
			}

			var retrievedBanner models.Banner
			err = json.Unmarshal([]byte(bannerData), &retrievedBanner)
			if err != nil {
				return nil, fmt.Errorf("can't unmarshal banner: %w", err)
			}
			return &retrievedBanner, nil
		}
	}
	return nil, nil
}

func (c *Cache) UpdateBanner(ctx context.Context, id int, banner *models.CreateBanner) error {
	oldBannerData, err := c.Client.HGet(ctx, fmt.Sprintf("banner:%d", id), "data").Result()
	if err != nil {
		return fmt.Errorf("can't get old banner from cache: %w", err)
	}

	var oldBanner models.Banner
	err = json.Unmarshal([]byte(oldBannerData), &oldBanner)
	if err != nil {
		return fmt.Errorf("can't unmarshal old banner: %w", err)
	}

	err = c.DeleteBannerWithRelation(ctx, &oldBanner)
	if err != nil {
		return fmt.Errorf("can't delete old banner from cache: %w", err)
	}

	newBanner := &models.Banner{
		ID:        id,
		TagIDs:    banner.TagIDS,
		FeatureID: banner.FeatureID,
		Content:   banner.Content,
		IsActive:  banner.IsActive,
		CreatedAt: oldBanner.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err = c.AddBanner(ctx, newBanner)
	if err != nil {
		return fmt.Errorf("can't add new banner to cache: %w", err)
	}

	return nil
}

func (c *Cache) DeleteBanner(ctx context.Context, bannerID int) error {
	oldBannerData, err := c.Client.HGet(ctx, fmt.Sprintf("banner:%d", bannerID), "data").Result()
	if err != nil {
		return fmt.Errorf("can't get old banner from cache: %w", err)
	}

	var oldBanner models.Banner
	err = json.Unmarshal([]byte(oldBannerData), &oldBanner)
	if err != nil {
		return fmt.Errorf("can't unmarshal old banner: %w", err)
	}

	err = c.DeleteBannerWithRelation(ctx, &oldBanner)
	if err != nil {
		return fmt.Errorf("can't delete old banner from cache: %w", err)
	}
	return nil
}

func (c *Cache) DeleteBannerWithRelation(ctx context.Context, banner *models.Banner) error {
	err := c.Client.Del(ctx, fmt.Sprintf("banner:%d", banner.ID)).Err()
	if err != nil {
		return fmt.Errorf("can't delete banner data from cache: %w", err)
	}

	for _, tagID := range banner.TagIDs {
		err = c.Client.SRem(ctx, fmt.Sprintf("tag:%d:banners", tagID), banner.ID).Err()
		if err != nil {
			return fmt.Errorf("can't remove banner ID from tag set: %w", err)
		}
	}

	err = c.Client.SRem(ctx, fmt.Sprintf("feature:%d:banners", banner.FeatureID), banner.ID).Err()
	if err != nil {
		return fmt.Errorf("can't remove banner ID from feature set: %w", err)
	}

	return nil
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func (c *Cache) Close() error {
	return c.Client.Close()
}
