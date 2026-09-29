package domain

import (
	"time"

	"github.com/google/uuid"
)

// Supported providers
const (
	ProviderGoogle = "google"
	ProviderApple  = "apple"
)

// IsValidProvider checks if the provider is supported
func IsValidProvider(provider string) bool {
	return provider == ProviderGoogle || provider == ProviderApple
}

type SocialAuthUser struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Provider    string    `json:"provider"`
	ProviderID  string    `json:"provider_id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	AvatarURL   string    `json:"avatar_url"`
	AccessToken string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
