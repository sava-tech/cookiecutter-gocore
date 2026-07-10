package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/socialauth/domain"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/socialauth/dto"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/socialauth/repository"
	userdomain "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/domain"

	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/token"
)

type SocialAuthService struct {
	socialRepo repository.SocialAuthRepository
	userRepo   userdomain.Repository
	tokenMaker token.Maker
	config     *Config
}

type Config struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	CallbackURL          string
}

func NewSocialAuthService(
	socialRepo repository.SocialAuthRepository,
	userRepo userdomain.Repository,
	tokenMaker token.Maker,
	config *Config,
) *SocialAuthService {
	return &SocialAuthService{
		socialRepo: socialRepo,
		userRepo:   userRepo,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

// ProcessSocialUser handles the OAuth callback and creates/links user
func (s *SocialAuthService) ProcessSocialUser(ctx context.Context, provider string, gothUser goth.User, userAgent, clientIP string) (*dto.SocialAuthResponse, error) {

	// Check if social user exists
	socialUser, err := s.socialRepo.GetByProviderID(ctx, provider, gothUser.UserID)
	if err != nil && !errors.Is(err, domain.ErrSocialUserNotFound) {
		return nil, err
	}

	var userID uuid.UUID

	if socialUser != nil {
		// Social user exists - use linked user
		userID = socialUser.UserID
	} else {
		// Check if user exists by email
		existingUser, err := s.userRepo.GetUserByEmail(ctx, gothUser.Email)
		if err != nil && !errors.Is(err, userdomain.ErrUserNotFound) {
			return nil, err
		}

		if existingUser != nil {
			// Link social account to existing user
			userID = existingUser.ID
		} else {
			// Create new user
			gothParams := &userdomain.User{
				Email:       gothUser.Email,
				FirstName:   gothUser.FirstName,
				LastName:    gothUser.LastName,
				FullName:    gothUser.Name,
				AccountType: "user",
				// TODO: ensure email is verified from provider
				EmailVerified: true, //which can be false and redirect for verification
			}
			newUser, err := s.userRepo.CreateUser(ctx, gothParams)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
			userID = newUser.ID
		}

		// Create social user record
		newSocialUser := &domain.SocialAuthUser{
			ID:          uuid.New(),
			UserID:      userID,
			Provider:    provider,
			ProviderID:  gothUser.UserID,
			Email:       gothUser.Email,
			Name:        gothUser.Name,
			AvatarURL:   gothUser.AvatarURL,
			AccessToken: gothUser.AccessToken,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.socialRepo.Create(ctx, newSocialUser); err != nil {
			return nil, fmt.Errorf("failed to create social user: %w", err)
		}

		socialUser = newSocialUser
	}

	// Generate Pasto or JWT tokens (using your existing token system)
	tokenDetails, err := s.generateTokens(ctx, socialUser, userAgent, clientIP)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &dto.SocialAuthResponse{
		SessionID:             tokenDetails.SessionID,
		AccessToken:           tokenDetails.AccessToken,
		AccessTokenExpiresAt:  tokenDetails.AccessTokenExpiresAt,
		RefreshToken:          tokenDetails.RefreshToken,
		RefreshTokenExpiresAt: tokenDetails.RefreshTokenExpiresAt,
		User: dto.SocialAuthUserResponse{
			ID:        socialUser.ID,
			UserID:    socialUser.UserID,
			Provider:  socialUser.Provider,
			Email:     socialUser.Email,
			Name:      socialUser.Name,
			AvatarURL: socialUser.AvatarURL,
			CreatedAt: socialUser.CreatedAt,
		},
	}, nil
}

func (s *SocialAuthService) generateTokens(ctx context.Context, user *domain.SocialAuthUser, userAgent, clientIP string) (*dto.SocialAuthResponse, error) {
	// Create access token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.Email,
		"user",
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.Email,
		"user",
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	// Create session
	session := &userdomain.Session{
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

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return &dto.SocialAuthResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}, nil
}
