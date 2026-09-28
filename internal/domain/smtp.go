package domain

type SMTP interface {
	GenerateRandomSequence() (string, error)
	SendCode(subject string, code string, to []string) error
	ValidateEmail(address string) bool
	HealthCheck() error
}
