package emailer

import (
	"github.com/michaelassa01/blog/utils"
)

func SendEmailOTP(identifier string, token string, cfg utils.Config) (string, error) {

	mailer, err := NewMailer(cfg)
	if err != nil {
		return "", err
	}

	res, err := mailer.SendEmailOTP(identifier, token)
	return res, err
}

func SendWelcomeMessage(identifier string, cfg utils.Config) (string, error) {

	mailer, err := NewMailer(cfg)
	if err != nil {
		return "", err
	}
	res, err := mailer.SendWelcomeMessage(identifier)
	return res, err
}


