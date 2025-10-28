package models

import "time"

type Event struct {
	EventUUID string `json:"event_uuid,omitempty"`

	CreationTime time.Time `json:"creation_time"`
	StartingTime time.Time `json:"starting_time"`
	EndingTime   time.Time `json:"ending_time"`

	Title string `json:"title"`

	Description string `json:"description,omitempty"`
	CategoryID  uint   `json:"category_id,omitempty"`

	Status string `json:"status"`

	Capacity uint `json:"capacity,omitempty"`

	Tags     []string `json:"tags"`
	ImageURL string   `json:"image_url,omitempty"`

	OrganizerID   uint `json:"organizer_id"`
	OrganizerName string `json:"organizer_name,omitempty"`

	Tickets []*Ticket `json:"ticketsId"`

	Location Location `json:"location"`
}

type EventStats struct {
	EventUUID         string    `json:"event_uuid"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	StartingTime      time.Time `json:"starting_time"`
	EndingTime        time.Time `json:"ending_time"`
	Status            string    `json:"status"`
	Capacity          int       `json:"capacity"`
	ImageURL          string    `json:"image_url"`
	CategoryName      string    `json:"category_name"`
	OrganizerUsername string    `json:"organizer_username"`
	TotalTicketsSold  int       `json:"total_tickets_sold"`
	Rank              int       `json:"rank"`

	// Derived field — computed in Go, not stored in DB
	FillRate float64 `json:"fill_rate"`
}
