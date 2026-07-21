package mailtrap

import "github.com/michaelassa01/blog/utils"

type MailtrapMailer struct {
	ApiKey string
	Config utils.Config
}

// EmailTemplateData Data structure passed into HTML template
type EmailTemplateData struct {
	Identifier string
	Token      string
	Link       string
}

// EmailAddress Email Sender
type MailtrapEmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type MailtrapEmailPayload struct {
	From     MailtrapEmailAddress   `json:"from"`
	To       []MailtrapEmailAddress `json:"to"`
	Subject  string                 `json:"subject"`
	HTML     string                 `json:"html,omitempty"`
	Category string                 `json:"category,omitempty"`
}
