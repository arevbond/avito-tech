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

type BannerParams struct {
	FeatureID int `db:"feature_id"`
	TagID     int `db:"tag_id"`
	Offset    int
	Limit     int
}
