package validinput

import (
	"regexp"
	"strings"
	"time"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
)

type ValidationError string

func (v ValidationError) Error() string {
	return string(v)
}

const (
	ErrorEmptyTitle               ValidationError = "Event title cannot be empty"
	ErrorEmptyTickets             ValidationError = "Event must have at least one ticket"
	ErrorEmptyLocation            ValidationError = "Event location cannot be empty"
	ErrorInvalidCapacity          ValidationError = "Event capacity is invalid"
	ErrorInvalidLocationPlacement ValidationError = "Event location coordinates are invalid"
	ErrorInvalidTime              ValidationError = "Event time is invalid"
	ErrorInvalidTicketPrice       ValidationError = "Ticket price is invalid"
	ErrorInvalidCurrencyFormat    ValidationError = "Ticket currency format is invalid"
)

var (
	uppercase = regexp.MustCompile(`[A-Z]`)
	lowercase = regexp.MustCompile(`[a-z]`)
	digit     = regexp.MustCompile(`[0-9]`)

	minLengthPassword = 8
	minLengthUsername = 3
)

func IsValidPassword(password string) bool {
	if len(password) < minLengthPassword {
		return false
	}

	if !uppercase.MatchString(password) {
		return false
	}

	if !lowercase.MatchString(password) {
		return false
	}

	if !digit.MatchString(password) {
		return false
	}

	return true
}

func IsValidUsername(username string) bool {
	if len(username) < minLengthUsername {
		return false
	}

	if !(lowercase.MatchString(username) || uppercase.MatchString(username)) {
		return false
	}

	return true
}

func IsValidFileName(filename string) bool {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return false
	}
	// Disallow common illegal filename characters
	illegal := regexp.MustCompile(`[\\/:\*\?"<>\|]`)
	return !illegal.MatchString(filename)
}

func ValidTimeOfEvent(starting, ending time.Time) bool {
	return (starting.After(time.Now()) && starting.Before(ending))
}

func ValidateEventSave(event *models.Event) error {
	if event.Title == "" {
		return ErrorEmptyTitle
	}

	if event.Tickets == nil {
		return ErrorEmptyTickets
	} else if err := ValidateTicket(event.Tickets); err != nil {
		return err
	}

	if ValidTimeOfEvent(event.StartingTime, event.EndingTime) {
		return ErrorInvalidTime
	}

	// (0.0, 0.0) is used as a sentinel value for "unset" or "invalid" coordinates.
	// It's a real location (Gulf of Guinea), but extremely unlikely for a real event.
	if event.Location.Longitude == 0.0 && event.Location.Latitude == 0.0 {
		return ErrorInvalidLocationPlacement
	}

	return nil
}

func ValidateTicket(tickets []*models.Ticket) error {
	currencyISO := regexp.MustCompile(`^[A-Z]{3}$`)
	for _, ticket := range tickets {
		if ticket.Price < 0 {
			return ErrorInvalidTicketPrice
		}
		if !currencyISO.MatchString(ticket.Currency) {
			return ErrorInvalidCurrencyFormat
		}
	}
	return nil
}
