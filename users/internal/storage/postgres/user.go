package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"users/internal/models"
	"users/internal/storage"
	"users/internal/utils"
)

func (s *Storage) CreateUser(ctx context.Context, params *models.UserRegister) (*models.User, error) {
	sql, args, err := sq.Insert(userTable).
		Columns(userFields...).
		Values(params.ID, params.Username, params.HashedPassword, params.IsAdmin).
		Suffix(returningUser).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("incorect sql: %w", err))
	}

	var entity userDBEntity
	err = s.db.GetContext(ctx, &entity, sql, args...)
	if err != nil {
		pgErr, ok := err.(*pq.Error)
		if ok && pgErr.Code == "23505" {
			return nil, storage.ErrUserExist
		}
		return nil, utils.WrapSqlError(err)
	}

	return entity.convertToModel(), nil
}

type userDBEntity struct {
	ID             string `db:"id"`
	Username       string `db:"username"`
	HashedPassword string `db:"hashed_password"`
	IsAdmin        bool   `db:"is_admin"`
}

func (u userDBEntity) convertToModel() *models.User {
	return &models.User{
		ID:       models.UserID(u.ID),
		Username: u.Username,
		IsAdmin:  u.IsAdmin,
	}
}
