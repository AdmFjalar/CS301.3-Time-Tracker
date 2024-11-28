package mailer

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridMailer is a struct that holds the SendGrid client, API key, and sender email address.
type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

// NewSendgrid creates a new SendGridMailer with the given API key and sender email address.
func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

// Send sends an email using the specified template file, recipient email address, data, and sandbox mode.
// It retries sending the email up to maxRetires times in case of failure.
func (m *SendGridMailer) Send(templateFile, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail("User", email)

	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	enableSandbox := false
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			// Enable: &isSandbox,
			Enable: &enableSandbox,
		},
	})

	var retryErr error
	for i := 0; i < maxRetires; i++ {
		response, retryErr := m.client.Send(message)
		if retryErr != nil {
			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return response.StatusCode, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempt, error: %v", maxRetires, retryErr)
}
