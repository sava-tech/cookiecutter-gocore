package emailer

type SMSender interface {
	// SendPhoneOTP OTP token or String of Pass code : SendPhoneOTP
	// If pass code is empty , it will auto-generate it
	SendPhoneOTP(identifier string, token string) (string, error)

	// SendPhoneWelcomeMessage Welcome message can either be sent to the user
	// through and email or phone number
	// TODO:: Whats Welcome message
	SendPhoneWelcomeMessage(identifier string) (string, error)
}
