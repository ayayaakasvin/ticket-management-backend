package smtptool

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/mail"
	"net/smtp"

	"github.com/ayayaakasvin/oneflick-ticket/mailtemplates"
)

const (
	max = 1_000_000
)

type SMTP_Client struct {
	auth smtp.Auth

	username string
	password string
	host     string
	port     int
}

func NewSMTPTool(
	username string,
	password string,
	host string,
	port int,
) *SMTP_Client {
	auth := smtp.PlainAuth(
		"",
		username,
		password,
		host,
	)

	return &SMTP_Client{
		auth:     auth,
		username: username,
		password: password,
		host:     host,
		port:     port,
	}
}

func NewSMTPToolWithPreHealthCheck(
	username string,
	password string,
	host string,
	port int,
) (*SMTP_Client, error) {
	s := NewSMTPTool(username, password, host, port)

	err := s.HealthCheck()
	if err != nil {
		return nil, fmt.Errorf("failed to healthcheck to SMTP: %v\n", err)
	}

	return s, nil
}

func (s *SMTP_Client) GenerateRandomSequence() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(9_000_000))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%07d", n.Int64()+max), nil
}

func (s *SMTP_Client) SendCode(subject string, code string, to []string) error {
	msg := fmt.Sprintf(mailtemplates.MailTemplate, subject, code, code)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.host, s.port),
		s.auth,
		s.username,
		to,
		[]byte(msg),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *SMTP_Client) ValidateEmail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func (s *SMTP_Client) HealthCheck() error {
	from := s.username
	to := []string{from}

	msg := []byte("Subject: SMTP Health Check\r\n\r\nThis is a test.")

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.host, s.port),
		s.auth,
		from,
		to,
		msg,
	)

	return err
}
