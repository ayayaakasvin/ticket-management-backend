// Handlers that serves for main http server, accessed via handlerd.Handler struct that contains necessary dependencies
package handlers

import (
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	userRepo 	inner.UserRepository
	cache 		inner.Cache

	logger 		*logrus.Logger
}

func NewHTTPHandlers(user inner.UserRepository, cache inner.Cache, logger *logrus.Logger) *Handlers {
	return &Handlers{
		userRepo: user,

		logger: logger,
	}
}