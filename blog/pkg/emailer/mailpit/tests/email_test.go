package tests

import (
	"fmt"
	"net/smtp"
	"testing"
)

func TestMailpitConnection(t *testing.T) {
	addr := "127.0.0.1:1025"
	from := "test@blog.local"
	to := "user@example.com"

	msg := []byte("To: user@example.com\r\n" +
		"From: test@blog.local\r\n" +
		"Subject: Test Email\r\n" +
		"\r\n" +
		"Hello from CNN Nigeria!\r\n")

	err := smtp.SendMail(addr, nil, from, []string{to}, msg)
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	t.Log("Email sent successfully to Mailpit")
}

func TestMailpitIPv4(t *testing.T) {
	// Test IPv4 address
	addr := "127.0.0.1:1025"
	from := "test@blog.local"
	to := "test@example.com"

	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: Test IPv4\r\n\r\nHello!\r\n", to, from)

	err := smtp.SendMail(addr, nil, from, []string{to}, []byte(msg))
	if err != nil {
		t.Fatalf("Failed to send via IPv4: %v", err)
	}
	t.Log("✓ Email sent successfully via IPv4!")
}
