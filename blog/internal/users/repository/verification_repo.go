package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michaelassa01/blog/internal/users/domain"
	models "github.com/michaelassa01/blog/internal/users/models"
	"github.com/michaelassa01/blog/utils"
)

type VerificationRepo struct {
	q  *models.Queries
	db *pgxpool.Pool
}

// NewVerificationRepo creates a new instance of VerificationRepo.
func NewVerificationRepo(conn *pgxpool.Pool) *VerificationRepo {
	return &VerificationRepo{
		q:  models.New(conn),
		db: conn,
	}
}

func (r *VerificationRepo) CreateVerification(ctx context.Context, verification *domain.Verification) error {
	params := models.CreateVerificationParams{
		Code:             verification.Code,
		IdentifierType:   verification.IdentifierType,
		Identifier:       verification.Identifier,
		VerificationType: verification.VerificationType,
		ExpiredAt:        verification.ExpiredAt,
	}

	dbVerification, err := r.q.CreateVerification(ctx, params)
	if err != nil {
		fmt.Printf("Failed to create verification record: %v\n", err)
		return err
	}

	r.mapDBVerificationToDomain(dbVerification, verification)
	return nil
}

func (r *VerificationRepo) GetValidVerificationCode(ctx context.Context, identifier, code, verificationType string) (*domain.Verification, error) {
	params := models.GetValidVerificationCodeParams{
		Identifier:       identifier,
		Code:             code,
		VerificationType: verificationType,
	}

	dbVerification, err := r.q.GetValidVerificationCode(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrVerificationNotFound
		}
		return nil, err
	}

	return r.mapDBVerificationToDomainPtr(dbVerification), nil
}

func (r *VerificationRepo) GetVerificationHash(ctx context.Context, identifier, verificationType string) (string, error) {
	params := models.GetVerificationHashParams{
		Identifier:       identifier,
		VerificationType: verificationType,
	}
	return r.q.GetVerificationHash(ctx, params)
}

func (r *VerificationRepo) MarkVerificationAsUsed(ctx context.Context, id uuid.UUID) error {
	nwId := utils.ConvertToPgUUID(id)
	return r.q.MarkVerificationAsUsed(ctx, nwId)
}

func (r *VerificationRepo) InvalidateVerificationCodes(ctx context.Context, identifier, verificationType string) error {
	params := models.InvalidateVerificationCodesParams{
		Identifier:       identifier,
		VerificationType: verificationType,
	}
	return r.q.InvalidateVerificationCodes(ctx, params)
}

// Mapping functions to convert between database models and domain models

func (r *VerificationRepo) mapDBVerificationToDomain(dbVerification models.Verification, verification *domain.Verification) {
	dbVerificationId, _ := utils.ConvertFromPgUUID(dbVerification.ID)
	verification.ID = dbVerificationId
	verification.Code = dbVerification.Code
	verification.IdentifierType = dbVerification.IdentifierType
	verification.Identifier = dbVerification.Identifier
	verification.VerificationType = dbVerification.VerificationType
	verification.ExpiredAt = dbVerification.ExpiredAt
	verification.Used = dbVerification.Used
	verification.CreatedAt = dbVerification.CreatedAt
}

func (r *VerificationRepo) mapDBVerificationToDomainPtr(dbVerification models.Verification) *domain.Verification {
	dbVerificationId, _ := utils.ConvertFromPgUUID(dbVerification.ID)
	return &domain.Verification{
		ID:               dbVerificationId,
		Code:             dbVerification.Code,
		IdentifierType:   dbVerification.IdentifierType,
		Identifier:       dbVerification.Identifier,
		VerificationType: dbVerification.VerificationType,
		ExpiredAt:        dbVerification.ExpiredAt,
		Used:             dbVerification.Used,
		CreatedAt:        dbVerification.CreatedAt,
	}
}
