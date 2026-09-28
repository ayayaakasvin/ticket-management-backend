package token

import (
	"time"

	genericjwtservice "github.com/ayayaakasvin/generic-jwt-service"
	"github.com/golang-jwt/jwt/v5"
)

var (
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
)

func SetTTL(attl, rttl time.Duration) {
	AccessTokenTTL = attl
	RefreshTokenTTL = rttl
}

type AccessTokenClaims struct {
	UserID    int64  `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewAccessTokenClaims(uID int64, sID string, ttl time.Duration) *AccessTokenClaims {
	return &AccessTokenClaims{
		UserID:           uID,
		SessionID:        sID,
		RegisteredClaims: genericjwtservice.StdClaims(ttl),
	}
}

func NewRefreshTokenClaims(uID int64, ttl time.Duration) *RefreshTokenClaims {
	return &RefreshTokenClaims{
		UserID:           uID,
		RegisteredClaims: genericjwtservice.StdClaims(ttl),
	}
}
