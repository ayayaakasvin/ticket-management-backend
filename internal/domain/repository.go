package domain

import (
	"context"
	"database/sql"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) (int64, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Update(ctx context.Context, user *User) error
	IsAdmin(ctx context.Context, id int64) bool
}

type EventRepository interface {
	InsertEventTx(ctx context.Context, tx *sql.Tx, event *Event) (string, error)
	GetEvent(ctx context.Context, eventUUID string) (*Event, error)
	GetEvents(ctx context.Context) ([]*Event, error)
	UpdateEventImage(ctx context.Context, eventUUID string, imageURL string) error
	DeleteEvent(ctx context.Context, eventUUID string) error
	GetEventsByCategory(ctx context.Context, categoryID uint) ([]*Event, error)

	CreateEvent(ctx context.Context, event *Event) (string, error)

	CategoryRepository
	TicketRepository
}

type CategoryRepository interface {
	GetCategories(ctx context.Context) ([]Category, error)
}

type TicketRepository interface {
	CreateTicket(ctx context.Context, ticket *Ticket) error
	DeleteTicket(ctx context.Context, ticketUUID string) error
	GetTicket(ctx context.Context, ticketUUID string) (*Ticket, error)
}
