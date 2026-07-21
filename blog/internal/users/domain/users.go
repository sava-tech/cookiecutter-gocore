package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User entity
type User struct {
	ID             uuid.UUID
	Email          string
	PhoneNumber    string
	FirstName      string
	LastName       string
	FullName       string
	AccountType    string
	HashedPassword string
	EmailVerified  bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TokenPayload entity
type TokenPayload struct {
	ID          uuid.UUID
	Email       string
	AccountType string
	IssuedAt    time.Time
	ExpiredAt   time.Time
}

type UserRepository interface {
	// User operations
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByPhoneNumber(ctx context.Context, phone string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	VerifyUserEmail(ctx context.Context, email string) error
	UpdateUserPassword(ctx context.Context, email, hashedPassword string) error
}

// Service interfaces
type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) error
}

type TokenService interface {
	CreateToken(email, accountType string, duration time.Duration) (string, *TokenPayload, error)
	VerifyToken(token string) (*TokenPayload, error)
}

type EmailService interface {
	SendEmailOTP(email, code string) (string, error)
	SendPasswordReset(email, resetLink string) error
}
