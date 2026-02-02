package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"{{ cookiecutter.module_path }}/internal/users/domain"
	db "{{ cookiecutter.module_path }}/internal/users/models"
	u "{{ cookiecutter.module_path }}/utils"
)

type PgRepo struct {
	q *db.Queries
}

// NewPgRepo creates a new instance of PgRepo.
func NewPgRepo(conn *pgxpool.Pool) *PgRepo {
	return &PgRepo{q: db.New(conn)}
}

func toDomain(u db.User) *domain.User {
	return &domain.User{
		ID:                u.ID,
		Email:             u.Email,
		Username:          u.Username,
		PhoneNumber:       u.PhoneNumber,
		Avatar:            u.Avatar,
		Age:               u.Age,
		Gender:            u.Gender,
		IsActive:          u.IsActive,
		PasswordChangedAt: u.PasswordChangedAt,
		UpdatedAt:         u.UpdatedAt,
		CreatedAt:         u.CreatedAt,
	}
}

func fromDomain(u *domain.User) db.CreateUserParams {
	return db.CreateUserParams{
		Email:          u.Email,
		Username:       u.Username,
		PhoneNumber:    u.PhoneNumber,
		Avatar:         u.Avatar,
		Age:            u.Age,
		Gender:         u.Gender,
		HashedPassword: u.HashedPassword,
	}
}

func (r *PgRepo) Create(u *domain.User) (*domain.User, error) {
	ctx := context.Background()
	params := fromDomain(u)
	createdUser, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return toDomain(createdUser), nil
}

func (r *PgRepo) GetByID(id string) (*domain.User, error) {
	ctx := context.Background()
	// convert string ID to pgtype.UUID
	newID, err := u.ConvertToPgUUIDFromString(id)
	if err != nil {
		return nil, err
	}

	user, err := r.q.GetUserByID(ctx, newID)
	if err != nil {
		return nil, err
	}
	return toDomain(user), nil
}

func (r *PgRepo) Update(u *domain.User) (*domain.User, error) {
	ctx := context.Background()
	params := db.UpdateUserParams{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		PhoneNumber: u.PhoneNumber,
		Avatar:      u.Avatar,
		Age:         u.Age,
	}
	updatedUser, err := r.q.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return toDomain(updatedUser), nil
}

func (r *PgRepo) Delete(id string) error {
	ctx := context.Background()
	// convert string ID to pgtype.UUID
	newID, err := u.ConvertToPgUUIDFromString(id)
	if err != nil {
		return err
	}

	return r.q.DeleteUser(ctx, newID)
}

func (r *PgRepo) List(limit, offset int32) ([]*domain.User, error) {
	ctx := context.Background()
	users, err := r.q.ListUsers(ctx, db.ListUsersParams{
		Limit:  limit,
		Offset: (offset - 1) * limit,
	})
	if err != nil {
		return nil, err
	}

	var domainUsers []*domain.User
	for _, u := range users {
		domainUsers = append(domainUsers, toDomain(u))
	}
	return domainUsers, nil
}
