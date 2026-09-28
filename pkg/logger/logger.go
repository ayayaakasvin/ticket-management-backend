package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/ayayaakasvin/oneflick-ticket/internal/config"
	"github.com/sirupsen/logrus"
)

func SetupLogrus() *logrus.Logger {
	logFile, err := os.OpenFile(".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	seperator := fmt.Sprintf("=== Logging started at %s ===\n", time.Now().Format("2006-01-02 15:04:05"))
	logFile.Write([]byte(seperator))

	multiwriter := io.MultiWriter(os.Stdout, logFile)

	logger := logrus.New()
	logger.SetOutput(multiwriter)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Logger has been set up")
	return logger
}

func SetupSlog(cfg config.LoggerConfig) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.Level(cfg.Level),
	}

	var handler slog.Handler

	if cfg.JSON {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler).With(
		slog.String("service", cfg.Service),
	)

	return logger
}
