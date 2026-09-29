package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michaelassa01/blog/internal/users/domain"
	"github.com/michaelassa01/blog/utils"

	models "github.com/michaelassa01/blog/internal/users/models"
)

type UserRepo struct {
	q  *models.Queries
	db *pgxpool.Pool
}

// NewUserRepo creates a new instance of UserRepo.
func NewUserRepo(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		q:  models.New(conn),
		db: conn,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *domain.User) ( *domain.User , error) {
	params := models.CreateUserParams{
		Email:          user.Email,
		PhoneNumber:    user.PhoneNumber,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		FullName:       user.FullName,
		AccountType:    user.AccountType,
		HashedPassword: user.HashedPassword,
		EmailVerified:  user.EmailVerified,
	}

	dbUser, err := r.q.CreateUser(ctx, params)
	if err != nil {
		return nil,err
	}

	r.mapDBUserToDomain(dbUser, user)
	return nil,nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	dbUser, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.mapDBUserToDomainPtr(dbUser), nil
}

func (r *UserRepo) GetUserByPhoneNumber(ctx context.Context, phone string) (*domain.User, error) {
	dbUser, err := r.q.GetUserByPhoneNumber(ctx, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.mapDBUserToDomainPtr(dbUser), nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	nwId := utils.ConvertToPgUUID(id)
	dbUser, err := r.q.GetUserByID(ctx, nwId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return r.mapDBUserToDomainPtr(dbUser), nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user *domain.User) error {
	params := models.UpdateUserParams{
		ID:          utils.ConvertToPgUUID(user.ID),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}

	_, err := r.q.UpdateUser(ctx, params)
	return err
}

func (r *UserRepo) VerifyUserEmail(ctx context.Context, email string) error {
	return r.q.VerifyUserEmail(ctx, email)
}

func (r *UserRepo) UpdateUserPassword(ctx context.Context, email, hashedPassword string) error {
	params := models.UpdateUserPasswordParams{
		Email:          email,
		HashedPassword: hashedPassword,
	}
	return r.q.UpdateUserPassword(ctx, params)
}

func (r *UserRepo) Delete(id string) error {
	ctx := context.Background()
	// convert string ID to pgtype.UUID
	newID, err := utils.ConvertToPgUUIDFromString(id)
	if err != nil {
		return err
	}

	return r.q.DeleteUser(ctx, newID)
}


// Mapping functions
func (r *UserRepo) mapDBUserToDomain(dbUser models.User, user *domain.User) {
	dbUserId, _ := utils.ConvertFromPgUUID(dbUser.ID)
	user.ID = dbUserId
	user.Email = dbUser.Email
	user.PhoneNumber = dbUser.PhoneNumber
	user.FirstName = dbUser.FirstName
	user.LastName = dbUser.LastName
	user.FullName = dbUser.FullName
	user.AccountType = dbUser.AccountType
	user.HashedPassword = dbUser.HashedPassword
	user.EmailVerified = dbUser.EmailVerified
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
}

func (r *UserRepo) mapDBUserToDomainPtr(dbUser models.User) *domain.User {
	dbUserId, _ := utils.ConvertFromPgUUID(dbUser.ID)
	return &domain.User{
		ID:             dbUserId,
		Email:          dbUser.Email,
		PhoneNumber:    dbUser.PhoneNumber,
		FirstName:      dbUser.FirstName,
		LastName:       dbUser.LastName,
		FullName:       dbUser.FullName,
		AccountType:    dbUser.AccountType,
		HashedPassword: dbUser.HashedPassword,
		EmailVerified:  dbUser.EmailVerified,
		CreatedAt:      dbUser.CreatedAt,
		UpdatedAt:      dbUser.UpdatedAt,
	}
}
