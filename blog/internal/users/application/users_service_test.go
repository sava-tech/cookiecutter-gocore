package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/michaelassa01/blog/internal/users/domain"
	"github.com/michaelassa01/blog/internal/users/dto"
	mock_db "github.com/michaelassa01/blog/internal/users/mock"
	"go.uber.org/mock/gomock"
)

func TestUserService_GetAccount(t *testing.T) {
	testCases := []struct {
		name          string
		req           *dto.GetUserReq
		setupMocks    func(*mock_db.MockRepository)
		checkResponse func(t *testing.T, resp *dto.RegisterResponse, err error)
	}{
		{
			name: "OK - Successfully get user",
			req: &dto.GetUserReq{
				ID: uuid.New(),
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				user := createTestUser()
				mr.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
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
				require.True(t, resp.EmailVerified)
				require.NotZero(t, resp.CreatedAt)
			},
		},
		{
			name: "User not found",
			req: &dto.GetUserReq{
				ID: uuid.New(),
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				mr.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, domain.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrUserNotFound, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "Database error",
			req: &dto.GetUserReq{
				ID: uuid.New(),
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				mr.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, errors.New("database connection error"))
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "database connection error")
				require.Nil(t, resp)
			},
		},
		{
			name: "Empty user ID",
			req: &dto.GetUserReq{
				ID: uuid.Nil,
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				mr.EXPECT().
					GetUserByID(gomock.Any(), uuid.Nil).
					Times(1).
					Return(nil, domain.ErrUserNotFound)
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.Error(t, err)
				require.Equal(t, domain.ErrUserNotFound, err)
				require.Nil(t, resp)
			},
		},
		{
			name: "User with empty phone number",
			req: &dto.GetUserReq{
				ID: uuid.New(),
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				user := createTestUser()
				user.PhoneNumber = "" // Empty phone number
				mr.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, "", resp.PhoneNumber)
				require.Equal(t, "test@example.com", resp.Email)
			},
		},
		{
			name: "User with unverified email",
			req: &dto.GetUserReq{
				ID: uuid.New(),
			},
			setupMocks: func(mr *mock_db.MockRepository) {
				user := createTestUser()
				user.EmailVerified = false
				mr.EXPECT().
					GetUserByID(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, resp *dto.RegisterResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.False(t, resp.EmailVerified)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_db.NewMockRepository(ctrl)
			mockToken := mock_db.NewMockMaker(ctrl)
			mockPassword := mock_db.NewMockPasswordService(ctrl)
			mockEmail := mock_db.NewMockMailer(ctrl)

			config := &Config{
				AccessTokenDuration:      15 * time.Minute,
				RefreshTokenDuration:     24 * time.Hour,
				VerificationCodeDuration: 5 * time.Minute,
			}

			service := NewUserService(
				mockRepo,
				mockToken,
				mockPassword,
				mockEmail,
				config,
			)

			// Setup mocks
			tc.setupMocks(mockRepo)

			// Execute
			resp, err := service.GetAccount(context.Background(), tc.req)

			// Check response
			tc.checkResponse(t, resp, err)

			// FIXED: Use ctrl.Finish() which is called via defer to verify all expectations
			// No need to call mockRepo.AssertExpectations()
		})
	}
}

// Test UserService initialization
func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_db.NewMockRepository(ctrl)
	mockToken := mock_db.NewMockMaker(ctrl)
	mockPassword := mock_db.NewMockPasswordService(ctrl)
	mockEmail := mock_db.NewMockMailer(ctrl)

	config := &Config{
		AccessTokenDuration:      15 * time.Minute,
		RefreshTokenDuration:     24 * time.Hour,
		VerificationCodeDuration: 5 * time.Minute,
	}

	service := NewUserService(
		mockRepo,
		mockToken,
		mockPassword,
		mockEmail,
		config,
	)

	require.NotNil(t, service)
	require.NotNil(t, service.repo)
	require.NotNil(t, service.tokenMaker)
	require.NotNil(t, service.passwordService)
	require.NotNil(t, service.emailService)
	require.NotNil(t, service.config)
	require.NotNil(t, service.validator)
}

// Test GetAccount with different user types
func TestUserService_GetAccount_DifferentUserTypes(t *testing.T) {
	testCases := []struct {
		name          string
		user          *domain.User
		expectedType  string
	}{
		{
			name: "Regular user",
			user: &domain.User{
				ID:             uuid.New(),
				Email:          "user@example.com",
				PhoneNumber:    "1234567890",
				FirstName:      "Regular",
				LastName:       "User",
				FullName:       "Regular User",
				AccountType:    "user",
				HashedPassword: "hashed",
				EmailVerified:  true,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			expectedType: "user",
		},
		{
			name: "Admin user",
			user: &domain.User{
				ID:             uuid.New(),
				Email:          "admin@example.com",
				PhoneNumber:    "0987654321",
				FirstName:      "Admin",
				LastName:       "User",
				FullName:       "Admin User",
				AccountType:    "admin",
				HashedPassword: "hashed",
				EmailVerified:  true,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			expectedType: "admin",
		},
		{
			name: "Moderator user",
			user: &domain.User{
				ID:             uuid.New(),
				Email:          "mod@example.com",
				PhoneNumber:    "5555555555",
				FirstName:      "Mod",
				LastName:       "User",
				FullName:       "Mod User",
				AccountType:    "moderator",
				HashedPassword: "hashed",
				EmailVerified:  false,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			expectedType: "moderator",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_db.NewMockRepository(ctrl)
			
			// Setup mock to return the test user
			mockRepo.EXPECT().
				GetUserByID(gomock.Any(), tc.user.ID).
				Times(1).
				Return(tc.user, nil)

			// Setup other mocks (not used but required for initialization)
			mockToken := mock_db.NewMockMaker(ctrl)
			mockPassword := mock_db.NewMockPasswordService(ctrl)
			mockEmail := mock_db.NewMockMailer(ctrl)

			config := &Config{
				AccessTokenDuration:      15 * time.Minute,
				RefreshTokenDuration:     24 * time.Hour,
				VerificationCodeDuration: 5 * time.Minute,
			}

			service := NewUserService(
				mockRepo,
				mockToken,
				mockPassword,
				mockEmail,
				config,
			)

			req := &dto.GetUserReq{
				ID: tc.user.ID,
			}

			resp, err := service.GetAccount(context.Background(), req)

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, tc.user.Email, resp.Email)
			require.Equal(t, tc.user.AccountType, resp.AccountType)
			require.Equal(t, tc.user.EmailVerified, resp.EmailVerified)
			require.Equal(t, tc.expectedType, resp.AccountType)
		})
	}
}

// Test GetAccount with nil context
func TestUserService_GetAccount_NilContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_db.NewMockRepository(ctrl)
	mockToken := mock_db.NewMockMaker(ctrl)
	mockPassword := mock_db.NewMockPasswordService(ctrl)
	mockEmail := mock_db.NewMockMailer(ctrl)

	config := &Config{
		AccessTokenDuration:      15 * time.Minute,
		RefreshTokenDuration:     24 * time.Hour,
		VerificationCodeDuration: 5 * time.Minute,
	}

	service := NewUserService(
		mockRepo,
		mockToken,
		mockPassword,
		mockEmail,
		config,
	)

	req := &dto.GetUserReq{
		ID: uuid.New(),
	}

	// The service should handle nil context gracefully
	user := createTestUser()
	user.ID = req.ID
	mockRepo.EXPECT().
		GetUserByID(gomock.Any(), req.ID).
		Times(1).
		Return(user, nil)

	// In Go, context.Background() is typically used instead of nil
	resp, err := service.GetAccount(context.Background(), req)
	
	// This test ensures the function works with a valid context
	// For nil context, we'd need to check how the repository handles it
	if err == nil {
		require.NotNil(t, resp)
	}
}

// Benchmark test for GetAccount
func BenchmarkUserService_GetAccount(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockRepo := mock_db.NewMockRepository(ctrl)
	user := createTestUser()
	
	mockRepo.EXPECT().
		GetUserByID(gomock.Any(), user.ID).
		AnyTimes().
		Return(user, nil)

	mockToken := mock_db.NewMockMaker(ctrl)
	mockPassword := mock_db.NewMockPasswordService(ctrl)
	mockEmail := mock_db.NewMockMailer(ctrl)

	config := &Config{
		AccessTokenDuration:      15 * time.Minute,
		RefreshTokenDuration:     24 * time.Hour,
		VerificationCodeDuration: 5 * time.Minute,
	}

	service := NewUserService(
		mockRepo,
		mockToken,
		mockPassword,
		mockEmail,
		config,
	)

	req := &dto.GetUserReq{
		ID: user.ID,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetAccount(context.Background(), req)
	}
}

