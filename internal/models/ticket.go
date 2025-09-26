package models

type Ticket struct {
	TicketUUID 	string 	`json:"ticket_uuid,omitempty"`
	EventUUID  	string 	`json:"event_uuid"`
	
	Name 		string 	`json:"name"`

	Price    	float64 `json:"price,omitempty"`
	Currency 	string  `json:"currency,omitempty"`
	Quantity 	uint    `json:"quantity"`
	Sold     	uint    `json:"sold,omitempty"`
}
