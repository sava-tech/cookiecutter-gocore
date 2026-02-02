// Interface definition

package emailer

type Mailer interface {
	// SendEmailOTP OTP token or String of Pass code : SendOTP
	// If pass code is empty , it will auto-generate it
	SendEmailOTP(identifier string, token string) (string, error)

	// SendWelcomeMessage Welcome message can either be sent to the user
	// through and email or phone number
	// TODO:: Whats Welcome message
	SendWelcomeMessage(identifier string) (string, error)
}
