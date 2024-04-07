package storage

import (
	"context"
	"users/internal/models"
)

type Storage interface {
	UserRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, params *models.UserRegister) (*models.User, error)
}
