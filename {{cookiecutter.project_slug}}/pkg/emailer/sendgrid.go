package emailer

type SendGridMailer struct {
	ApiKey string
}

func (s *SendGridMailer) SendOTP(identifier string, token string) (string, error) {
	return "SendGrid OTP sent", nil
}

func (s *SendGridMailer) SendWelcomeMessage(identifier string) (string, error) {
	return "SendGrid Welcome message sent", nil
}
