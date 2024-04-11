package storage

import (
	"context"
	"errors"
	"users/internal/models"
)

var (
	ErrUserExist = errors.New("user already exist")
)

type Storage interface {
	UserRepository
	TokenRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, params *models.UserRegister) (*models.User, error)
}

type TokenRepository interface {
	CreateToken(ctx context.Context, params *models.Token) (*models.Token, error)
	VerifyToken(ctx context.Context, token string) (bool, error)
	IsAdmin(ctx context.Context, token string) (bool, error)
}
