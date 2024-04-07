package models

import "time"

type UserRegister struct {
	ID             string
	Username       string
	HashedPassword string
	IsAdmin        bool
}

type UserLogin struct {
	Username       string
	HashedPassword string
}

type UserID string

type User struct {
	ID       UserID
	Username string
	IsAdmin  bool
}

type TokenID string

type Token struct {
	ID             TokenID
	UserID         UserID
	Value          string
	ExpirationDate time.Time
}
