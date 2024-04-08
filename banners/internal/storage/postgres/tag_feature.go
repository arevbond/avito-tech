package postgres

import (
	"banners/internal/models"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) CreateTag(ctx context.Context, params *models.Tag) error {
	sql, args, err := sq.Insert(tagTable).
		Columns(tagFields...).
		Values(params.ID, params.Name).
		PlaceholderFormat(sq.Dollar).
		Suffix(suffixDoNothing).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't build sql query: %w", err)
	}

	_, err = s.Master().ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("can't insert tag into %s: %w", tagTable, err)
	}
	return nil
}

func (s *Storage) CreateFeature(ctx context.Context, params *models.Feature) error {
	sql, args, err := sq.Insert(featureTable).
		Columns(featureFields...).
		Values(params.ID, params.Name).
		PlaceholderFormat(sq.Dollar).
		Suffix(suffixDoNothing).
		ToSql()
	if err != nil {
		return fmt.Errorf("can't build sql query: %w", err)
	}

	_, err = s.Master().ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("can't insert tag into %s: %w", featureTable, err)
	}
	return nil
}
