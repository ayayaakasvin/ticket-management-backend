package models

import "time"

type Ticket struct {
	TicketUUID   string    `json:"ticket_uuid"`
	EventUUID    string    `json:"event_uuid"`

	Name         string    `json:"name"`
	
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	Quantity     uint      `json:"quantity"`
	Sold         uint      `json:"sold"`
	
	StartingTime time.Time `json:"starting_time"`
	EndingTime   time.Time `json:"ending_time"`
}
