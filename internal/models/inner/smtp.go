package inner

type SMTP interface {
	GenerateRandomSequence() int
	SendCode(subject string, code int, to []string) error
	ValidateEmail(address string) bool
	HealthCheck() error
}