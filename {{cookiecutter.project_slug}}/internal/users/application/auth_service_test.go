package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/shared/helpers"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/domain"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/dto"
	mock_db "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/mock"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/token"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Test helpers
func createTestUser() *domain.User {
	return &domain.User{
		ID:             uuid.New(),
		Email:          "test@example.com",
		PhoneNumber:    "1234567890",
		FirstName:      "John",
		LastName:       "Doe",
		FullName:       "John Doe",
		AccountType:    "user",
		HashedPassword: "hashed_password",
		EmailVerified:  true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func createTestConfig() *Config {
	return &Config{
		AccessTokenDuration:      15 * time.Minute,
		RefreshTokenDuration:     24 * time.Hour,
		VerificationCodeDuration: 5 * time.Minute,
	}
}

func setupTestService(t *testing.T) (auth *AuthService, repo *mock_db.MockRepository, tokenMaker *mock_db.MockMaker, password *mock_db.MockPasswordService, emailier *mock_db.MockMailer) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_db.NewMockRepository(ctrl)
	mockToken := mock_db.NewMockMaker(ctrl)
	mockPassword := mock_db.NewMockPasswordService(ctrl)
	mockEmail := mock_db.NewMockMailer(ctrl)

	config := createTestConfig()

	service := NewAuthService(
		mockRepo,
		mockToken,
		mockPassword,
		mockEmail,
		config,
	)

	return service, mockRepo, mockToken, mockPassword, mockEmail
}

// Test Register
func TestAuthService_Register(t *testing.T) {
	testCases := []struct {
		name          string
		req           *dto.RegisterRequest
		setupMocks    func(*mock_db.MockRepository, *mock_db.MockPasswordService, *mock_db.MockMailer)
		checkResponse func(t *testing.T, resp *dto.RegisterResponse, err error)
	}{
		{
			name: "OK - Successful registration",
			req: &dto.RegisterRequest{
				Email:       "test@example.com",
				PhoneNumber: "1234567890",
				FirstName:   "John",
				LastName:    "Doe",
				Password:    "password123",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(nil, domain.ErrUserNotFound)

				mr.EXPECT().
					GetUserByPhoneNumber(gomock.Any(), "1234567890").
					Times(1).
					Return(nil, domain.ErrUserNotFound)

				mp.EXPECT().
					HashPassword("password123").
					Times(1).
					Return("hashed_password", nil)

				mr.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx context.Context, user *domain.User) (*domain.User, error) {
						user.ID = uuid.New()
						return user, nil
					})

				mr.EXPECT().
					CreateVerification(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				me.EXPECT().
					SendEmailOTP("test@example.com", gomock.Any()).
					Times(1).
					Return("message_id", nil)
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "test@example.com", resp.Email)
				require.Equal(t, "1234567890", resp.PhoneNumber)
				require.Equal(t, "John", resp.FirstName)
				require.Equal(t, "Doe", resp.LastName)
				require.Equal(t, "John Doe", resp.FullName)
				require.Equal(t, "user", resp.AccountType)
				require.False(t, resp.EmailVerified)
			},
		},
		{
			name: "Email already exists",
			req: &dto.RegisterRequest{
				Email: "existing@example.com",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				existingUser := createTestUser()
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "existing@example.com").
					Times(1).
					Return(existingUser, nil)
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrEmailAlreadyExists, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "Phone number already exists",
			req: &dto.RegisterRequest{
				Email:       "test@example.com",
				PhoneNumber: "existing_phone",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(nil, domain.ErrUserNotFound)

				existingUser := createTestUser()
				mr.EXPECT().
					GetUserByPhoneNumber(gomock.Any(), "existing_phone").
					Times(1).
					Return(existingUser, nil)
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrPhoneAlreadyExists, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "Password hashing fails",
			req: &dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(nil, domain.ErrUserNotFound)

				mr.EXPECT().
					GetUserByPhoneNumber(gomock.Any(), "").
					Times(1).
					Return(nil, domain.ErrUserNotFound)

				mp.EXPECT().
					HashPassword("password123").
					Times(1).
					Return("", errors.New("hashing error"))
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "hashing error")
				require.Nil(t, resp)
			},
		},
		{
			name: "Create user fails",
			req: &dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(nil, domain.ErrUserNotFound)

				mr.EXPECT().
					GetUserByPhoneNumber(gomock.Any(), "").
					Times(1).
					Return(nil, domain.ErrUserNotFound)

				mp.EXPECT().
					HashPassword("password123").
					Times(1).
					Return("hashed_password", nil)

				mr.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("database error"))
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "database error")
				require.Nil(t, resp)
			},
		},
		{
			name: "Send OTP fails",
			req: &dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(nil, domain.ErrUserNotFound)

				mr.EXPECT().
					GetUserByPhoneNumber(gomock.Any(), "").
					Times(1).
					Return(nil, domain.ErrUserNotFound)

				mp.EXPECT().
					HashPassword("password123").
					Times(1).
					Return("hashed_password", nil)

				mr.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(func(ctx context.Context, user *domain.User) (*domain.User, error) {
						user.ID = uuid.New()
						return user, nil
					})

				mr.EXPECT().
					CreateVerification(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("verification creation failed"))
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "verification creation failed")
				require.Nil(t, resp)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, mockRepo, _, mockPassword, mockEmail := setupTestService(t)
			tc.setupMocks(mockRepo, mockPassword, mockEmail)

			resp, err := service.Register(context.Background(), tc.req)
			tc.checkResponse(t, resp, err)
		})
	}
}

// Test Login
func TestAuthService_Login(t *testing.T) {
	testCases := []struct {
		name          string
		req           *dto.LoginRequest
		setupMocks    func(*mock_db.MockRepository, *mock_db.MockMaker, *mock_db.MockPasswordService, *mock_db.MockMailer)
		checkResponse func(t *testing.T, resp *dto.LoginResponse, err error)
	}{
		{
			name: "OK - Successful login",
			req: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				user := createTestUser()
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(user, nil)

				mp.EXPECT().
					ComparePassword("password123", user.HashedPassword).
					Times(1).
					Return(nil)

				// Mock token generation
				accessPayload := &token.Payload{
					ID:        uuid.New(),
					ExpiredAt: time.Now().Add(15 * time.Minute),
				}
				refreshPayload := &token.Payload{
					ID:        uuid.New(),
					ExpiredAt: time.Now().Add(24 * time.Hour),
				}

				mt.EXPECT().
					CreateToken(user.Email, user.AccountType, 15*time.Minute).
					Times(1).
					Return("access_token", accessPayload, nil)

				mt.EXPECT().
					CreateToken(user.Email, user.AccountType, 24*time.Hour).
					Times(1).
					Return("refresh_token", refreshPayload, nil)

				mr.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, resp *dto.LoginResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotEmpty(t, resp.AccessToken)
				require.NotEmpty(t, resp.RefreshToken)
				require.NotEmpty(t, resp.SessionID)
				require.NotNil(t, resp.User)
				require.Equal(t, "test@example.com", resp.User.Email)
			},
		},
		{
			name: "User not found",
			req: &dto.LoginRequest{
				Email: "nonexistent@example.com",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "nonexistent@example.com").
					Times(1).
					Return(nil, domain.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, resp *dto.LoginResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidCredentials, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "Email not verified",
			req: &dto.LoginRequest{
				Email: "unverified@example.com",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				user := createTestUser()
				user.EmailVerified = false
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "unverified@example.com").
					Times(1).
					Return(user, nil)

				mr.EXPECT().
					CreateVerification(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				me.EXPECT().
					SendEmailOTP("unverified@example.com", gomock.Any()).
					Times(1).
					Return("message_id", nil)
			},
			checkResponse: func(t *testing.T, resp *dto.LoginResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrEmailNotVerified, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "Invalid password",
			req: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrong_password",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				user := createTestUser()
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(user, nil)

				mp.EXPECT().
					ComparePassword("wrong_password", user.HashedPassword).
					Times(1).
					Return(errors.New("invalid password"))
			},
			checkResponse: func(t *testing.T, resp *dto.LoginResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidCredentials, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "Create session fails",
			req: &dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker, mp *mock_db.MockPasswordService, me *mock_db.MockMailer) {
				user := createTestUser()
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(user, nil)

				mp.EXPECT().
					ComparePassword("password123", user.HashedPassword).
					Times(1).
					Return(nil)

				accessPayload := &token.Payload{
					ID:        uuid.New(),
					ExpiredAt: time.Now().Add(15 * time.Minute),
				}
				refreshPayload := &token.Payload{
					ID:        uuid.New(),
					ExpiredAt: time.Now().Add(24 * time.Hour),
				}

				mt.EXPECT().
					CreateToken(user.Email, user.AccountType, 15*time.Minute).
					Times(1).
					Return("access_token", accessPayload, nil)

				mt.EXPECT().
					CreateToken(user.Email, user.AccountType, 24*time.Hour).
					Times(1).
					Return("refresh_token", refreshPayload, nil)

				mr.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("session creation failed"))
			},
			checkResponse: func(t *testing.T, resp *dto.LoginResponse, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "session creation failed")
				require.Nil(t, resp)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, mockRepo, mockToken, mockPassword, mockEmail := setupTestService(t)
			tc.setupMocks(mockRepo, mockToken, mockPassword, mockEmail)

			resp, err := service.Login(context.Background(), tc.req, "test-agent", "127.0.0.1")
			tc.checkResponse(t, resp, err)
		})
	}
}

// Test VerifyEmail
func TestAuthService_VerifyEmail(t *testing.T) {
	testCases := []struct {
		name          string
		req           *dto.VerifyEmailRequest
		setupMocks    func(*mock_db.MockRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK - Successful verification",
			req: &dto.VerifyEmailRequest{
				Email: "test@example.com",
				Code:  "1234",
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				hashedCode := helpers.HashPINFast("1234")
				verification := &domain.Verification{
					ID:               uuid.New(),
					Code:             hashedCode,
					Identifier:       "test@example.com",
					VerificationType: "registration",
					Used:             false,
					ExpiredAt:        time.Now().Add(5 * time.Minute),
				}

				mr.EXPECT().
					GetVerificationHash(gomock.Any(), "test@example.com", "registration").
					Times(1).
					Return(hashedCode, nil)

				mr.EXPECT().
					GetValidVerificationCode(gomock.Any(), "test@example.com", hashedCode, "registration").
					Times(1).
					Return(verification, nil)

				mr.EXPECT().
					MarkVerificationAsUsed(gomock.Any(), verification.ID).
					Times(1).
					Return(nil)

				mr.EXPECT().
					VerifyUserEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(nil)

				mr.EXPECT().
					InvalidateVerificationCodes(gomock.Any(), "test@example.com", "registration").
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Invalid verification - verification not found after code match",
			req: &dto.VerifyEmailRequest{
				Email: "test@example.com",
				Code:  "1234",
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				hashedCode := helpers.HashPINFast("1234")
				mr.EXPECT().
					GetVerificationHash(gomock.Any(), "test@example.com", "registration").
					Times(1).
					Return(hashedCode, nil)

				// Now expect GetValidVerificationCode to be called
				mr.EXPECT().
					GetValidVerificationCode(gomock.Any(), "test@example.com", hashedCode, "registration").
					Times(1).
					Return(nil, domain.ErrVerificationNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidVerificationCode, err)
			},
		},
		{
			name: "Verification not found",
			req: &dto.VerifyEmailRequest{
				Email: "test@example.com",
				Code:  "1234",
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				hashedCode := helpers.HashPINFast("1234")
				mr.EXPECT().
					GetVerificationHash(gomock.Any(), "test@example.com", "registration").
					Times(1).
					Return(hashedCode, nil)

				mr.EXPECT().
					GetValidVerificationCode(gomock.Any(), "test@example.com", hashedCode, "registration").
					Times(1).
					Return(nil, domain.ErrVerificationNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidVerificationCode, err)
			},
		},
		{
			name: "Mark verification as used fails",
			req: &dto.VerifyEmailRequest{
				Email: "test@example.com",
				Code:  "1234",
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				hashedCode := helpers.HashPINFast("1234")
				verification := &domain.Verification{
					ID:               uuid.New(),
					Code:             hashedCode,
					Identifier:       "test@example.com",
					VerificationType: "registration",
					Used:             false,
					ExpiredAt:        time.Now().Add(5 * time.Minute),
				}

				mr.EXPECT().
					GetVerificationHash(gomock.Any(), "test@example.com", "registration").
					Times(1).
					Return(hashedCode, nil)

				mr.EXPECT().
					GetValidVerificationCode(gomock.Any(), "test@example.com", hashedCode, "registration").
					Times(1).
					Return(verification, nil)

				mr.EXPECT().
					MarkVerificationAsUsed(gomock.Any(), verification.ID).
					Times(1).
					Return(errors.New("mark failed"))
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "mark failed")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, mockRepo, _, _, _ := setupTestService(t)
			tc.setupMocks(mockRepo)

			err := service.VerifyEmail(context.Background(), tc.req)
			tc.checkResponse(t, err)
		})
	}
}

// Test ForgotPassword
func TestAuthService_ForgotPassword(t *testing.T) {
	testCases := []struct {
		name          string
		req           *dto.ForgotPasswordRequest
		setupMocks    func(*mock_db.MockRepository, *mock_db.MockMailer)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK - Successful forgot password",
			req: &dto.ForgotPasswordRequest{
				Email: "test@example.com",
			},
			setupMocks: func(mr *mock_db.MockRepository, me *mock_db.MockMailer) {
				user := createTestUser()
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "test@example.com").
					Times(1).
					Return(user, nil)

				mr.EXPECT().
					CreateVerification(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				me.EXPECT().
					SendEmailOTP("test@example.com", gomock.Any()).
					Times(1).
					Return("message_id", nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "User not found - returns success for security",
			req: &dto.ForgotPasswordRequest{
				Email: "nonexistent@example.com",
			},
			setupMocks: func(mr *mock_db.MockRepository, me *mock_db.MockMailer) {
				mr.EXPECT().
					GetUserByEmail(gomock.Any(), "nonexistent@example.com").
					Times(1).
					Return(nil, domain.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err) // Should return nil for security
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, mockRepo, _, _, mockEmail := setupTestService(t)
			tc.setupMocks(mockRepo, mockEmail)

			err := service.ForgotPassword(context.Background(), tc.req)
			tc.checkResponse(t, err)
		})
	}
}

// Test ResetPassword
func TestAuthService_ResetPassword(t *testing.T) {
	testCases := []struct {
		name          string
		req           *dto.ResetPasswordRequest
		setupMocks    func(*mock_db.MockRepository, *mock_db.MockPasswordService)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name: "OK - Successful password reset",
			req: &dto.ResetPasswordRequest{
				Email:       "test@example.com",
				Code:        "1234",
				NewPassword: "newpassword123",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService) {
				hashedCode := helpers.HashPINFast("1234")
				verification := &domain.Verification{
					ID:               uuid.New(),
					Code:             hashedCode,
					Identifier:       "test@example.com",
					VerificationType: "password_reset",
					Used:             false,
					ExpiredAt:        time.Now().Add(5 * time.Minute),
				}

				mr.EXPECT().
					GetVerificationHash(gomock.Any(), "test@example.com", "password_reset").
					Times(1).
					Return(hashedCode, nil)

				mr.EXPECT().
					GetValidVerificationCode(gomock.Any(), "test@example.com", hashedCode, "password_reset").
					Times(1).
					Return(verification, nil)

				mp.EXPECT().
					HashPassword("newpassword123").
					Times(1).
					Return("new_hashed_password", nil)

				mr.EXPECT().
					UpdateUserPassword(gomock.Any(), "test@example.com", "new_hashed_password").
					Times(1).
					Return(nil)

				mr.EXPECT().
					MarkVerificationAsUsed(gomock.Any(), verification.ID).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Invalid verification code",
			req: &dto.ResetPasswordRequest{
				Email: "test@example.com",
				Code:  "9999",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService) {
				hashedCode := helpers.HashPINFast("1234")
				mr.EXPECT().
					GetVerificationHash(gomock.Any(), "test@example.com", "password_reset").
					Times(1).
					Return(hashedCode, nil)

				// Now expect GetValidVerificationCode to be called
				mr.EXPECT().
					GetValidVerificationCode(gomock.Any(), "test@example.com", hashedCode, "password_reset").
					Times(1).
					Return(nil, domain.ErrVerificationNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidVerificationCode, err)
			},
		},
		{
			name: "Verification not found",
			req: &dto.ResetPasswordRequest{
				Email: "test@example.com",
				Code:  "1234",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService) {
				hashedCode := helpers.HashPINFast("1234")
				mr.EXPECT().
					GetVerificationHash(gomock.Any(), "test@example.com", "password_reset").
					Times(1).
					Return(hashedCode, nil)

				mr.EXPECT().
					GetValidVerificationCode(gomock.Any(), "test@example.com", hashedCode, "password_reset").
					Times(1).
					Return(nil, domain.ErrVerificationNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidVerificationCode, err)
			},
		},
		{
			name: "Hash new password fails",
			req: &dto.ResetPasswordRequest{
				Email:       "test@example.com",
				Code:        "1234",
				NewPassword: "newpassword123",
			},
			setupMocks: func(mr *mock_db.MockRepository, mp *mock_db.MockPasswordService) {
				hashedCode := helpers.HashPINFast("1234")
				verification := &domain.Verification{
					ID:               uuid.New(),
					Code:             hashedCode,
					Identifier:       "test@example.com",
					VerificationType: "password_reset",
					Used:             false,
					ExpiredAt:        time.Now().Add(5 * time.Minute),
				}

				mr.EXPECT().
					GetVerificationHash(gomock.Any(), "test@example.com", "password_reset").
					Times(1).
					Return(hashedCode, nil)

				mr.EXPECT().
					GetValidVerificationCode(gomock.Any(), "test@example.com", hashedCode, "password_reset").
					Times(1).
					Return(verification, nil)

				mp.EXPECT().
					HashPassword("newpassword123").
					Times(1).
					Return("", errors.New("hashing error"))
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "hashing error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, mockRepo, _, mockPassword, _ := setupTestService(t)
			tc.setupMocks(mockRepo, mockPassword)

			err := service.ResetPassword(context.Background(), tc.req)
			tc.checkResponse(t, err)
		})
	}
}

// Test RefreshToken
func TestAuthService_RefreshToken(t *testing.T) {
	testCases := []struct {
		name          string
		req           *dto.RefreshTokenRequest
		setupMocks    func(*mock_db.MockRepository, *mock_db.MockMaker)
		checkResponse func(t *testing.T, resp *dto.TokenResponse, err error)
	}{
		{
			name: "OK - Successful refresh",
			req: &dto.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker) {
				sessionID := uuid.New()
				userID := uuid.New()
				payload := &token.Payload{
					ID: sessionID,
				}

				mt.EXPECT().
					VerifyToken("valid_refresh_token").
					Times(1).
					Return(payload, nil)

				session := &domain.Session{
					ID:           sessionID,
					UserID:       userID,
					RefreshToken: "valid_refresh_token",
					IsBlocked:    false,
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				}

				mr.EXPECT().
					GetSessionByID(gomock.Any(), sessionID).
					Times(1).
					Return(session, nil)

				user := createTestUser()
				user.ID = userID
				mr.EXPECT().
					GetUserByID(gomock.Any(), userID).
					Times(1).
					Return(user, nil)

				accessPayload := &token.Payload{
					ID:        uuid.New(),
					ExpiredAt: time.Now().Add(15 * time.Minute),
				}

				mt.EXPECT().
					CreateToken(user.Email, user.AccountType, 15*time.Minute).
					Times(1).
					Return("new_access_token", accessPayload, nil)
			},
			checkResponse: func(t *testing.T, resp *dto.TokenResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "new_access_token", resp.AccessToken)
				require.Equal(t, "valid_refresh_token", resp.RefreshToken)
				require.NotNil(t, resp.ExpiresAt)
			},
		},
		{
			name: "Invalid refresh token",
			req: &dto.RefreshTokenRequest{
				RefreshToken: "invalid_token",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker) {
				mt.EXPECT().
					VerifyToken("invalid_token").
					Times(1).
					Return(nil, errors.New("invalid token"))
			},
			checkResponse: func(t *testing.T, resp *dto.TokenResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidToken, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "Session not found",
			req: &dto.RefreshTokenRequest{
				RefreshToken: "valid_token",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker) {
				sessionID := uuid.New()
				payload := &token.Payload{
					ID: sessionID,
				}

				mt.EXPECT().
					VerifyToken("valid_token").
					Times(1).
					Return(payload, nil)

				mr.EXPECT().
					GetSessionByID(gomock.Any(), sessionID).
					Times(1).
					Return(nil, domain.ErrSessionNotFound)
			},
			checkResponse: func(t *testing.T, resp *dto.TokenResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidToken, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "Blocked session",
			req: &dto.RefreshTokenRequest{
				RefreshToken: "blocked_token",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker) {
				sessionID := uuid.New()
				payload := &token.Payload{
					ID: sessionID,
				}

				mt.EXPECT().
					VerifyToken("blocked_token").
					Times(1).
					Return(payload, nil)

				session := &domain.Session{
					ID:           sessionID,
					RefreshToken: "blocked_token",
					IsBlocked:    true,
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				}

				mr.EXPECT().
					GetSessionByID(gomock.Any(), sessionID).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(t *testing.T, resp *dto.TokenResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidToken, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "Expired session",
			req: &dto.RefreshTokenRequest{
				RefreshToken: "expired_token",
			},
			setupMocks: func(mr *mock_db.MockRepository, mt *mock_db.MockMaker) {
				sessionID := uuid.New()
				payload := &token.Payload{
					ID: sessionID,
				}

				mt.EXPECT().
					VerifyToken("expired_token").
					Times(1).
					Return(payload, nil)

				session := &domain.Session{
					ID:           sessionID,
					RefreshToken: "expired_token",
					IsBlocked:    false,
					ExpiresAt:    time.Now().Add(-1 * time.Hour),
				}

				mr.EXPECT().
					GetSessionByID(gomock.Any(), sessionID).
					Times(1).
					Return(session, nil)
			},
			checkResponse: func(t *testing.T, resp *dto.TokenResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrInvalidToken, err)
				require.Nil(t, resp)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, mockRepo, mockToken, _, _ := setupTestService(t)
			tc.setupMocks(mockRepo, mockToken)

			resp, err := service.RefreshToken(context.Background(), tc.req)
			tc.checkResponse(t, resp, err)
		})
	}
}

// Test Logout
func TestAuthService_Logout(t *testing.T) {
	testCases := []struct {
		name          string
		sessionID     uuid.UUID
		setupMocks    func(*mock_db.MockRepository)
		checkResponse func(t *testing.T, err error)
	}{
		{
			name:      "OK - Successful logout",
			sessionID: uuid.New(),
			setupMocks: func(mr *mock_db.MockRepository) {
				mr.EXPECT().
					DeleteSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name:      "Session not found",
			sessionID: uuid.New(),
			setupMocks: func(mr *mock_db.MockRepository) {
				mr.EXPECT().
					DeleteSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(domain.ErrSessionNotFound)
			},
			checkResponse: func(t *testing.T, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrSessionNotFound, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, mockRepo, _, _, _ := setupTestService(t)
			tc.setupMocks(mockRepo)

			err := service.Logout(context.Background(), tc.sessionID)
			tc.checkResponse(t, err)
		})
	}
}
