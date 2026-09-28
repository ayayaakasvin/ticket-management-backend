package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
	"github.com/google/uuid"
)

const eventColumns = `event_uuid, creation_time, starting_time, ending_time, title,
description, category_id, status, capacity, image_url, organizer_id`

type scanner interface {
	Scan(dest ...any) error
}

func scanEvent(row scanner) (*domain.Event, error) {
	event := new(domain.Event)
	var description, imageURL sql.NullString
	var categoryID, capacity, organizerID sql.NullInt64
	err := row.Scan(
		&event.EventUUID,
		&event.CreationTime,
		&event.StartingTime,
		&event.EndingTime,
		&event.Title,
		&description,
		&categoryID,
		&event.Status,
		&capacity,
		&imageURL,
		&organizerID,
	)
	if err != nil {
		return nil, err
	}
	event.Description = description.String
	event.ImageURL = imageURL.String
	if categoryID.Valid {
		event.CategoryID = uint(categoryID.Int64)
	}
	if capacity.Valid {
		event.Capacity = uint(capacity.Int64)
	}
	if organizerID.Valid {
		event.OrganizerID = uint(organizerID.Int64)
	}
	return event, nil
}

func (s *SQLite) DeleteEvent(ctx context.Context, eventUUID string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM events WHERE event_uuid = ?`, eventUUID)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	if n, err := result.RowsAffected(); err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *SQLite) GetEvents(ctx context.Context) ([]*domain.Event, error) {
	return s.getEvents(ctx, `SELECT `+eventColumns+` FROM events ORDER BY starting_time`)
}

func (s *SQLite) GetEventsByCategory(ctx context.Context, categoryID uint) ([]*domain.Event, error) {
	return s.getEvents(ctx, `SELECT `+eventColumns+` FROM events WHERE category_id = ? ORDER BY starting_time`, categoryID)
}

func (s *SQLite) getEvents(ctx context.Context, query string, args ...any) ([]*domain.Event, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	events := make([]*domain.Event, 0)
	for rows.Next() {
		event, err := scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}
	return events, nil
}

func (s *SQLite) GetEvent(ctx context.Context, eventUUID string) (*domain.Event, error) {
	event, err := scanEvent(s.db.QueryRowContext(ctx,
		`SELECT `+eventColumns+` FROM events WHERE event_uuid = ?`, eventUUID,
	))
	if err != nil {
		return nil, err
	}

	event.Tickets, err = s.getEventTickets(ctx, eventUUID)
	if err != nil {
		return nil, err
	}
	var location domain.Location
	err = s.db.QueryRowContext(ctx, `
		SELECT location_id, event_uuid, name, address, latitude, longitude
		FROM locations WHERE event_uuid = ? LIMIT 1`, eventUUID,
	).Scan(&location.LocationID, &location.EventUUID, &location.Name, &location.Address, &location.Latitude, &location.Longitude)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("get event location: %w", err)
	}
	if err == nil {
		event.Location = location
	}
	event.Tags, err = s.getEventTags(ctx, eventUUID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *SQLite) InsertEventTx(ctx context.Context, tx *sql.Tx, eventObj *domain.Event) (string, error) {
	if eventObj == nil {
		return "", errors.New("event is nil")
	}
	if tx == nil {
		return "", errors.New("transaction is nil")
	}
	if eventObj.EventUUID == "" {
		eventObj.EventUUID = uuid.NewString()
	}
	var categoryID, capacity, organizerID any
	if eventObj.CategoryID != 0 {
		categoryID = eventObj.CategoryID
	}
	if eventObj.Capacity != 0 {
		capacity = eventObj.Capacity
	}
	if eventObj.OrganizerID != 0 {
		organizerID = eventObj.OrganizerID
	}
	_, err := tx.ExecContext(ctx, `
		INSERT INTO events (
			event_uuid, starting_time, ending_time, title, description,
			category_id, status, capacity, image_url, organizer_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		eventObj.EventUUID,
		eventObj.StartingTime,
		eventObj.EndingTime,
		eventObj.Title,
		eventObj.Description,
		categoryID,
		eventObj.Status,
		capacity,
		nullString(eventObj.ImageURL),
		organizerID,
	)
	if err != nil {
		return "", fmt.Errorf("insert event: %w", err)
	}
	return eventObj.EventUUID, nil
}

func (s *SQLite) CreateEvent(ctx context.Context, event *domain.Event) (string, error) {
	if event == nil {
		return "", errors.New("event is nil")
	}
	if event.EventUUID == "" {
		event.EventUUID = uuid.NewString()
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin create event transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := s.InsertEventTx(ctx, tx, event); err != nil {
		return "", err
	}
	for _, ticket := range event.Tickets {
		if ticket == nil {
			return "", errors.New("event contains a nil ticket")
		}
		if ticket.TicketUUID == "" {
			ticket.TicketUUID = uuid.NewString()
		}
		ticket.EventUUID = event.EventUUID
		if err := insertTicket(ctx, tx, ticket); err != nil {
			return "", err
		}
	}
	event.Location.EventUUID = event.EventUUID
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO locations (event_uuid, name, address, latitude, longitude)
		VALUES (?, ?, ?, ?, ?)`,
		event.Location.EventUUID,
		event.Location.Name,
		event.Location.Address,
		event.Location.Latitude,
		event.Location.Longitude,
	); err != nil {
		return "", fmt.Errorf("insert event location: %w", err)
	}
	for _, tag := range event.Tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO tags(name) VALUES (?)`, tag); err != nil {
			return "", fmt.Errorf("insert event tag: %w", err)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT OR IGNORE INTO event_tags(event_uuid, tag_id)
			SELECT ?, tag_id FROM tags WHERE name = ?`, event.EventUUID, tag); err != nil {
			return "", fmt.Errorf("link event tag: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit create event transaction: %w", err)
	}
	return event.EventUUID, nil
}

func (s *SQLite) UpdateEventImage(ctx context.Context, eventUUID string, imageURL string) error {
	result, err := s.db.ExecContext(ctx, `UPDATE events SET image_url = ? WHERE event_uuid = ?`, imageURL, eventUUID)
	if err != nil {
		return fmt.Errorf("update event image: %w", err)
	}
	if n, err := result.RowsAffected(); err == nil && n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func nullString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func (s *SQLite) getEventTickets(ctx context.Context, eventUUID string) ([]*domain.Ticket, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT ticket_uuid, event_uuid, name, price, currency, quantity, sold
		FROM tickets WHERE event_uuid = ? ORDER BY rowid`, eventUUID)
	if err != nil {
		return nil, fmt.Errorf("query event tickets: %w", err)
	}
	defer rows.Close()
	tickets := make([]*domain.Ticket, 0)
	for rows.Next() {
		var ticket domain.Ticket
		var quantity, sold sql.NullInt64
		var currency sql.NullString
		if err := rows.Scan(&ticket.TicketUUID, &ticket.EventUUID, &ticket.Name, &ticket.Price, &currency, &quantity, &sold); err != nil {
			return nil, fmt.Errorf("scan event ticket: %w", err)
		}
		ticket.Currency = currency.String
		if quantity.Valid {
			ticket.Quantity = uint(quantity.Int64)
		}
		if sold.Valid {
			ticket.Sold = uint(sold.Int64)
		}
		tickets = append(tickets, &ticket)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate event tickets: %w", err)
	}
	return tickets, nil
}

func (s *SQLite) getEventTags(ctx context.Context, eventUUID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT t.name FROM tags t
		JOIN event_tags et ON et.tag_id = t.tag_id
		WHERE et.event_uuid = ? ORDER BY t.name`, eventUUID)
	if err != nil {
		return nil, fmt.Errorf("query event tags: %w", err)
	}
	defer rows.Close()
	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("scan event tag: %w", err)
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate event tags: %w", err)
	}
	return tags, nil
}
