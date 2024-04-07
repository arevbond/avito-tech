package service

import (
	"context"
	"fmt"
	"log/slog"
	"users/internal/models"
	"users/internal/storage"
	"users/internal/utils"
)

type Service interface {
	Register(ctx context.Context, params *RegisterParams) (*models.Token, error)
}

type UserService struct {
	db  storage.Storage
	log *slog.Logger
}

func New(db storage.Storage, log *slog.Logger) *UserService {
	return &UserService{
		db:  db,
		log: log,
	}
}

type RegisterParams struct {
	Username string
	Password string
	IsAdmin  bool
}

func (s *UserService) Register(ctx context.Context, params *RegisterParams) (*models.Token, error) {
	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		return nil, fmt.Errorf("can't hash password: %w", err)
	}

	user, err := s.db.CreateUser(ctx, &models.UserRegister{
		ID:             utils.GenerateUUID(),
		Username:       params.Username,
		HashedPassword: hashedPassword,
		IsAdmin:        params.IsAdmin,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create user: %w", err)
	}

	jwtToken, err := utils.GenerateJWTToken(user)
	if err != nil {
		return nil, fmt.Errorf("can't generate jwt token :%w", err)
	}

	token, err := s.db.CreateToken(ctx, &models.Token{
		ID:             models.TokenID(utils.GenerateUUID()),
		UserID:         user.ID,
		Value:          jwtToken,
		ExpirationDate: utils.GetExpirationDate(),
	})
	if err != nil {
		return nil, fmt.Errorf("can't create token in storage: %w", err)
	}

	return token, nil
}
