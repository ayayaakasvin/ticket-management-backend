package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
	"github.com/google/uuid"
)

func (s *SQLite) GetTicket(ctx context.Context, ticketUUID string) (*domain.Ticket, error) {
	var ticket domain.Ticket
	var currency sql.NullString
	var quantity, sold sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		SELECT ticket_uuid, event_uuid, name, price, currency, quantity, sold
		FROM tickets WHERE ticket_uuid = ?`, ticketUUID,
	).Scan(&ticket.TicketUUID, &ticket.EventUUID, &ticket.Name, &ticket.Price, &currency, &quantity, &sold)
	if err != nil {
		return nil, err
	}
	ticket.Currency = currency.String
	if quantity.Valid {
		ticket.Quantity = uint(quantity.Int64)
	}
	if sold.Valid {
		ticket.Sold = uint(sold.Int64)
	}
	return &ticket, nil
}

func (s *SQLite) DeleteTicket(ctx context.Context, ticketUUID string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM tickets WHERE ticket_uuid = ?`, ticketUUID)
	if err != nil {
		return fmt.Errorf("delete ticket: %w", err)
	}
	if n, err := result.RowsAffected(); err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *SQLite) CreateTicket(ctx context.Context, ticket *domain.Ticket) error {
	if ticket == nil {
		return errors.New("ticket is nil")
	}
	if ticket.TicketUUID == "" {
		ticket.TicketUUID = uuid.NewString()
	}
	if err := insertTicket(ctx, s.db, ticket); err != nil {
		return err
	}
	return nil
}

type execContexter interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func insertTicket(ctx context.Context, exec execContexter, ticket *domain.Ticket) error {
	_, err := exec.ExecContext(ctx, `
		INSERT INTO tickets (ticket_uuid, event_uuid, name, price, currency, quantity, sold)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		ticket.TicketUUID,
		ticket.EventUUID,
		ticket.Name,
		ticket.Price,
		nullString(ticket.Currency),
		ticket.Quantity,
		ticket.Sold,
	)
	if err != nil {
		return fmt.Errorf("insert ticket: %w", err)
	}
	return nil
}
