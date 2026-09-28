// Handlers that serves for main http server, accessed via handlerd.Handler struct that contains necessary dependencies
package handlers

import (
	"log/slog"

	genericjwtservice "github.com/ayayaakasvin/generic-jwt-service"
	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
)

type Handlers struct {
	ur   domain.UserRepository
	er   domain.EventRepository
	cc   domain.Cache
	lfs  domain.FS
	smtp domain.SMTP
	jwtM *genericjwtservice.JWTService

	logger *slog.Logger
}

func NewHTTPHandlers(
	user domain.UserRepository, 
	eventRepo domain.EventRepository, 
	cache domain.Cache, 
	lfs domain.FS, 
	smtp domain.SMTP, 
	jwtM *genericjwtservice.JWTService,
	logger *slog.Logger,
	) *Handlers {
	return &Handlers{
		ur:   user,
		er:   eventRepo,
		cc:   cache,
		lfs:  lfs,
		smtp: smtp,
		jwtM: jwtM,

		logger: logger,
	}
}
