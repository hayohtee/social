package mailer

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"time"
)

type SendGridMailer struct {
	senderEmail string
	apiKey      string
	client      *sendgrid.Client
}

func NewSendGridMailer(senderEmail, apiKey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		senderEmail: senderEmail,
		apiKey:      apiKey,
		client:      client,
	}
}

func (s SendGridMailer) Send(templateFile, username, email string, data any) error {
	from := mail.NewEmail(senderName, s.senderEmail)
	to := mail.NewEmail(username, email)

	//TODO: Template parsing

	subject := new(bytes.Buffer)
	body := new(bytes.Buffer)
	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	for i := 0; i < maxRetries; i++ {
		response, err := s.client.Send(message)
		if err != nil {
			fmt.Printf("failed to send email to %v, attempt %d of %d\n", email, i+1, maxRetries)
			fmt.Println("Error:", err)

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Println("Email sent with status code", response.StatusCode)
		return nil
	}

	return fmt.Errorf("failed to send email after %d the attempts", maxRetries)
}
