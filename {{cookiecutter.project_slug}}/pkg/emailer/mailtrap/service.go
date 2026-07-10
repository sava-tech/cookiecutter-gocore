package mailtrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/utils"
)

// Sender handles sending email using mailtrap
// It takes in the identifier (email address), subject, html content, and config
// It returns an error if the email sending fails
// It constructs the payload, marshals it to JSON, and makes an HTTP POST request to the Mailtrap API
func Sender(identifier string, subject string, html string, config utils.Config) error {

	url := config.MailtrapURL
	method := "POST"

	// Prepare payload
	payload := MailtrapEmailPayload{
		From: MailtrapEmailAddress{
			Email: config.DefaultFromEmail,
			Name:  config.EmailSubjectPrefit,
		},
		To: []MailtrapEmailAddress{
			{Email: identifier},
		},
		Subject: subject,
		HTML:    html,
	}
	// Marshal payload to json
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("email json marshal failed", err)
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
