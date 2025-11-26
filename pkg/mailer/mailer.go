package mailer

import "embed"

const (
	FromName             = "GopherSocial"
	maxRetries           = 3
	UserWelcomeTemplate  = "user_invitation.tmpl"
	UserResetPwsTemplate = "user_reset.tmpl"
	UserStaffActivate    = "user_staff_activate.tmpl"
	InvoiceTemplate      = "invoice.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile string, username, email string, data any, isSandbox bool) (int, error)
}
