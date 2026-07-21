package emailer

type Mailer interface {
	SendEmailOTP(identifier string, token string) (string, error)
	SendWelcomeMessage(identifier string) (string, error)
}
