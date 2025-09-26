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
