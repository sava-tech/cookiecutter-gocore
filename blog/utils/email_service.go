package utils

type EmailService struct {
	config Config
}

func NewEmailService(config Config) *EmailService {
	return &EmailService{
		config: config,
	}
}

func (s *EmailService) SendOTP(email string, code string) error {

	// your resend/sendgrid/mailgun logic here

	return nil
}