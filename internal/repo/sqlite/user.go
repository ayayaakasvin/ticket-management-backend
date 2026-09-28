package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
)

// Create inserts a user and returns its generated ID.
func (s *SQLite) Create(ctx context.Context, user *domain.User) (int64, error) {
	if user == nil {
		return 0, errors.New("user is nil")
	}
	role := user.Role
	if role == "" {
		role = domain.Client
	}
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO users (username, email, password) VALUES (?, ?, ?)`,
		user.Username, user.Email, user.PasswordHash,
	)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get inserted user ID: %w", err)
	}
	user.ID = id
	user.Role = role
	return id, nil
}

// GetByID returns nil when no matching user exists.
func (s *SQLite) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	var createdAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT user_id, username, email, password, created_at FROM users WHERE user_id = ?`, id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}
	user.Role = domain.Client
	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	return &user, nil
}

// GetByUsername returns nil when no matching user exists.
func (s *SQLite) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	var createdAt sql.NullTime
	err := s.db.QueryRowContext(ctx,
		`SELECT user_id, username, email, password, created_at FROM users WHERE username = ?`, username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	user.Role = domain.Client
	if createdAt.Valid {
		user.CreatedAt = createdAt.Time
	}
	return &user, nil
}

// IsAdmin reports false for missing users and database errors.
func (s *SQLite) IsAdmin(ctx context.Context, id int64) bool {
	// Role is not persisted by the supplied SQLite schema.
	return false
}

func (s *SQLite) Update(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	result, err := s.db.ExecContext(ctx,
		`UPDATE users SET username = ?, email = ?, password = ? WHERE user_id = ?`,
		user.Username, user.Email, user.PasswordHash, user.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if n, err := result.RowsAffected(); err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
