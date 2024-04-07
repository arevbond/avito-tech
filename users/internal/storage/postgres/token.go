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
