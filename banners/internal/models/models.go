package models

import "time"

type Tag struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Feature struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Banner struct {
	ID        int       `json:"banner_id"`
	TagIDs    []int     `json:"tag_ids"`
	FeatureID int       `json:"feature_id"`
	Content   Content   `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBanner struct {
	TagIDS    []int   `json:"tag_ids"`
	FeatureID int     `json:"feature_id"`
	Content   Content `json:"content"`
	IsActive  bool    `json:"is_active"`
}

type Content struct {
	Title string `json:"title" db:"title"`
	Text  string `json:"text" db:"text"`
	Url   string `json:"url" db:"url"`
}

type BannersParams struct {
	FeatureID int
	TagID     int
	Offset    int
	Limit     int
}

type BannerParams struct {
	FeatureID int
	TagID     int
}
