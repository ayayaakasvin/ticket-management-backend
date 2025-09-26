// Handlers that serves for main http server, accessed via handlerd.Handler struct that contains necessary dependencies
package handlers

import (
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	userRepo 	inner.UserRepository
	eventRepo 	inner.EventRepository
	cache 		inner.Cache
	lfs			inner.FS

	logger 		*logrus.Logger
}

func NewHTTPHandlers(user inner.UserRepository, eventRepo inner.EventRepository, cache inner.Cache, logger *logrus.Logger, lfs inner.FS) *Handlers {
	return &Handlers{
		userRepo: user,
		eventRepo: eventRepo,
		cache: cache,
		lfs: lfs,

		logger: logger,
	}
}