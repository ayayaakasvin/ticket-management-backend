package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ayayaakasvin/oneflick-ticket/internal/lib/bcrypthashing"
)

func (p *PostgreSQL) RegisterUser(ctx context.Context, username, hashedPassword, email string) error {
	if exists, err := p.UsernameExists(ctx, username); err != nil {
		return err
	} else if exists {
		return errors.New("user already exists")
	}

	if exists, err := p.EmailExists(ctx, username); err != nil {
		return err
	} else if exists {
		return errors.New("email already in use")
	}

	_, err := p.conn.ExecContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)",username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) AuthentificateUser(ctx context.Context, username, password string) (uint, error) {
	var (
		userId         uint
		hashedPassword string
	)
	
	if err := p.conn.QueryRowContext(ctx, "SELECT user_id, password FROM users WHERE username = $1", username).Scan(&userId, &hashedPassword); 
		err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New(NotFound)
		}
		return 0, err
	}

	if err := bcrypthashing.ComparePasswordAndHash(password, hashedPassword); err != nil {
		return 0, errors.New(UnAuthorized)
	}

	return userId, nil
}

func (p *PostgreSQL) UsernameExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := p.conn.QueryRowContext(ctx,"SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)", name).Scan(&exists)

	return exists, err
}

func (p *PostgreSQL) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := p.conn.QueryRowContext(ctx,"SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)

	return exists, err
}