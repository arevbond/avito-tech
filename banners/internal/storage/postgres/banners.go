package postgres

import (
	"banners/internal/models"
	"banners/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"log"
	"regexp"
	"strconv"
	"time"
)

func (s *Storage) GetBanner(ctx context.Context, params *models.BannerParams) (*models.Banner, error) {
	sql, args, err := sq.Select(
		bannerTable+"."+fieldID,
		bannerTable+"."+fieldFeatureID,
		bannerTable+"."+fieldContent,
		bannerTable+"."+fieldIsActive,
		bannerTable+"."+fieldCreatedAt,
		bannerTable+"."+fieldUpdatedAt,
	).From(bannerTable).
		Join(
			bannerToTagsTable + " ON " + bannerTable + "." + fieldID + " = " + bannerToTagsTable + "." + fieldBannerID,
		).
		Where(
			sq.Eq{
				bannerTable + "." + fieldFeatureID:   params.FeatureID,
				bannerToTagsTable + "." + fieldTagID: params.TagID,
			},
		).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("can't create sql: %w", err))
	}

	var banner bannerDBEntity
	err = sqlx.GetContext(ctx, s.Slave(), &banner, sql, args...)
	if err != nil {
		return nil, utils.WrapSqlError(fmt.Errorf("can't get banner: %w", err))
	}
	result, err := banner.toModel()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("can't convert db entity to model: %w", err))
	}
	return result, nil
}

func (s *Storage) GetBanners(ctx context.Context, params *models.BannersParams) ([]*models.Banner, error) {
	banners := []*bannerWithTagsDBEntity{}
	err := sqlx.SelectContext(ctx, s.Slave(), &banners, queryAllBannersWithFilters, params.TagID, params.FeatureID,
		params.Offset, params.Limit)
	if err != nil {
		return nil, utils.WrapSqlError(fmt.Errorf("can't get banners: %w", err))
	}

	result := make([]*models.Banner, 0, len(banners))
	for _, banner := range banners {
		bannerModel, err := banner.toModel()
		if err != nil {
			s.storage.log.Error("can't get banner model from db entity", "error", err)
			continue
		}
		result = append(result, bannerModel)
	}
	return result, nil
}

func (s *Storage) CreateBanner(ctx context.Context, params *models.CreateBanner) (int, error) {
	tx, err := s.BeginTransaction()
	if err != nil {
		return -1, fmt.Errorf("can't begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().Truncate(time.Millisecond)
	insertBannerQuery, insertBannerArgs, err := sq.Insert(bannerTable).
		Columns(bannerFields[1:]...).
		Values(params.FeatureID, params.Content, params.IsActive, now, now).
		Suffix(returningID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return -1, utils.WrapInternalError(fmt.Errorf("can't create query: %w", err))
	}

	var bannerID int
	err = tx.GetContext(ctx, &bannerID, insertBannerQuery, insertBannerArgs...)
	if err != nil {
		return -1, utils.WrapSqlError(fmt.Errorf("can't insert banner into %s: %w", bannerTable, err))
	}

	builder := sq.Insert(bannerToTagsTable).Columns(bannerToTagsFields...)
	for _, tgID := range params.TagIDS {
		builder = builder.Values(bannerID, tgID)
	}
	insertRelatedTagsQuery, insertRelatedTagsArgs, err2 := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err2 != nil {
		return -1, utils.WrapSqlError(fmt.Errorf("can't create query: %w", err))
	}

	_, err = tx.ExecContext(ctx, insertRelatedTagsQuery, insertRelatedTagsArgs...)
	if err != nil {
		return -1, fmt.Errorf("can't insert values into %s: %w",
			bannerToTagsTable, err)
	}

	err = tx.Commit()
	if err != nil {
		return -1, utils.WrapSqlError(fmt.Errorf("can't commit transactions: %w", err))
	}
	return bannerID, nil
}

func (s *Storage) UpdateBanner(ctx context.Context, id int, params *models.CreateBanner) error {
	tx, err := s.BeginTransaction()
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't begin transaction: %w", err))
	}
	defer tx.Rollback()

	now := time.Now().Truncate(time.Millisecond)
	updateBannerQuery, updateBannerArgs, err := sq.Update(bannerTable).
		Set(fieldFeatureID, params.FeatureID).
		Set(fieldContent, params.Content).
		Set(fieldIsActive, params.IsActive).
		Set(fieldUpdatedAt, now).
		Where(sq.Eq{
			fieldID: id,
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't create updateBannerQuery: %w", err))
	}

	_, err = tx.ExecContext(ctx, updateBannerQuery, updateBannerArgs...)
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't update banner in %s: %w", bannerTable, err))
	}

	deleteOldTagsQuery, deleteOldTagsArgs, err2 := sq.Delete(bannerToTagsTable).Where(sq.Eq{fieldBannerID: id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err2 != nil {
		return utils.WrapSqlError(fmt.Errorf("can't create updateBannerQuery: %w", err))
	}
	_, err = tx.ExecContext(ctx, deleteOldTagsQuery, deleteOldTagsArgs...)
	if err != nil {
		return utils.WrapSqlError(fmt.Errorf("can't delete old relationships: %w", err))
	}

	builder := sq.Insert(bannerToTagsTable).Columns(bannerToTagsFields...)
	for _, tgID := range params.TagIDS {
		builder = builder.Values(id, tgID)
	}
	insertTagQuery, insertTagArgs, err3 := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err3 != nil {
		return utils.WrapInternalError(fmt.Errorf("can't create updateBannerQuery: %w", err3))
	}

	_, err = tx.ExecContext(ctx, insertTagQuery, insertTagArgs...)
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't insert banner_id and tag_id into %s: %w",
			bannerToTagsTable, err))
	}

	err = tx.Commit()
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't commit transactions: %w", err))
	}
	return nil
}

func (s *Storage) DeleteBanner(ctx context.Context, bannerID int) error {
	tx, err := s.BeginTransaction()
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't begin transaction: %w", err))
	}
	defer tx.Rollback()

	deleteTagQuery, tagQueryArgs, err := sq.Delete(bannerToTagsTable).Where(sq.Eq{
		fieldBannerID: bannerID,
	}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't create sql: %w", err))
	}

	_, err = tx.ExecContext(ctx, deleteTagQuery, tagQueryArgs...)
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't delete related tags: %w", err))
	}

	deleteBannerQuery, deleteBannerArgs, err2 := sq.Delete(bannerTable).Where(sq.Eq{fieldID: bannerID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err2 != nil {
		return utils.WrapInternalError(fmt.Errorf("can't create sql: %w", err2))
	}

	_, err = tx.ExecContext(ctx, deleteBannerQuery, deleteBannerArgs...)
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't delete banner: %w", err))
	}

	err = tx.Commit()
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("can't commit transaction: %w", err))
	}

	return nil
}

type bannerDBEntity struct {
	ID        int       `db:"id"`
	FeatureID int       `db:"feature_id"`
	Content   []byte    `db:"content"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (b bannerDBEntity) toModel() (*models.Banner, error) {
	var content models.Content
	err := json.Unmarshal(b.Content, &content)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal json: %w", err)
	}
	return &models.Banner{
		ID:        b.ID,
		FeatureID: b.FeatureID,
		Content:   content,
		IsActive:  b.IsActive,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}, nil
}

type bannerWithTagsDBEntity struct {
	ID        int       `db:"id"`
	TagIDs    string    `db:"tag_ids"`
	FeatureID int       `db:"feature_id"`
	Content   []byte    `db:"content"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (b *bannerWithTagsDBEntity) toModel() (*models.Banner, error) {
	tagIDs := b.parseIntArray(b.TagIDs)
	var content models.Content
	err := json.Unmarshal(b.Content, &content)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal json: %w", err)
	}
	return &models.Banner{
		ID:        b.ID,
		TagIDs:    tagIDs,
		FeatureID: b.FeatureID,
		Content:   content,
		IsActive:  b.IsActive,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}, nil
}

func (b *bannerWithTagsDBEntity) parseIntArray(input string) []int {
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(input, -1)

	result := make([]int, 0)
	for _, match := range matches {
		num, err := strconv.Atoi(match)
		if err != nil {
			log.Println("Error parsing integer:", err)
			continue
		}
		result = append(result, num)
	}

	return result
}
