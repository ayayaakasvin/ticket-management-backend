package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type EventMatch struct {
	EventUUID 		string
	MatchPercentage float64
}

// Insert tags to the tags table
func (p *PostgreSQL) InsertTags(ctx context.Context, tx *sql.Tx, eventUUID string, tags []string) error {
	for _, tag := range tags {
		var tagID int
		err := tx.QueryRowContext(ctx, `
			INSERT INTO tags (name)
			VALUES ($1)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING tag_id
		`, strings.ToLower(tag)).Scan(&tagID)
		if err != nil {
			return fmt.Errorf("failed to insert tag %q", tag)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO event_tags (event_uuid, tag_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, eventUUID, tagID)
		if err != nil {
			return fmt.Errorf("failed to link tag %q to event %s: %w", tag, eventUUID, err)
		}
	}
	return nil
}

func (p *PostgreSQL) GetEventsByTags(ctx context.Context, tags []string) ([]*EventMatch, error) {
	rows, err := p.conn.QueryContext(ctx, `
		SELECT 
			et.event_uuid,
			COUNT(DISTINCT et.tag_id)::float / $1 AS match_percentage
		FROM 
			event_tags et
		WHERE 
			et.tag_id = ANY ($2)
		GROUP BY 
			et.event_uuid
		ORDER BY 
			match_percentage DESC
		LIMIT 10;
	`, len(tags), fmt.Sprintf("ARRAY[%s]", strings.Join(tags, ", ")),
	)
	if err != nil {
		return nil, err
	}

	var eventMatched []*EventMatch
	for rows.Next() {
		var matched *EventMatch = new(EventMatch)
		err := rows.Scan(&matched.EventUUID, &matched.MatchPercentage)
		if err != nil {
			return nil, err
		}

		eventMatched = append(eventMatched, matched)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return eventMatched, nil
}

func (p *PostgreSQL) DeleteTag(ctx context.Context, tagID uint) error {
	_, err := p.conn.ExecContext(ctx, `
		DELETE FROM tags
		WHERE tag_id = $1
	`, tagID)
	if err != nil {
		return err
	}

	return nil
}