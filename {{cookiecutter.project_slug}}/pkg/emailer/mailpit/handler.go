package mailpit

import (
	"fmt"
	"net/smtp"
	"os"
)

func (m *MailpitMailer) SendEmailOTP(identifier string, token string) (string, error) {
	// Implementation for sending OTP via Mailpit
	// Check if template exists
	templatePath := "./pkg/emailer/templates/otp.html"
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", templatePath)
	}
	data := EmailTemplateData{
		Identifier: identifier,
		Token:      token,
		Link:       fmt.Sprintf("https://{{ cookiecutter.project_name }}.co/verify?token=%s", token),
		Year:       2026,
		AppName:    "CNN Nigeria",
	}
	// Generate HTML content
	htmlContent, err := GenerateHTML(templatePath, data)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Generate plain text fallback
	textContent := fmt.Sprintf("Email Verification Code\n\nYour verification code is: %s\n\nCopy this code and paste it in the app to verify your email.\n\nThank you,\n%s Team", token, data.AppName)

	// Build email message with both plain text and HTML
	msg := m.buildEmailMessage(identifier, "Email Verification Code", textContent, htmlContent)

	addr := fmt.Sprintf("%s:%s", m.SmtpHost, m.SmtpPort)
	// Send OTP mailer using SMTP
	if err := smtp.SendMail(addr, nil, m.From, []string{identifier}, []byte(msg)); err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return "OTP sent via Mailpit", nil
}

func (m *MailpitMailer) SendPasswordReset(identifier string) (string, error) {
	// Implementation for sending password reset via Mailpit
	// Prepare template data
	resetLink := fmt.Sprintf("https://{{ cookiecutter.project_name }}.co/reset-password?email=%s", identifier)
	data := EmailTemplateData{
		Identifier: identifier,
		Link:       resetLink,
		Year:       2026,
		AppName:    "CNN Nigeria",
	}

	// Generate HTML content
	htmlContent, err := GenerateHTML("./pkg/emailer/templates/password_with_link.html", data)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML: %w", err)
	}
	fmt.Println("3. HTML generated, length: ", len(htmlContent), " bytes")

	// Generate plain text fallback
	textContent := fmt.Sprintf("Password Reset Request\n\nClick this link to reset your password: %s\n\nIf you didn't request this, please ignore this email.\n\nThank you,\n%s Team", resetLink, data.AppName)

	// Build email message
	msg := m.buildEmailMessage(identifier, "Password Reset Request", textContent, htmlContent)

	// Send email via Mailpit
	addr := fmt.Sprintf("%s:%s", m.SmtpHost, m.SmtpPort)
	if err := smtp.SendMail(addr, nil, m.From, []string{identifier}, []byte(msg)); err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return "", nil
}

func (m *MailpitMailer) SendWelcomeMessage(identifier string) (string, error) {
	// Implementation for sending welcome message via Mailpit
	return "Welcome message sent via Mailpit", nil
}

func (m *MailpitMailer) buildEmailMessage(to, subject, textContent, htmlContent string) string {
	// Create MIME multipart message
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
