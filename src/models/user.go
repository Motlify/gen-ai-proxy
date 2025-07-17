package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID           pgtype.UUID `json:"id"`
	Username     string      `json:"username"`
	PasswordHash string      `json:"-"`
	CreatedAt    time.Time   `json:"created_at"`
}

type JwtCustomClaims struct {
	UserID pgtype.UUID `json:"user_id"`
	jwt.RegisteredClaims
}