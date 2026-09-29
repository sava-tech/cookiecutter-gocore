package emailer

import (
	"github.com/michaelassa01/blog/pkg/emailer/mailpit"
	"github.com/michaelassa01/blog/pkg/emailer/mailtrap"
	"github.com/michaelassa01/blog/utils"
)

type Config struct {
	Provider          string
	MailtrapAuthToken string
	SendGridAuthToken string
}

func NewMailer(cfg utils.Config) (Mailer, error) {
	switch cfg.Provider {
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
