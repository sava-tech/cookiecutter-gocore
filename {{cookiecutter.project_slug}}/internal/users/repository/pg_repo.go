package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/domain"
)

type PgRepo struct {
	domain.UserRepository
	domain.SessionRepository
	domain.VerificationRepository
}

func NewPgRepo(db *pgxpool.Pool) *PgRepo {
	return &PgRepo{
		UserRepository:         NewUserRepo(db),
		SessionRepository:      NewSessionRepo(db),
		VerificationRepository: NewVerificationRepo(db),
	}
}
