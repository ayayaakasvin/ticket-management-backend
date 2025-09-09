package main

import (
	"os"
	"sync"

	"github.com/ayayaakasvin/oneflick-ticket/internal/config"
	httpserver "github.com/ayayaakasvin/oneflick-ticket/internal/http-server"
	"github.com/ayayaakasvin/oneflick-ticket/internal/logger"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
	"github.com/ayayaakasvin/oneflick-ticket/internal/repo/postgresql"
	"github.com/ayayaakasvin/oneflick-ticket/internal/repo/valkey"
)

func main() {
	cfg := config.MustLoadConfig()
	logger := logger.SetupLogger()

	shutdownChan := inner.NewShutdownChannel()
	go func() {
		logger.Errorf("Error during setup: %s, %v", shutdownChan.Value(), cfg)
		os.Exit(1)
	}()

	repo := postgresql.NewPostgreSQLConnection(cfg.Database, shutdownChan)
	logger.Info("Postgresql conn has been established")

	cache := valkey.NewValkeyClient(cfg.Valkey, shutdownChan)
	logger.Info("Valkey conn has been established")

	wg := new(sync.WaitGroup)
	wg.Add(1) // to wait for server

	app := httpserver.NewServerApp(&cfg.HTTPServer, logger, wg, repo, cache)

	app.Run()
}
