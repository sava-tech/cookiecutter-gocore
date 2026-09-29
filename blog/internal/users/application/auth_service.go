package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/michaelassa01/blog/internal/shared/helpers"
	"github.com/michaelassa01/blog/internal/users/domain"
	"github.com/michaelassa01/blog/internal/users/dto"
	"github.com/michaelassa01/blog/pkg/emailer"
	"github.com/michaelassa01/blog/pkg/token"
)

type AuthService struct {
	repo            domain.Repository
	tokenMaker      token.Maker
	passwordService domain.PasswordService
	emailService    emailer.Mailer
	config          *Config
	validator       *validator.Validate
}

func NewAuthService(
	repo domain.Repository,
	tokenMaker token.Maker,
	passwordService domain.PasswordService,
	emailService emailer.Mailer,
	config *Config,
) *AuthService {
	return &AuthService{
		repo:            repo,
		tokenMaker:      tokenMaker,
		passwordService: passwordService,
		emailService:    emailService,
		config:          config,
		validator:       validator.New(),
	}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Check if email exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	// Check if phone exists
	if req.PhoneNumber != "" {
		existingUser, err = s.repo.GetUserByPhoneNumber(ctx, req.PhoneNumber)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		if existingUser != nil {
			return nil, domain.ErrPhoneAlreadyExists
		}
	}

	// Hash password
	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &domain.User{
		Email:          req.Email,
		PhoneNumber:    req.PhoneNumber,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		FullName:       fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		AccountType:    "user",
		HashedPassword: hashedPassword,
		EmailVerified:  false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if _, err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Send verification OTP
	if err := s.sendVerificationOTP(ctx, req.Email, "registration"); err != nil {
		// Log error but don't fail registration
		fmt.Printf("Failed to send OTP: %v\n", err)
		return nil, err
	}

	return &dto.RegisterResponse{
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		FullName:      user.FullName,
		AccountType:   user.AccountType,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest, userAgent, clientIP string) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		fmt.Println("user not found", err)
		return nil, err
	}

	// Check email verification
	if !user.EmailVerified {
		// Resend OTP
		if err := s.sendVerificationOTP(ctx, req.Email, "registration"); err != nil {
			return nil, err
		}
		return nil, domain.ErrEmailNotVerified
	}

	// Verify password
	if err := s.passwordService.ComparePassword(req.Password, user.HashedPassword); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate tokens
	tokenDetails, err := s.generateTokens(ctx, user, userAgent, clientIP)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		SessionID:             tokenDetails.SessionID,
		AccessToken:           tokenDetails.AccessToken,
		AccessTokenExpiresAt:  tokenDetails.AccessTokenExpiresAt,
		RefreshToken:          tokenDetails.RefreshToken,
		RefreshTokenExpiresAt: tokenDetails.RefreshTokenExpiresAt,
		User: dto.RegisterResponse{
			Email:         user.Email,
			PhoneNumber:   user.PhoneNumber,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			FullName:      user.FullName,
			AccountType:   user.AccountType,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) OtpEmailLogin(ctx context.Context, req *dto.OtpEmailLoginReq) error {
	// Check if user exists (don't reveal for security)
	_, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// Return success even if user doesn't exist (security)
			return nil
		}
		return err
	}

	// Send password reset OTP
	return s.sendVerificationOTP(ctx, req.Email, "otp_login")
}

func (s *AuthService) VerifyOtpEmailLogin(ctx context.Context, req *dto.VerifyOtpEmailLoginReq, userAgent, clientIP string) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// Check email verification
	if !user.EmailVerified {
		// Resend OTP
		if err := s.sendVerificationOTP(ctx, req.Email, "otp_login"); err != nil {
			return nil, err
		}
		return nil, domain.ErrEmailNotVerified
	}

	// Verify password
	if err := s.VerifyEmail(ctx, &dto.VerifyEmailRequest{
		Email:            req.Email,
		Code:             req.Code,
		VerificationType: "otp_login",
	}); err != nil {
		return nil, domain.ErrInvalidVerificationCode
	}

	// Generate tokens
	tokenDetails, err := s.generateTokens(ctx, user, userAgent, clientIP)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		SessionID:             tokenDetails.SessionID,
		AccessToken:           tokenDetails.AccessToken,
		AccessTokenExpiresAt:  tokenDetails.AccessTokenExpiresAt,
		RefreshToken:          tokenDetails.RefreshToken,
		RefreshTokenExpiresAt: tokenDetails.RefreshTokenExpiresAt,
		User: dto.RegisterResponse{
			Email:         user.Email,
			PhoneNumber:   user.PhoneNumber,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			FullName:      user.FullName,
			AccountType:   user.AccountType,
			EmailVerified: user.EmailVerified,
			CreatedAt:     user.CreatedAt,
		},
	}, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, req *dto.VerifyEmailRequest) error {
	// Get valid verification code
	// verify the code with the hashed code
	hashedCode, err := s.repo.GetVerificationHash(ctx, req.Email, req.VerificationType)
	if err != nil {
		if errors.Is(err, domain.ErrVerificationNotFound) {
			return domain.ErrInvalidVerificationCode
		}
	}

	// Compare provided code with hashed code
	if helpers.VerifyPINFast(req.Code, hashedCode); err != nil {
		return domain.ErrInvalidVerificationCode
	}

	verification, err := s.repo.GetValidVerificationCode(ctx, req.Email, hashedCode, req.VerificationType)
	if err != nil {
		if errors.Is(err, domain.ErrVerificationNotFound) {
			return domain.ErrInvalidVerificationCode
		}
		return err
	}

	// Mark code as used
	if err := s.repo.MarkVerificationAsUsed(ctx, verification.ID); err != nil {
		return err
	}

	// Verify user email
	if err := s.repo.VerifyUserEmail(ctx, req.Email); err != nil {
		return err
	}

	// Invalidate all other verification codes
	if err := s.repo.InvalidateVerificationCodes(ctx, req.Email, req.VerificationType); err != nil {
		// Log error but don't fail
		fmt.Printf("Failed to invalidate codes: %v\n", err)
	}

	return nil
}

func (s *AuthService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {
	// Check if user exists (don't reveal for security)
	_, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// Return success even if user doesn't exist (security)
			return nil
		}
		return err
	}

	// Send password reset OTP
	return s.sendVerificationOTP(ctx, req.Email, "password_reset")
}

func (s *AuthService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	// Get valid verification code
	// verify the code with the hashed code
	hashedCode, err := s.repo.GetVerificationHash(ctx, req.Email, "password_reset")
	if err != nil {
		if errors.Is(err, domain.ErrVerificationNotFound) {
			return domain.ErrInvalidVerificationCode
		}
	}

	// Compare provided code with hashed code
	if helpers.VerifyPINFast(req.Code, hashedCode); err != nil {
		return domain.ErrInvalidVerificationCode
	}

	verification, err := s.repo.GetValidVerificationCode(ctx, req.Email, hashedCode, "password_reset")
	if err != nil {
		if errors.Is(err, domain.ErrVerificationNotFound) {
			return domain.ErrInvalidVerificationCode
		}
		return err
	}

	// Hash new password
	hashedPassword, err := s.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	if err := s.repo.UpdateUserPassword(ctx, req.Email, hashedPassword); err != nil {
		return err
	}

	// Mark code as used
	if err := s.repo.MarkVerificationAsUsed(ctx, verification.ID); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	// Verify refresh token
	payload, err := s.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Get session
	session, err := s.repo.GetSessionByID(ctx, payload.ID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	// Check if session is valid
	if session.IsBlocked || session.RefreshToken != req.RefreshToken || session.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrInvalidToken
	}

	// Get user
	user, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Email,
		user.AccountType,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		ExpiresAt:    accessPayload.ExpiredAt,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return s.repo.DeleteSession(ctx, sessionID)
}

// Private helper methods
func (s *AuthService) sendVerificationOTP(ctx context.Context, identifier, verificationType string) error {
	// TODO: add validation for the verification type
	// Generate random code
	codePin, err := helpers.GeneratePIN(4) // Assuming you have a utility function to generate a random 4-digit
	if err != nil {
		return err
	}
	// Hash the code for storage
	hashedCode := helpers.HashPINFast(codePin)

	// Create verification record
	verification := &domain.Verification{
		Code:             hashedCode,
		IdentifierType:   "email",
		Identifier:       identifier,
		VerificationType: verificationType,
		ExpiredAt:        time.Now().Add(s.config.VerificationCodeDuration),
		Used:             false,
		CreatedAt:        time.Now(),
	}

	if err := s.repo.CreateVerification(ctx, verification); err != nil {
		return err
	}

	// Send email
	_, err = s.emailService.SendEmailOTP(identifier, codePin)
	if err != nil {
		return err
	}

	return err
}

func (s *AuthService) generateTokens(ctx context.Context, user *domain.User, userAgent, clientIP string) (*dto.LoginResponse, error) {
	// Create access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Email,
		user.AccountType,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Email,
		user.AccountType,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	// Create session
	session := &domain.Session{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		Email:        user.Email,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}, nil
}
