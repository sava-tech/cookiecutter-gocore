package smtp

import (
	"fmt"
	"net/smtp"

	"{{ cookiecutter.module_path }}/pkg/emailer/templates"
)

func (m *SMTPMailer) SendEmailOTP(identifier string, token string) (string, error) {
	data := EmailTemplateData{
		Identifier: identifier,
		Token:      token,
		Link:       fmt.Sprintf("https://{{ cookiecutter.project_name }}.co/verify?token=%s", token),
		Year:       2026,
		AppName:    "{{ cookiecutter.project_name }}",
	}

	htmlContent, err := templates.GenerateHTML("./pkg/emailer/templates/otp.html", data)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML: %w", err)
	}

	textContent := fmt.Sprintf("Email Verification Code\n\nYour verification code is: %s\n\nCopy this code and paste it in the app to verify your email.\n\nThank you,\n%s Team", token, data.AppName)

	msg := m.buildEmailMessage(identifier, "Email Verification Code", textContent, htmlContent)

	if err := m.send(identifier, msg); err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return "OTP sent via SMTP", nil
}

func (m *SMTPMailer) SendPasswordReset(identifier string) (string, error) {
	resetLink := fmt.Sprintf("https://{{ cookiecutter.project_name }}.co/reset-password?email=%s", identifier)
	data := EmailTemplateData{
		Identifier: identifier,
		Link:       resetLink,
		Year:       2026,
		AppName:    "{{ cookiecutter.project_name }}",
	}

	htmlContent, err := templates.GenerateHTML("./pkg/emailer/templates/password_with_link.html", data)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML: %w", err)
	}

	textContent := fmt.Sprintf("Password Reset Request\n\nClick this link to reset your password: %s\n\nIf you didn't request this, please ignore this email.\n\nThank you,\n%s Team", resetLink, data.AppName)

	msg := m.buildEmailMessage(identifier, "Password Reset Request", textContent, htmlContent)

	if err := m.send(identifier, msg); err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return "", nil
}

func (m *SMTPMailer) SendWelcomeMessage(identifier string) (string, error) {
	return "Welcome message sent via SMTP", nil
}

// send authenticates (when credentials are set) and relays the message
// through the configured SMTP host — this is what actually varies between
// Google/Mailgun/Zoho/SendGrid/Amazon SES/Postmark.
func (m *SMTPMailer) send(to, msg string) error {
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)

	var auth smtp.Auth
	if m.Username != "" {
		auth = smtp.PlainAuth("", m.Username, m.Password, m.Host)
	}

	return smtp.SendMail(addr, auth, m.From, []string{to}, []byte(msg))
}

func (m *SMTPMailer) buildEmailMessage(to, subject, textContent, htmlContent string) string {
	msg := fmt.Sprintf("From: %s\r\n", m.From)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: multipart/alternative; boundary=boundary\r\n"
	msg += "\r\n"
	msg += "--boundary\r\n"
	msg += "Content-Type: text/plain; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += textContent + "\r\n"
	msg += "\r\n"
	msg += "--boundary\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n"
	msg += htmlContent + "\r\n"
	msg += "\r\n"
	msg += "--boundary--\r\n"

	return msg
}
