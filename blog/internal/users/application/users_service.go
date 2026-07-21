package application

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/michaelassa01/blog/internal/users/domain"
	"github.com/michaelassa01/blog/internal/users/dto"
	"github.com/michaelassa01/blog/pkg/emailer"
	"github.com/michaelassa01/blog/pkg/token"
)

type UserService struct {
	repo            domain.Repository
	tokenMaker      token.Maker
	passwordService domain.PasswordService
	emailService    emailer.Mailer
	config          *Config
	validator       *validator.Validate
}

type Config struct {
	AccessTokenDuration      time.Duration
	RefreshTokenDuration     time.Duration
	VerificationCodeDuration time.Duration
}

func NewUserService(
	repo domain.Repository,
	tokenMaker token.Maker,
	passwordService domain.PasswordService,
	emailService emailer.Mailer,
	config *Config,
) *UserService {
	return &UserService{
		repo:            repo,
		tokenMaker:      tokenMaker,
		passwordService: passwordService,
		emailService:    emailService,
		config:          config,
		validator:       validator.New(),
	}
}

func (s *UserService) GetAccount(ctx context.Context, req *dto.GetUserReq) (*dto.RegisterResponse, error) {
	user, err := s.repo.GetUserByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		Email:         user.Email,
		PhoneNumber:   user.PhoneNumber,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		FullName:      user.FullName,
		AccountType:   user.AccountType,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
	}, nil
}
