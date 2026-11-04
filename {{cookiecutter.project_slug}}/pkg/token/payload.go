package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken  = errors.New("token has expired")
	ErrInvalideToken = errors.New("token is invalide")
)

// Payload contains the payload data of the token
type Payload struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	AccountType string    `json:"account_type"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(email string, accountType string, duration time.Duration) (*Payload, error) {
	// Generate random uuid
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:          tokenID,
		Email:       email,
		AccountType: accountType,
		IssuedAt:    time.Now(),
		ExpiredAt:   time.Now().Add(duration),
	}

	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

// Validate implements jwt.ClaimsValidator so JWT verification applies the same
// expiry check as Valid and returns ErrExpiredToken
func (payload *Payload) Validate() error {
	return payload.Valid()
}

// The registered-claim getters below are required by jwt.Claims. They return
// nil because expiry is carried in the custom expired_at claim and checked by
// Validate, which keeps the token format unchanged.

func (payload *Payload) GetExpirationTime() (*jwt.NumericDate, error) { return nil, nil }
func (payload *Payload) GetIssuedAt() (*jwt.NumericDate, error)       { return nil, nil }
func (payload *Payload) GetNotBefore() (*jwt.NumericDate, error)      { return nil, nil }
func (payload *Payload) GetIssuer() (string, error)                   { return "", nil }
func (payload *Payload) GetSubject() (string, error)                  { return "", nil }
func (payload *Payload) GetAudience() (jwt.ClaimStrings, error)       { return nil, nil }
