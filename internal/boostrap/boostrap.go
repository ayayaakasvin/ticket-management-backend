package boostrap

import (
	"context"
	"log/slog"

	genericjwtservice "github.com/ayayaakasvin/generic-jwt-service"
	"github.com/ayayaakasvin/goroutinesupervisor"
	"github.com/golang-jwt/jwt/v5"
)

func SetupSupervisor(ctx context.Context, log *slog.Logger) *goroutinesupervisor.GoRoutineSupervisor {
	gs := goroutinesupervisor.NewSupervisor(ctx)
	gs.WithHandler(func(e goroutinesupervisor.Event) {
		switch e.Type {
		case goroutinesupervisor.EventTaskStarted:
			log.Info("Task started", "task", e.Task, "time", e.Started.String())
		case goroutinesupervisor.EventTaskFinished:
			log.Info("Task finished", "task", e.Task, "time", e.Ended.String())
		case goroutinesupervisor.EventTaskFailed:
			log.Info("Task failed", "task", e.Task, "time", e.Ended.String())
		default:
		}
	})

	return gs
}

func JWTManagerWithHS256(secret []byte) *genericjwtservice.JWTService {
	return genericjwtservice.NewJWTService(secret, jwt.SigningMethodHS256)
}