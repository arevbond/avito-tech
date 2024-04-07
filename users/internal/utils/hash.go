package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password string) (string, error) {
	h := sha256.New()
	h.Write([]byte(password))
	resultHash := h.Sum(nil)
	resultString := hex.EncodeToString(resultHash)
	return resultString, nil
}
