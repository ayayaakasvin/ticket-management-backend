package postgresql

import (
	"context"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
	"github.com/google/uuid"
)

// Whole transaction function to add event. should return uuid of event or error in case of issue
func (p *PostgreSQL) AddEventToDatabase(ctx context.Context, event *models.Event, tickets []*models.Ticket, location models.Location, tags []string) (string, error) {
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

	for _, ticket := range tickets {
		newTicketUUID := uuid.NewString()
		ticket.TicketUUID = newTicketUUID
		ticket.EventUUID = newEventUUID

		err = p.InsertTicket(ctx, tx, ticket)
		if err != nil {
			return "", err
		}
	}

	location.EventUUID = newEventUUID
	err = p.InsertEventLocation(ctx, tx, location)
	if err != nil {
		return "", err
	}

	err = p.InsertTags(ctx, tx, newEventUUID, tags)
	if err != nil {
		return "", err
	}

	err = p.UpdateEventImageURL(ctx, tx, newEventUUID, event.ImageURL)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("commit error: %v", err)
	}

	return newEventUUID, nil
}