package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
)

// Insert Ticket model to the tickets table
func (p *PostgreSQL) InsertTicket(ctx context.Context, tx *sql.Tx, ticketObj *models.Ticket) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO tickets (ticket_uuid, event_uuid, name, price, currency, quantity)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		ticketObj.TicketUUID,
		ticketObj.EventUUID,
		ticketObj.Name,
		ticketObj.Price,
		ticketObj.Currency,
		ticketObj.Quantity,
	)
	if err != nil {
		return err
	}

	return nil
}

// Modify Ticket record to decrement the quantity, rn without the payment part maybe it will be completed by *sql.Tx
func (p *PostgreSQL) UpdateTicketWasSolt(ctx context.Context, ticketUUID string) error {
	var updatedID string
	err := p.conn.QueryRowContext(ctx, `
        UPDATE tickets
        SET 
            quantity = CASE 
                WHEN quantity IS NULL THEN NULL
                WHEN quantity > 0 THEN quantity - 1
                ELSE quantity
            END,
            sold = sold + 1
        WHERE ticket_uuid = $1
          AND (quantity IS NULL OR quantity > 0)
        RETURNING ticket_uuid
    `, ticketUUID).Scan(&updatedID)

	if err == sql.ErrNoRows {
		return errors.New(NotFound)
	}
	if err != nil {
		return err
	}

	return nil
}

// Modify Ticket record to decrement the quantity, rn without the payment part maybe it will be completed by *sql.Tx
func (p *PostgreSQL) DeleteTicket(ctx context.Context, ticketUUID string) error {
	var updatedID string
	err := p.conn.QueryRowContext(ctx, `
		DELETE FROM tickets WHERE ticket_uuid = $1
		RETURNING ticket_uuid
	`, ticketUUID).Scan(&updatedID)

	if err == sql.ErrNoRows {
		return errors.New(NotFound)
	}
	if err != nil {
		return err
	}

	return nil
}

// Get the record of Ticket
func (p *PostgreSQL) GetTicket(ctx context.Context, ticketUUID string) (*models.Ticket, error) {
	var ticketObj *models.Ticket = new(models.Ticket)
	err := p.conn.QueryRowContext(ctx, `
		SELECT ticket_uuid, event_uuid, name, price, currency, quantity, sold
		FROM tickets
		WHERE ticket_uuid = $1
	`, ticketUUID).Scan(
		&ticketObj.TicketUUID,
		&ticketObj.EventUUID,
		&ticketObj.Name,
		&ticketObj.Price,
		&ticketObj.Currency,
		&ticketObj.Quantity,
		&ticketObj.Sold,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(NotFound)
	}
	if err != nil {
		return nil, err
	}

	return ticketObj, nil
}

func (p *PostgreSQL) GetEventTickets(ctx context.Context, eventUUID string) ([]*models.Ticket, error) {
	rows, err := p.conn.QueryContext(ctx, `
		SELECT ticket_uuid, event_uuid, name, price, currency, quantity, sold
		FROM tickets
		WHERE event_uuid = $1
	`, eventUUID,
	)
	if err != nil {
		return nil, err
	}

	var tickets []*models.Ticket
	for rows.Next() {
		var ticketObj *models.Ticket = new(models.Ticket)
		err := rows.Scan(
			&ticketObj.TicketUUID,
			&ticketObj.EventUUID,
			&ticketObj.Name,
			&ticketObj.Price,
			&ticketObj.Currency,
			&ticketObj.Quantity,
			&ticketObj.Sold,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		tickets = append(tickets, ticketObj)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}

	return tickets, nil
}

// Insert Ticket model to the tickets table after main transaction
func (p *PostgreSQL) InsertTicketAfterwards(ctx context.Context, ticketObj *models.Ticket) error {
	_, err := p.conn.ExecContext(ctx, `
		INSERT INTO tickets (ticket_uuid, event_uuid, name, price, currency, quantity)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		ticketObj.TicketUUID,
		ticketObj.EventUUID,
		ticketObj.Name,
		ticketObj.Price,
		ticketObj.Currency,
		ticketObj.Quantity,
	)
	if err != nil {
		return err
	}

	return nil
}