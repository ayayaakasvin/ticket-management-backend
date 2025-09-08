package models

import "time"

type Event struct {
	EventUUID      	string    	`json:"event_uuid"`

	CreationTime   	time.Time 	`json:"creation_time"`
	StartingTime   	time.Time 	`json:"starting_time"`
	EndingTime     	time.Time 	`json:"ending_time"`

	Title          	string    	`json:"title"`
	Description    	string    	`json:"description"`
	CategoryID	   	uint		`json:"category_id"`

	Status         	string    	`json:"status"`
	
	Capacity       	uint      	`json:"capacity"`
	
	Tags           	[]string  	`json:"tags"`
	ImageURL       	string    	`json:"image_url"`
	
	OrganizerID    	uint64    	`json:"organizer_id"`
	OrganizerName  	string    	`json:"organizer_name"`
	
	Tickets        	[]string  	`json:"ticketsId"`
	
	Location       	Location  	`json:"location"`
}