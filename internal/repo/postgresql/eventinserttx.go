package postgresql

import (
	"context"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
	"github.com/google/uuid"
)

// Whole transaction function to add event. should return uuid of event or error in case of issue
func (p *PostgreSQL) InsertEventObjectToDatabase(ctx context.Context, event *models.Event) (string, error) {
	tx, err := p.conn.BeginTx(ctx, nil)
	defer func() {
    	if err != nil {
			err = fmt.Errorf("event add tx error: %v", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("tx rollback error: %v\n", rbErr)
			}
    	}
	}()
	
	if err != nil {
		return "", err
	}

	newEventUUID := uuid.NewString()
	event.EventUUID = newEventUUID
	_, err = p.InsertEvent(ctx, tx, event)
	if err != nil {
		return "", err
	}

	for _, ticket := range event.Tickets {
		newTicketUUID := uuid.NewString()
		ticket.TicketUUID = newTicketUUID
		ticket.EventUUID = newEventUUID

		err = p.InsertTicket(ctx, tx, ticket)
		if err != nil {
			return "", err
		}
	}

	event.Location.EventUUID = newEventUUID
	err = p.InsertEventLocation(ctx, tx, event.Location)
	if err != nil {
		return "", err
	}

	err = p.InsertTags(ctx, tx, newEventUUID, event.Tags)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("commit error: %v", err)
	}

	return newEventUUID, nil
}