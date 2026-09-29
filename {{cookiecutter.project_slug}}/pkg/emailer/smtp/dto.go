package smtp

import "{{ cookiecutter.module_path }}/utils"

// SMTPMailer sends mail over an authenticated SMTP relay. It's provider
// agnostic — Google, Mailgun, Zoho, SendGrid, Amazon SES and Postmark all
// expose a standard SMTP endpoint, so one client covers all of them; only
// Host/Port/Username/Password differ per provider (see env.example).
type SMTPMailer struct {
	Config   utils.Config
	From     string
	Host     string
	Port     string
	Username string
	Password string
}

// EmailTemplateData Data structure passed into HTML template
type EmailTemplateData struct {
	Identifier string
	Token      string
	Link       string
	Year       int
	AppName    string
}
