package mailtrap

import "fmt"



func (m *MailtrapMailer) SendEmailOTP(identifier string, token string) (string, error) {
	// build payload and send using Mailtrap
	subject := "[Cnn-Nigeria]-OTP Token"
	data := EmailTemplateData{
		Identifier: identifier,
		Token:      token,
		Link:       fmt.Sprintf("https://cnn-Nigeria.co/verify?token=%s", token),
	}

	// Generate HTML content
	html, err := GenerateHTML("./pkg/emailer/templates/otp.html", data)
	if err != nil {
		return "", err
	}

	// Send OTP mailer
	err = Sender(identifier, subject, html, m.Config)
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("OTP token sent to the identifier %s", identifier)
	return res, nil
}

func (m *MailtrapMailer) SendWelcomeMessage(identifier string) (string, error) {
	return "Mailtrap Welcome message sent", nil
}
