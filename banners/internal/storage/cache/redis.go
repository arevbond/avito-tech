package cache

import (
	"banners/cmd/avito-tech/config"
	"banners/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
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
