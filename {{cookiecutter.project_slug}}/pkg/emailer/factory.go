package emailer

import (
	"fmt"

	c "{{ cookiecutter.module_path }}/utils"
)

type Config struct {
	Provider          string
	MailtrapAuthToken string
	SendGridAuthToken string
	SMSProvider       string
	TermiiApiKey      string
}

func NewMailer(cfg c.Config) (Mailer, error) {
	switch cfg.Provider {
	case "mailtrap":
		return &MailtrapMailer{ApiKey: cfg.MailtrapAuthToken, Config: cfg}, nil

	//case "sendgrid":
	//	return &SendGridMailer{ApiKey: cfg.SendGridAuthToken}, nil

	default:
		return nil, fmt.Errorf("unsupported mailer provider: %s", cfg.Provider)
	}
}

func NewSMSMessageSender(cfg c.Config) (SMSender, error) {
	switch cfg.SMSProvider {
	case "termii":
		return &TermiiSender{APIKey: cfg.TermiiApiKey, Config: cfg}, nil

	default:
		return nil, fmt.Errorf("unsupported mailer provider: %s", cfg.Provider)

	}
}
