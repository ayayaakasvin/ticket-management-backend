package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
	_ "modernc.org/sqlite"
)

type SQLite struct {
	db *sql.DB
}

var (
	_ domain.UserRepository  = (*SQLite)(nil)
	_ domain.EventRepository = (*SQLite)(nil)
)

func NewSqliteConnection(dbPath string) (*SQLite, error) {
	if dbPath == "" {
		return nil, fmt.Errorf("database path is empty")
	}

	dbPath, err := normalizeSqliteDBPath(dbPath)
	if err != nil {
		return nil, err
	}

	if dir := filepath.Dir(dbPath); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed to create sqlite directory %s: %w", dir, err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db %s: %w", dbPath, err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	if _, err := db.ExecContext(context.Background(), "PRAGMA busy_timeout = 5000"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set sqlite busy timeout: %w", err)
	}

	if _, err := db.ExecContext(context.Background(), "PRAGMA journal_mode = WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set sqlite journal mode: %w", err)
	}

	if _, err := db.ExecContext(context.Background(), "PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable sqlite foreign keys: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping sqlite db %s: %w", dbPath, err)
	}
	return &SQLite{db: db}, nil
}

func normalizeSqliteDBPath(dbPath string) (string, error) {
	if dbPath == "~" || strings.HasPrefix(dbPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}
		if dbPath == "~" {
			dbPath = home
		} else {
			dbPath = filepath.Join(home, dbPath[2:])
		}
	}

	dbPath = filepath.Clean(dbPath)
	if !filepath.IsAbs(dbPath) {
		absPath, err := filepath.Abs(dbPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve sqlite db path %s: %w", dbPath, err)
		}
		dbPath = absPath
	}

	return dbPath, nil
}
