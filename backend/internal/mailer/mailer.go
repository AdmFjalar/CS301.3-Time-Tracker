package mailer

import "embed"

const (
	FromName            = "Thyme Flies"
	maxRetires          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, email string, data any, isSandbox bool) (int, error)
}
