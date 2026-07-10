package domain

import "errors"

// Domain errors
var (
	ErrSocialUserNotFound      = errors.New("social user not found")
	ErrSocialUserAlreadyExists = errors.New("social user already exists")
	ErrProviderNotSupported    = errors.New("provider not supported")
	ErrUserNotFound            = errors.New("user not found")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrPhoneAlreadyExists      = errors.New("phone number already exists")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrEmailNotVerified        = errors.New("email not verified")
	ErrInvalidVerificationCode = errors.New("invalid or expired verification code")
	ErrVerificationNotFound    = errors.New("verification code not found")
	ErrInvalidToken            = errors.New("invalid or expired token")
	ErrSessionNotFound         = errors.New("session not found")
)
