package middlewares

import (
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
	"github.com/sirupsen/logrus"
)

type Middlewares struct {
	cache 	inner.Cache
	logger *logrus.Logger
}

func NewHTTPMiddlewares(logger *logrus.Logger, cache inner.Cache) *Middlewares {
	return &Middlewares{
		logger: logger,
		cache: cache,
	}
}
