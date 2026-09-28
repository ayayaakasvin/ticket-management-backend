package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ayayaakasvin/oneflick-ticket/docs"
	"github.com/ayayaakasvin/oneflick-ticket/internal/boostrap"
	"github.com/ayayaakasvin/oneflick-ticket/internal/config"
	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
	"github.com/ayayaakasvin/oneflick-ticket/internal/fs"
	"github.com/ayayaakasvin/oneflick-ticket/internal/repo/gcache"
	"github.com/ayayaakasvin/oneflick-ticket/internal/repo/sqlite"
	"github.com/ayayaakasvin/oneflick-ticket/internal/smtptool"

	httpserver "github.com/ayayaakasvin/oneflick-ticket/internal/http-server"
	"github.com/ayayaakasvin/oneflick-ticket/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoadConfig()
	log := logger.SetupSlog(cfg.Logger)

	// smtp, err := smtptool.NewSMTPToolWithPreHealthCheck(cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host, cfg.SMTP.Port)
	// if err != nil {
	// 	return err
	// }

	smtp, err := smtptool.NewSMTP_MockUp()
	if err != nil {
		return err
	}

	lfs, err := fs.NewFS(cfg.LFS.BasePath)
	if err != nil {
		return err
	}

	repo, err := sqlite.NewSqliteConnection(cfg.Database.FilePath)
	if err != nil {
		return err
	}

	gc := gcache.NewGCache(cfg.JWT.AccessTokenTTL)
	
	rlm := domain.NewRateLimiter()

	jwtM := boostrap.JWTManagerWithHS256([]byte(cfg.JWT.Secret))

	gs := boostrap.SetupSupervisor(ctx, log)

	app := httpserver.NewServerApp(&cfg.HTTP, &cfg.CORS, log, repo, repo, gc, lfs, rlm, jwtM, smtp)
	// TODO: reformat whole code to your new standard
	gs.Go("HTTP-server", app.Start)

	err = gs.Wait()
	if err != nil {
		return fmt.Errorf("gs wait error: %s", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.Stop(shutdownCtx)

	return nil
}
