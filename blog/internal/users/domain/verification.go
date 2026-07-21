package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Verification entity
type Verification struct {
	ID               uuid.UUID
	Code             string
	IdentifierType   string
	Identifier       string
	VerificationType string
	ExpiredAt        time.Time
	Used             bool
	CreatedAt        time.Time
}

type VerificationRepository interface {
	CreateVerification(ctx context.Context, verification *Verification) error
	GetValidVerificationCode(ctx context.Context, identifier, code, verificationType string) (*Verification, error)
	GetVerificationHash(ctx context.Context, identifier, verificationType string) (string, error)
	MarkVerificationAsUsed(ctx context.Context, id uuid.UUID) error
	InvalidateVerificationCodes(ctx context.Context, identifier, verificationType string) error
}
