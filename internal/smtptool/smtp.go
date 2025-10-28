package smtptool

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"

	"github.com/ayayaakasvin/oneflick-ticket/internal/config"
	"github.com/ayayaakasvin/oneflick-ticket/internal/models/inner"
	"github.com/ayayaakasvin/oneflick-ticket/mailtemplates"
)

const (
	min = 100000
	max = 999999
	origin = "SMTP"
)

type SMTPTool struct {
	auth smtp.Auth

	cfg *config.SMTPConfig
}

func NewSMTPTool(cfg *config.SMTPConfig) *SMTPTool {
	auth := smtp.PlainAuth(
		"",
		cfg.Username,
		cfg.Password,
		cfg.Host,
	)

	return &SMTPTool{
		auth: auth,
		cfg:  cfg,
	}
}


func NewSMTPToolWithPreHealthCheck(cfg *config.SMTPConfig, shutdownChannel inner.ShutdownChannel) *SMTPTool {
	s := NewSMTPTool(cfg)

	err := s.HealthCheck()
	if err != nil {
		msg := fmt.Sprintf("failed to healthcheck to SMTP: %v\n", err)
		shutdownChannel.Send(inner.ShutdownMessage, origin, msg)
		return nil
	}

	return s
}

func (s *SMTPTool) GenerateRandomSequence() int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max+1-min))
	if err != nil {
		return int(time.Now().UnixNano()%(max+1-min)) + min
	}

	return min + int(nBig.Int64())
}

func (s *SMTPTool) SendCode(subject string, code int, to []string) error {
	codeStr := strconv.Itoa(code)

	msg := fmt.Sprintf(mailtemplates.MailTemplate, subject, codeStr, codeStr)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		s.auth,
		s.cfg.Username,
		to,
		[]byte(msg),
	)
	if err != nil {
		return err
	}	

	return nil
}

func (s *SMTPTool) ValidateEmail(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func (s *SMTPTool) HealthCheck() error {
	from := s.cfg.Username
	to := []string{from}

	msg := []byte("Subject: SMTP Health Check\r\n\r\nThis is a test.")

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
		s.auth,
		from,
		to,
		msg,
	)

	return err
}
