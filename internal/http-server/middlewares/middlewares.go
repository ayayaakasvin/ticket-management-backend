package middlewares

import (
	"log/slog"
	"strings"

	genericjwtservice "github.com/ayayaakasvin/generic-jwt-service"
	"github.com/ayayaakasvin/oneflick-ticket/internal/config"
	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
)

type Middlewares struct {
	logger *slog.Logger

	// CORS config
	allowedOrigins   string
	allowedMethods   string
	allowedHeaders   string
	allowCredentials bool

	cc   domain.Cache
	jwtM *genericjwtservice.JWTService
	ur   domain.UserRepository
	rlm *domain.RateLimiter
}

func NewHTTPMiddlewares(logger *slog.Logger, cc domain.Cache, jwtM *genericjwtservice.JWTService, corsCfg *config.CorsConfig, ur domain.UserRepository, rlm *domain.RateLimiter) *Middlewares {
	return &Middlewares{
		logger: logger,
		cc:     cc,
		jwtM:   jwtM,
		ur: ur,

		allowedOrigins:   strings.Join(corsCfg.AllowedOrigins, ","),
		allowedMethods:   strings.Join(corsCfg.AllowedMethods, ","),
		allowedHeaders:   strings.Join(corsCfg.AllowedHeaders, ","),
		allowCredentials: corsCfg.AllowedCredentials,
	}
}
