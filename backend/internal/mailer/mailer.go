package mailer

import "embed"

const (
	// FromName is the name used in the "From" field of the email.
	FromName = "Thyme Flies"
	// maxRetires is the maximum number of retries for sending an email.
	maxRetires = 3
	// UserWelcomeTemplate is the template file for the user welcome email.
	UserWelcomeTemplate = "user_invitation.tmpl"
	// PasswordResetTemplate is the template file for the password reset email.
	PasswordResetTemplate = "password_reset.tmpl"
)

//go:embed "templates"
var FS embed.FS

// Client is an interface for sending emails using different templates.
type Client interface {
	// Send sends an email using the specified template file, email address, data, and sandbox mode.
	Send(templateFile, email string, data any, isSandbox bool) (int, error)
}
