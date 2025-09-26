package inner

import (
	"context"
	"database/sql"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, username, hashedpassword string) error
	AuthentificateUser(ctx context.Context, username, password string) (uint, error)
}

type EventRepository interface {
	InsertEvent(ctx context.Context, tx *sql.Tx, eventObj *models.Event) (string, error)
	GetEventByUUID(ctx context.Context, eventUUID string) (*models.Event, error)
	GetAllEvents(ctx context.Context) ([]*models.Event, error)
	UpdateEventImageURL(ctx context.Context, eventUUID string, imageURL string) error
	DeleteEventByUUID(ctx context.Context, eventUUID string) error
	GetEventsByCategoryID(ctx context.Context, categoryID uint) ([]*models.Event, error)

	InsertEventObjectToDatabase(ctx context.Context, event *models.Event) (string, error)

	CategoryRepository
	TicketRepository
}

type CategoryRepository interface{
	GetAllCategories(ctx context.Context) ([]models.Category, error)
}

type TicketRepository interface {
	InsertTicketAfterwards(ctx context.Context, ticketObj *models.Ticket) error
	DeleteTicket(ctx context.Context, ticketUUID string) error
	GetTicket(ctx context.Context, ticketUUID string) (*models.Ticket, error)
}