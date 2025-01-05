package models

import (
	"fmt"
	"html/template"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultSender = "support@lenslocked.com"
)

type Email struct {
	To        string
	From      string
	Subject   string
	Plaintext string
	HTML      template.Template
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailService struct {
	//DefaultSender is used as default when no sender is provided
	DefaultSender string
	dialer        *mail.Dialer
}

func NewEmailService(config SMTPConfig) (*EmailService, error) {
	es := EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
	return &es, nil
}

func (es *EmailService) ForgotPassword(recipient, resetUrl string, tpl template.Template) error {
	email := Email{
		To:        recipient,
		Subject:   "Reset Password",
		Plaintext: "To reset your password please visit the following link " + resetUrl,
		HTML:      tpl,
	}

	err := es.Send(email)
	if err != nil {
		return fmt.Errorf("forgot password email: %w", err)
	}
	return nil
}

func (es *EmailService) Send(email Email) error {

	msg := mail.NewMessage()
	msg.SetHeader("To", email.To)
	es.setFrom(msg, email)
	msg.SetHeader("Subject", email.Subject)

	switch {
	case email.Plaintext != "" && email.HTML != "":
		msg.SetBody("text/plain", email.Plaintext)
		msg.AddAlternative("text/html", email.HTML)
	case email.Plaintext != "":
		msg.SetBody("text/plain", email.Plaintext)
	case email.HTML != "":
		msg.SetBody("text/html", email.HTML)
	}

	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	return nil
}

func (es *EmailService) setFrom(msg *mail.Message, email Email) {
	var from string
	switch {
	case email.From != "":
		from = email.From
	case es.DefaultSender != "":
		from = es.DefaultSender
	default:
		from = DefaultSender
	}
	msg.SetHeader("From", from)
}
