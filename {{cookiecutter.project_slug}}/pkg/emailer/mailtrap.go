package emailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	c "{{ cookiecutter.module_path }}/utils"
)

type MailtrapMailer struct {
	ApiKey string
	Config c.Config
}

// EmailTemplateData Data structure passed into HTML template
type EmailTemplateData struct {
	Identifier string
	Token      string
	Link       string
}

func (m *MailtrapMailer) SendEmailOTP(identifier string, token string) (string, error) {
	// build payload and send using Mailtrap
	subject := "[{{ cookiecutter.project_slug}}]-OTP Token"
	data := EmailTemplateData{
		Identifier: identifier,
		Token:      token,
		Link:       fmt.Sprintf("https://{{ cookiecutter.project_slug}}.co/verify?token=%s", token),
	}

	// Generate HTML content
	html, err := GenerateHTML("./internal/emailer/templates/otp.html", data)
	if err != nil {
		panic(err)
	}

	// Send OTP mailer
	err = Sender(identifier, subject, html, m.Config)
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf(" OTP token sent to the identifier %s", identifier)
	return res, nil
}

// SendWelcomeMessage Send Welcome Message
func (m *MailtrapMailer) SendWelcomeMessage(identifier string) (string, error) {
	return "Mailtrap Welcome message sent", nil
}

// EmailAddress Email Sender
type EmailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type EmailPayload struct {
	From     EmailAddress   `json:"from"`
	To       []EmailAddress `json:"to"`
	Subject  string         `json:"subject"`
	HTML     string         `json:"html,omitempty"`
	Category string         `json:"category,omitempty"`
}

// Sender handles sending email using mailtrap
func Sender(identifier string, subject string, html string, config c.Config) error {

	url := config.MailtrapURL
	method := "POST"

	// Prepare payload
	payload := EmailPayload{
		From: EmailAddress{
			Email: config.DefaultFromEmail,
			Name:  config.EmailSubjectPrefit,
		},
		To: []EmailAddress{
			{Email: identifier},
		},
		Subject: subject,
		HTML:    html,
	}
	// Marshal payload to json
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		fmt.Println(err)
		return err
	}
	// API Bearer token
	BearerToken := fmt.Sprintf("Bearer %s", config.MailtrapAuthToken)

	req.Header.Add("Authorization", BearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(body))
	return err
}
