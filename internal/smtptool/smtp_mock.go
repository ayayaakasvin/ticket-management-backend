package smtptool

import (
	"net/mail"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
)

type MockSMTP struct {
	GenerateRandomSequenceFunc func() (string, error)
	SendCodeFunc               func(subject string, code string, to []string) error
	ValidateEmailFunc          func(address string) bool
	HealthCheckFunc            func() error
}

var _ domain.SMTP = (*MockSMTP)(nil)

func (m *MockSMTP) GenerateRandomSequence() (string, error) {
	if m.GenerateRandomSequenceFunc == nil {
		panic("MockSMTP.GenerateRandomSequenceFunc is not configured")
	}
	return m.GenerateRandomSequenceFunc()
}

func (m *MockSMTP) SendCode(subject string, code string, to []string) error {
	if m.SendCodeFunc == nil {
		panic("MockSMTP.SendCodeFunc is not configured")
	}
	return m.SendCodeFunc(subject, code, to)
}

func (m *MockSMTP) ValidateEmail(address string) bool {
	if m.ValidateEmailFunc == nil {
		panic("MockSMTP.ValidateEmailFunc is not configured")
	}
	return m.ValidateEmailFunc(address)
}

func (m *MockSMTP) HealthCheck() error {
	if m.HealthCheckFunc == nil {
		panic("MockSMTP.HealthCheckFunc is not configured")
	}
	return m.HealthCheckFunc()
}

func NewSMTP_MockUp() (*MockSMTP, error) {
	smtp := &MockSMTP{
		GenerateRandomSequenceFunc: func() (string, error) {
			return "1234567", nil
		},
		SendCodeFunc: func(subject string, code string, to []string) error {
			return nil
		},
		ValidateEmailFunc: func(address string) bool {
			_, err := mail.ParseAddress(address)
			return err == nil
		},
		HealthCheckFunc: func() error {
			return nil
		},
	}

	return smtp, nil
}
