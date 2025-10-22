# HTTP Routes

This document lists all HTTP routes exposed by the server, the HTTP method, the handler function, whether the route requires authentication, and expected request/response shapes.

Note: handler request/response shapes were inferred from handler code and the `models` package. Where exact model fields are required, check the corresponding model definitions in `internal/models`.

## Global

- GET /ping
  - Handler: inline in `internal/http-server/http-server.go`
  - Auth: none
  - Response: 200 "pong"

## API group (/api)

### Authentication

- POST /api/login
  - Handler: `handlers.LogIn()`
  - Auth: none
  - Request JSON: {"username": "string", "password": "string"}
  - Response JSON (200): {"data": {"access-token": "string", "refresh-token": "string"}}

- POST /api/register
  - Handler: `handlers.Register()`
  - Auth: none
  - Request JSON: {"username": "string", "password": "string"}
  - Response: 201 Created, no body (or empty data)

- DELETE /api/logout
  - Handler: `handlers.LogOut()`
  - Auth: JWT required
  - Request: Authorization header with Bearer access-token
  - Response: 200 OK

- POST /api/refresh
  - Handler: `handlers.RefreshTheToken()`
  - Auth: none (but requires refresh token in Authorization header)
  - Request: Authorization: Bearer <refresh-token>
  - Response JSON (200): {"data": {"access-token": "string"}}

### Event routes (/api/event)

All `/api/event` routes (except `/top-ten` and `/category`) require JWT authentication.

- GET /api/event/all
  - Handler: `handlers.GetAllEvents()`
  - Auth: JWT
  - Response JSON: {"data": {"events": [Event]}}

- GET /api/event/category
  - Handler: `handlers.GetEventsByCategoryID()`
  - Auth: JWT
  - Query params: `category_id` (uint)
  - Response JSON: {"data": {"events": [Event]}}

- GET /api/event/top-ten
  - Handler: `handlers.GetTop10Events()`
  - Auth: JWT
  - Response JSON: {"data": {"trending": [Event]}}

- POST /api/event
  - Handler: `handlers.SaveEvent()`
  - Auth: JWT
  - Request JSON: models.Event (see `internal/models/event.go`)
  - Response JSON (201): {"data": {"event_uuid": "uuid"}}

- POST /api/event/update/upload?event_uuid=<uuid>
  - Handler: `handlers.UpdateEventImageURLByUploading()`
  - Auth: JWT
  - Request: multipart/form-data with file field `image` (max 10MB). Valid mime types: image/png, image/jpeg, image/webp
  - Response: 200 OK

- POST /api/event/update/image?event_uuid=<uuid>&image_url=<url>
  - Handler: `handlers.UpdateEventImageURLUsingExternalSource()`
  - Auth: JWT
  - Request: query param `image_url` (string)
  - Response: 200 OK

- GET /api/event?event_uuid=<uuid>
  - Handler: `handlers.GetEventByUUID()`
  - Auth: JWT
  - Query params: `event_uuid` (string UUID)
  - Response JSON: {"data": {"event": Event}}

- DELETE /api/event?event_uuid=<uuid>
  - Handler: `handlers.DeleteEventByUUID()`
  - Auth: JWT
  - Query params: `event_uuid` (string UUID)
  - Response: 200 OK

### Category

- GET /api/category
  - Handler: `handlers.GetAllCategories()`
  - Auth: JWT
  - Response JSON: {"data": {"categories": [Category]}}

### Images

- GET /api/images/<path>
  - Handler: `handlers.ServeImages()`
  - Auth: JWT
  - Response: serves file from custom FS; path is trimmed from `/images/`

### Ticket routes (/api/ticket)

- GET /api/ticket?ticket_uuid=<uuid>
  - Handler: `handlers.GetTicket()`
  - Auth: none (no middleware attached in route registration)
  - Query params: `ticket_uuid` (string UUID)
  - Response JSON: {"data": {"ticket": Ticket}}

- POST /api/ticket
  - Handler: `handlers.InsertTicketAfterwards()`
  - Auth: none
  - Request JSON: models.Ticket (must include EventUUID)
  - Response JSON (200): {"data": {"ticket_uuid": "uuid"}}

- DELETE /api/ticket?ticket_uuid=<uuid>
  - Handler: `handlers.DeleteTicket()`
  - Auth: none
  - Query params: `ticket_uuid` (string UUID)
  - Response: 200 OK


## Models referenced (high level)
- Event: see `internal/models/event.go`
- Ticket: see `internal/models/event.go` (ticket struct)
- Category: see `internal/repo/postgresql` migration (category table)

## Exact JSON shapes (from models)

Event (JSON)

{
  "event_uuid": "string (uuid, optional)",
  "creation_time": "string (RFC3339 timestamp)",
  "starting_time": "string (RFC3339 timestamp)",
  "ending_time": "string (RFC3339 timestamp)",
  "title": "string",
  "description": "string (optional)",
  "category_id": "uint (optional)",
  "status": "string",
  "capacity": "uint (optional)",
  "tags": ["string", ...],
  "image_url": "string (optional)",
  "organizer_id": "uint",
  "organizer_name": "string (optional)",
  "ticketsId": [ { Ticket }, ... ],
  "location": {
    "name": "string",
    "address": "string",
    "latitude": "number (float, optional)",
    "longitude": "number (float, optional)"
  }
}

Ticket (JSON)

{
  "ticket_uuid": "string (uuid, optional)",
  "event_uuid": "string (uuid)",
  "name": "string",
  "price": "number (float, optional)",
  "currency": "string (optional)",
  "quantity": "uint",
  "sold": "uint (optional)"
}

## Notes
- Routes that used `mws.JWTAuthMiddleware` during registration require a valid access token in the `Authorization` header: `Authorization: Bearer <token>` unless the handler explicitly checks otherwise.
- Many handlers expect query parameters (e.g., event_uuid, ticket_uuid, category_id); ensure clients provide them.
- Request validation happens in handlers (e.g., `validinput.ValidateEventSave`), so the handlers return 400/401/404 appropriately.

If you'd like, I can:
- Expand each `Event`/`Ticket` model field in this document by reading `internal/models/*.go` and adding the exact JSON shapes.
- Generate an OpenAPI spec (yaml) from these routes.
- Add example curl requests for each route.
