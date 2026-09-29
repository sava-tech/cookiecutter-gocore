package emailer

import (
	"{{ cookiecutter.module_path }}/pkg/emailer/mailpit"
	"{{ cookiecutter.module_path }}/pkg/emailer/mailtrap"
	"{{ cookiecutter.module_path }}/pkg/emailer/smtp"
	"{{ cookiecutter.module_path }}/utils"
)

// NewMailer picks a Mailer based on EMAIL_PROVIDER:
//   - "smtp"     -> generic SMTP relay (Google, Mailgun, Zoho, SendGrid, Amazon SES, Postmark — see env.example)
//   - "mailtrap" -> Mailtrap's HTTP send API
//   - anything else -> Mailpit, the local dev mail catcher
func NewMailer(cfg utils.Config) (Mailer, error) {
	switch cfg.Provider {
	case "smtp":
		return &smtp.SMTPMailer{
			Config:   cfg,
			From:     cfg.DefaultFromEmail,
			Host:     cfg.SmtpHost,
			Port:     cfg.SmtpPort,
			Username: cfg.SmtpUsername,
			Password: cfg.SmtpPassword,
		}, nil

	case "mailtrap":
		return &mailtrap.MailtrapMailer{ApiKey: cfg.MailtrapAuthToken, Config: cfg}, nil

	default:
		return &mailpit.MailpitMailer{
			Config:   cfg,
			From:     cfg.EmailFrom,
			SmtpHost: cfg.MailpitHost,
			SmtpPort: cfg.MailpitPort,
		}, nil
	}
}
