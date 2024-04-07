package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"users/internal/models"
)

const ExpirationDuration = time.Hour * 24

func GenerateJWTToken(user *models.User) (string, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if user == nil {
		return "", fmt.Errorf("user cannot be empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": user.Username,
	})

	s, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("signed string: %w", err)
	}

	return s, nil
}

func GetExpirationDate() time.Time {
	return time.Now().Add(ExpirationDuration)
}
