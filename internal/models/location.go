package models

type Location struct {
    LocationID  uint    `json:"location_id,omitempty"`
    EventUUID   string  `json:"event_uuid.omitempty"`
    Name        string  `json:"name,omitempty"`      // "Dostyk Plaza"
    Address     string  `json:"address,omitempty"`   // "Dostyk Ave 85, Almaty, Kazakhstan"
    Latitude    float64 `json:"latitude"`  // 43.2405
    Longitude   float64 `json:"longitude"` // 76.9312
}