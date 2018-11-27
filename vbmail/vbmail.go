package vbmail

import (
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	client *sendgrid.Client
)

func Init(sendgridSecret string) {
	client = sendgrid.NewSendClient(sendgridSecret)
}

// SendTo sends a new email using sendgrid service
func SendTo(subject string, receiverName string, receiverEmail string, plainText string, htmlText string) error {
	from := mail.NewEmail("Vikebot", "noreply@vikebot.com")
	to := mail.NewEmail(receiverName, receiverEmail)

	message := mail.NewSingleEmail(from, subject, to, plainText, htmlText)

	_, err := client.Send(message)
	return err
}
