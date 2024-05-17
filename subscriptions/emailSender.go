package subscriptions

import (
	"gopkg.in/gomail.v2"
	"test_project/settings"
)

type EmailSender struct {
	dialer *gomail.Dialer
}

func NewEmailSender(settings settings.AppSettings) *EmailSender {
	auth := gomail.NewDialer(settings.Email.Host, settings.Email.Port, settings.Email.Email, settings.Email.Password)
	return &EmailSender{dialer: auth}
}

func (e *EmailSender) SendEmail(to string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", e.dialer.Username)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Daily rates email")
	message.SetBody("text/html", body)

	return e.dialer.DialAndSend(message)
}
