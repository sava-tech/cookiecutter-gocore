package domain

import (
    "context"
    "github.com/google/uuid"
)

type SocialAuthRepository interface {
    Create(ctx context.Context, socialUser *SocialAuthUser) error
    GetByProviderID(ctx context.Context, provider, providerID string) (*SocialAuthUser, error)
    GetByUserID(ctx context.Context, userID uuid.UUID) (*SocialAuthUser, error)
    Update(ctx context.Context, socialUser *SocialAuthUser) error
}