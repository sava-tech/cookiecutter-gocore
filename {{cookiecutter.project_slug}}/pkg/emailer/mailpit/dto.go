package mailpit

import "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/utils"


type MailpitMailer struct {
	Config utils.Config
	From     string
	SmtpHost string
	SmtpPort string
}

// EmailTemplateData Data structure passed into HTML template
type EmailTemplateData struct {
	Identifier string
	Token      string
	Link       string
	Year       int
	AppName    string
}

func NewMailpitMailer(from, smtpHost, smtpPort string, config utils.Config) *MailpitMailer {
	return &MailpitMailer{
		From:     from,
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
		Config:   config,
	}
}


