package storage

import (
	"context"
	"users/internal/models"
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
