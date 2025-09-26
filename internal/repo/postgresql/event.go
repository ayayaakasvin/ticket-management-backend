package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
)

// CRUD for event model in models.event

// Insert Event model to the events table
func (p *PostgreSQL) InsertEvent(ctx context.Context, tx *sql.Tx, eventObj *models.Event) (string, error) {
	var newEventUUID string
	err := tx.QueryRowContext(ctx, `
		INSERT INTO events (event_uuid, starting_time, ending_time, title, description, category_id, status, capacity, organizer_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING event_uuid
		`,
		eventObj.EventUUID,
		eventObj.StartingTime,
		eventObj.EndingTime,
		eventObj.Title,
		eventObj.Description,
		eventObj.CategoryID,
		eventObj.Status,
		eventObj.Capacity,
		eventObj.OrganizerID,
	).Scan(&newEventUUID)
	if err != nil {
		return "", err
	}

	return newEventUUID, nil
}

func (p *PostgreSQL) GetEventByUUID(ctx context.Context, eventUUID string) (*models.Event, error) {
	var event *models.Event = new(models.Event)
	err := p.conn.QueryRowContext(ctx, `
		SELECT 
			event_uuid, 
			creation_time, 
			starting_time, 
			ending_time, 
			title, 
			description,
			category_id, 
			status, 
			capacity, 
			image_url, 
			organizer_id
		FROM events
		WHERE event_uuid = $1
	`, eventUUID).Scan(
		&event.EventUUID,
		&event.CreationTime,
		&event.StartingTime,
		&event.EndingTime,
		&event.Title,
		&event.Description,
		&event.CategoryID,
		&event.Status,
		&event.Capacity,
		&event.ImageURL,
		&event.OrganizerID,
	)
	if err != nil {
		return nil, err
	}
	
	if tickets, err := p.GetEventTickets(ctx, eventUUID); err != nil {
		return nil, err
	} else {
		event.Tickets = tickets
	}

	if location, err := p.GetEventLocationByEventUUID(ctx, eventUUID); err != nil {
		return  nil, err
	} else {
		event.Location = *location
	}

	return event, nil
}

// Get all events records without details such as Tickets and Location
func (p *PostgreSQL) GetAllEvents(ctx context.Context) ([]*models.Event, error) {
	rows, err := p.conn.QueryContext(ctx, `
		SELECT event_uuid, creation_time, starting_time, ending_time, title, description, category_id, status, capacity, image_url, organizer_id
		FROM events`)
	if err != nil {
		return nil, err
	}

	var events []*models.Event
	for rows.Next() {
		var event *models.Event = new(models.Event)
		err := rows.Scan(
			&event.EventUUID,
			&event.CreationTime,
			&event.StartingTime,
			&event.EndingTime,
			&event.Title,
			&event.Description,
			&event.CategoryID,
			&event.Status,
			&event.Capacity,
			&event.ImageURL,
			&event.OrganizerID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}

	return  events, nil
}

func (p *PostgreSQL) GetEventsByCategoryID(ctx context.Context, categoryID uint) ([]*models.Event, error) {
	rows, err := p.conn.QueryContext(ctx, `
		SELECT event_uuid, creation_time, starting_time, ending_time, title, description, category_id, status, capacity, image_url, organizer_id
		FROM events
		WHERE category_id = $1`, categoryID)
	if err != nil {
		return nil, err
	}

	var events []*models.Event
	for rows.Next() {
		var event *models.Event = new(models.Event)
		err := rows.Scan(
			&event.EventUUID,
			&event.CreationTime,
			&event.StartingTime,
			&event.EndingTime,
			&event.Title,
			&event.Description,
			&event.CategoryID,
			&event.Status,
			&event.Capacity,
			&event.ImageURL,
			&event.OrganizerID,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}

	return  events, nil
}

// Update the image_url for a specific event
func (p *PostgreSQL) UpdateEventImageURL(ctx context.Context, eventUUID string, imageURL string) error {
	_, err := p.conn.ExecContext(ctx, `
		UPDATE events
		SET image_url = $1
		WHERE event_uuid = $2
	`, imageURL, eventUUID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) DeleteEventByUUID(ctx context.Context, eventUUID string) error {
	_, err := p.conn.ExecContext(ctx, `
		DELETE FROM events
		WHERE event_uuid = $1
	`, eventUUID)
	if err != nil {
		return err
	}

	return nil
}
