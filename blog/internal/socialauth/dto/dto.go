package dto

import (
    "time"
    "github.com/google/uuid"
)

// Response
type SocialAuthUserResponse struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"user_id"`
    Provider  string    `json:"provider"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    AvatarURL string    `json:"avatar_url"`
    CreatedAt time.Time `json:"created_at"`
}


type SocialAuthResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  SocialAuthUserResponse `json:"user"`
}