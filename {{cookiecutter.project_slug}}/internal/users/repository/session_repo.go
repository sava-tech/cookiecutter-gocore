package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/domain"
	models "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/models"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/utils"
)

type SessionRepo struct {
	q  *models.Queries
	db *pgxpool.Pool
}

// NewSessionRepo creates a new instance of SessionRepo.
func NewSessionRepo(conn *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{
		q:  models.New(conn),
		db: conn,
	}
}


func (r *SessionRepo) CreateSession(ctx context.Context, session *domain.Session) error {
	// convert session id to pgtype.UUID
	nwSessionId := utils.ConvertToPgUUID(session.ID)
	nwUserId := utils.ConvertToPgUUID(session.UserID)
	params := models.CreateSessionParams{
		ID:           nwSessionId,
		UserID:       nwUserId,
		Email:        session.Email,
		RefreshToken: session.RefreshToken,
		UserAgent:    session.UserAgent,
		ClientIp:     session.ClientIP,
		IsBlocked:    session.IsBlocked,
		ExpiresAt:    session.ExpiresAt,
	}

	dbSession, err := r.q.CreateSession(ctx, params)
	if err != nil {
		return err
	}

	r.mapDBSessionToDomain(dbSession, session)
	return nil
}

func (r *SessionRepo) GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	nwId := utils.ConvertToPgUUID(id)
	dbSession, err := r.q.GetSessionByID(ctx, nwId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}

	return r.mapDBSessionToDomainPtr(dbSession), nil
}

func (r *SessionRepo) UpdateSession(ctx context.Context, session *domain.Session) error {
	nwSessionId := utils.ConvertToPgUUID(session.ID)
	params := models.UpdateSessionParams{
		ID:           nwSessionId,
		IsBlocked:    session.IsBlocked,
		RefreshToken: session.RefreshToken,
	}

	_, err := r.q.UpdateSession(ctx, params)
	return err
}

func (r *SessionRepo) DeleteSession(ctx context.Context, id uuid.UUID) error {
	nwId := utils.ConvertToPgUUID(id)
	return r.q.DeleteSession(ctx, nwId)
}


// Mapping functions to convert between database models and domain models
func (r *SessionRepo) mapDBSessionToDomain(dbSession models.Session, session *domain.Session) {
	dbSessionId, _ := utils.ConvertFromPgUUID(dbSession.ID)
	nwUserId, _ := utils.ConvertFromPgUUID(dbSession.UserID)

	session.ID = dbSessionId
	session.UserID = nwUserId
	session.Email = dbSession.Email
	session.RefreshToken = dbSession.RefreshToken
	session.UserAgent = dbSession.UserAgent
	session.ClientIP = dbSession.ClientIp
	session.IsBlocked = dbSession.IsBlocked
	session.ExpiresAt = dbSession.ExpiresAt
	session.CreatedAt = dbSession.CreatedAt
}

func (r *SessionRepo) mapDBSessionToDomainPtr(dbSession models.Session) *domain.Session {
	dbSessionId, _ := utils.ConvertFromPgUUID(dbSession.ID)
	nwUserId, _ := utils.ConvertFromPgUUID(dbSession.UserID)
	return &domain.Session{
		ID:           dbSessionId,
		UserID:       nwUserId,
		Email:        dbSession.Email,
		RefreshToken: dbSession.RefreshToken,
		UserAgent:    dbSession.UserAgent,
		ClientIP:     dbSession.ClientIp,
		IsBlocked:    dbSession.IsBlocked,
		ExpiresAt:    dbSession.ExpiresAt,
		CreatedAt:    dbSession.CreatedAt,
	}
}
