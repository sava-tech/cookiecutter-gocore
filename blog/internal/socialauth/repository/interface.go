package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/michaelassa01/blog/internal/socialauth/domain"
)

type SocialAuthRepository interface {
    Create(ctx context.Context, socialAuthUser *domain.SocialAuthUser) error
    GetByProviderID(ctx context.Context, provider, providerID string) (*domain.SocialAuthUser, error)
    GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.SocialAuthUser, error)
    Update(ctx context.Context, socialAuthUser *domain.SocialAuthUser) error
}