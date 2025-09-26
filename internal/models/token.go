package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	UserID    uint    `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}
