package inner

import (
	"context"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
)

type UserRepository interface {
	RegisterUser		(ctx context.Context, username, hashedpassword string) 		error
	AuthentificateUser	(ctx context.Context, username, password string) 			(int, error)
}

type EventRepository interface {
	InsertEvent 		(ctx context.Context, eventObj *models.Event)				error
	GetEvent 			(ctx context.Context, eventUUID string)						(*models.Event, error)
	GetEvents	 		(ctx context.Context, eventObj *models.Event)				error
}