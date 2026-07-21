package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Request DTOs
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyEmailRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Code             string `json:"code" binding:"required,len=4"`
	VerificationType string `json:"verification_type" binding:"required"`
}

type OtpEmailLoginReq struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyOtpEmailLoginReq struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=4"`
}

type OtpPhnLoginReq struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,len=4"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Response DTOs
type RegisterResponse struct {
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	FullName      string    `json:"full_name"`
	AccountType   string    `json:"account_type"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	PhoneNumber   string    `json:"phone_number"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	FullName      string    `json:"full_name"`
	AccountType   string    `json:"account_type"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

type LoginResponse struct {
	SessionID             uuid.UUID        `json:"session_id"`
	AccessToken           string           `json:"access_token"`
	AccessTokenExpiresAt  time.Time        `json:"access_token_expires_at"`
	RefreshToken          string           `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time        `json:"refresh_token_expires_at"`
	User                  RegisterResponse `json:"user"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type GetUserReq struct {
	ID uuid.UUID `json:"id" binding:"required,uuid"`
}

type GetUserRes struct {
	ID          pgtype.UUID `json:"id"`
	Email       string      `json:"email"`
	Username    string      `json:"username"`
	PhoneNumber string      `json:"phone_number"`
	Avatar      string      `json:"avatar"`
	Age         int32       `json:"age"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type UpdateUserReq struct {
	Email       string `json:"email" binding:"omitempty,email"`
	Username    string `json:"username" binding:"omitempty"`
	PhoneNumber string `json:"phone_number" binding:"omitempty"`
	Avatar      string `json:"avatar" binding:"omitempty"`
	Age         int32  `json:"age" binding:"omitempty,gte=13"`
}

type UpdateUserRes struct {
	ID          pgtype.UUID `json:"id"`
	Email       string      `json:"email"`
	Username    string      `json:"username"`
	PhoneNumber string      `json:"phone_number"`
	Avatar      string      `json:"avatar"`
	Age         int32       `json:"age"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ListUsersReq struct {
	Limit  int32 `json:"limit" binding:"omitempty,min=1,max=100"`
	Offset int32 `json:"offset" binding:"omitempty,min=0"`
}

type ListUsersRes struct {
	Users []GetUserRes `json:"users"`
}

type DeleteUserReq struct {
	ID string `json:"id" binding:"required,uuid"`
}
