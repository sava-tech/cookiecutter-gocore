package emailer

import (
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/emailer/mailpit"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/emailer/mailtrap"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/utils"
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
