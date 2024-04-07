package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"time"
	"users/internal/models"
	"users/internal/utils"
)

func (s *Storage) CreateToken(ctx context.Context, params *models.Token) (*models.Token, error) {
	now := time.Now().Truncate(time.Millisecond)
	sql, args, err := sq.Insert(tokenTable).
		Columns(tokenFields...).
		Values(params.ID, params.UserID, params.Value, now, params.ExpirationDate).
		Suffix(returningToken).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("incorrect sql: %w", err))
	}

	var entity tokenDBEntity
	err = s.db.GetContext(ctx, &entity, sql, args...)
	if err != nil {
		return nil, utils.WrapSqlError(err)
	}
	return entity.convertToModel(), nil
}

type tokenDBEntity struct {
	ID             string    `db:"id"`
	UserID         string    `db:"user_id"`
	Value          string    `db:"value"`
	CreatedAt      time.Time `db:"created_at"`
	ExpirationDate time.Time `db:"expiration_date"`
}

func (t tokenDBEntity) convertToModel() *models.Token {
	return &models.Token{
		ID:             models.TokenID(t.ID),
		UserID:         models.UserID(t.UserID),
		Value:          t.Value,
		ExpirationDate: t.ExpirationDate,
	}
}

func (s *Storage) VerifyToken(ctx context.Context, token string) (bool, error) {
	sql, args, err := sq.Select("1").
		Prefix("SELECT EXISTS (").
		From(tokenTable).
		Where(sq.Eq{
			fieldTokenValue: token,
		}).
		Suffix(")").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, utils.WrapInternalError(fmt.Errorf("can't get token: %w", err))
	}

	var result bool
	err = s.db.GetContext(ctx, &result, sql, args...)
	if err != nil {
		return false, utils.WrapSqlError(err)
	}
	return result, nil
}

func (s *Storage) IsAdmin(ctx context.Context, token string) (bool, error) {
	sql, args, err := sq.Select("u.is_admin").
		From("tokens t").
		Join("users u ON t.user_id = u.id").
		Where(sq.Eq{"t.value": token}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return false, fmt.Errorf("sql error: %w", err)
	}

	var isAdmin bool
	err = s.db.GetContext(ctx, &isAdmin, sql, args...)
	if err != nil {
		return false, utils.WrapSqlError(err)
	}
	return isAdmin, nil
}
