package emailer

import (
	"{{ cookiecutter.module_path }}/pkg/emailer"
	c "{{ cookiecutter.module_path }}/utils"
)

func SendEmailOTP(identifier string, token string, cfg c.Config) (string, error) {

	mailer, err := emailer.NewMailer(cfg)
	if err != nil {
		return "", err
	}

	res, err := mailer.SendEmailOTP(identifier, token)
	return res, err
}

func SendSMSMessage(identifier string, token string, cfg c.Config) (string, error) {
	mailer, err := emailer.NewSMSMessageSender(cfg)
	if err != nil {
		return "", err
	}

	res, err := mailer.SendPhoneOTP(identifier, token)
	return res, err
}
