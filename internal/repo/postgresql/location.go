package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
)

// Insert Ticket model to the tickets table
func (p *PostgreSQL) InsertEventLocation(ctx context.Context, tx *sql.Tx, location models.Location) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO locations (event_uuid, name, address, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5)
	`,
		location.EventUUID,
		location.Name,
		location.Address,
		location.Latitude,
		location.Longitude,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) GetEventLocation(ctx context.Context, locationID uint) (*models.Location, error) {
	var location *models.Location = new(models.Location)
	err := p.conn.QueryRowContext(ctx, `
		SELECT event_uuid, name, address, latitude, longitude
		FROM locations
		WHERE location_id = $1
	`, locationID).Scan(&location.EventUUID, &location.Name, &location.Address, &location.Latitude, &location.Longitude)
	if err != nil {
		return nil, err
	}

	location.LocationID = locationID
	return location, nil
}

func (p *PostgreSQL) GetEventLocationByEventUUID(ctx context.Context, eventUUID string) ([]*models.Location, error) {
	rows, err := p.conn.QueryContext(ctx, `
		SELECT location_id, event_uuid, name, address, latitude, longitude
		FROM locations
		WHERE event_uuid = $1
	`, eventUUID)
	if err != nil {
		return nil, err
	}

	var locationsOfEvent []*models.Location
	for rows.Next() {
		var location *models.Location = new(models.Location)
		err := rows.Scan(&location.LocationID, &location.EventUUID, &location.Name, &location.Address, &location.Latitude, &location.Longitude)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		locationsOfEvent = append(locationsOfEvent, location)
	}
	
	if rows.Err() != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}
	
	return locationsOfEvent, nil
}

func (p *PostgreSQL) UpdateEventLocation(ctx context.Context, eventUUID string, location models.Location) error {
	_, err := p.conn.ExecContext(ctx, `
		UPDATE locations
		SET name = $1, address = $2, latitude = $3, longitude = $4
		WHERE event_uuid = $5
	`, location.Name, location.Address, location.Latitude, location.Longitude, eventUUID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreSQL) DeleteEventLocation(ctx context.Context, locationID uint) error {
	_, err := p.conn.ExecContext(ctx, `
		DELETE FROM locations
		WHERE location_id = $1
	`, locationID)
	if err != nil {
		return err
	}

	return nil
}