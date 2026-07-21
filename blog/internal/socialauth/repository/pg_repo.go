package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michaelassa01/blog/internal/socialauth/domain"
	models "github.com/michaelassa01/blog/internal/socialauth/models"
	"github.com/michaelassa01/blog/utils"
)

type PgRepo struct {
	q  *models.Queries
	db *pgxpool.Pool
}

func NewPgRepo(db *pgxpool.Pool) *PgRepo {
	return &PgRepo{
		q:  models.New(db),
		db: db,
	}
}

func (r *PgRepo) Create(ctx context.Context, socialUser *domain.SocialAuthUser) error {
	params := models.CreateSocialUserParams{
		ID:          utils.ConvertToPgUUID(socialUser.ID),
		UserID:      utils.ConvertToPgUUID(socialUser.UserID),
		Provider:    socialUser.Provider,
		ProviderID:  socialUser.ProviderID,
		Email:       socialUser.Email,
		Name:        &socialUser.Name,
		AvatarUrl:   &socialUser.AvatarURL,
		AccessToken: &socialUser.AccessToken,
	}

	_, err := r.q.CreateSocialUser(ctx, params)
	return err
}

func (r *PgRepo) GetByProviderID(ctx context.Context, provider, providerID string) (*domain.SocialAuthUser, error) {
	dbSocialUser, err := r.q.GetSocialUserByProviderID(ctx, models.GetSocialUserByProviderIDParams{
		Provider:   provider,
		ProviderID: providerID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSocialUserNotFound
		}
		return nil, err
	}

	return r.mapToDomain(dbSocialUser), nil
}

func (r *PgRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.SocialAuthUser, error) {
	nwUserID := utils.ConvertToPgUUID(userID)
	dbSocialUser, err := r.q.GetSocialUserByUserID(ctx, nwUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSocialUserNotFound
		}
		return nil, err
	}

	return r.mapToDomain(dbSocialUser), nil
}

func (r *PgRepo) Update(ctx context.Context, socialUser *domain.SocialAuthUser) error {
	params := models.UpdateSocialUserParams{
		ID:          utils.ConvertToPgUUID(socialUser.ID),
		AccessToken: &socialUser.AccessToken,
		AvatarUrl:   &socialUser.AvatarURL,
	}

	_, err := r.q.UpdateSocialUser(ctx, params)
	return err
}

// Mapping helper
func (r *PgRepo) mapToDomain(dbSocialUser models.SocialAuthUser) *domain.SocialAuthUser {
	newsocialID, _ := utils.ConvertFromPgUUID(dbSocialUser.ID)
	newsocialUserID, _ := utils.ConvertFromPgUUID(dbSocialUser.UserID)
	
	return &domain.SocialAuthUser{
		ID:          newsocialID,
		UserID:      newsocialUserID,
		Provider:    dbSocialUser.Provider,
		ProviderID:  dbSocialUser.ProviderID,
		Email:       dbSocialUser.Email,
		Name:        *dbSocialUser.Name,
		AvatarURL:   *dbSocialUser.AvatarUrl,
		AccessToken: *dbSocialUser.AccessToken,
		CreatedAt:   dbSocialUser.CreatedAt,
		UpdatedAt:   dbSocialUser.UpdatedAt,
	}
}
