package middlewares

import (
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	cache 	inner.Cache
	logger 	*logrus.Logger

	rlm 	*inner.RateLimiter
}

func NewHTTPMiddlewares(logger *logrus.Logger, cache inner.Cache, rlm *inner.RateLimiter) *Middlewares {
	return &Middlewares{
		logger: logger,
		cache: cache,
		rlm: rlm,
	}
}
