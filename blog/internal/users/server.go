package users

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/michaelassa01/blog/internal/users/application"
	"github.com/michaelassa01/blog/internal/users/repository"
	"github.com/michaelassa01/blog/pkg/emailer"
	"github.com/michaelassa01/blog/pkg/password"
	"github.com/michaelassa01/blog/pkg/token"
	"github.com/michaelassa01/blog/utils"
)

type Services struct {
	User *application.UserService
	Auth *application.AuthService
}

func NewServices(
	dbConn *pgxpool.Pool,
	tokenMaker token.Maker,
	config utils.Config,
) *Services {

	// repositories
	userRepo := repository.NewPgRepo(dbConn)
	authRepo := repository.NewPgRepo(dbConn)

	// infrastructure services
	passwordService := password.NewPasswordService()
	emailService, err := emailer.NewMailer(config)
	if err != nil {
		fmt.Println("email service not handled or config")
		panic(err)
	}

	// auth config
	authConfig := &application.Config{
		AccessTokenDuration:      config.AccessTokenDuration,
		RefreshTokenDuration:     config.RefreshTokenDuration,
		VerificationCodeDuration: config.VerificationCodeDuration,
	}

	return &Services{
		User: application.NewUserService(
			userRepo,
			tokenMaker,
			passwordService,
			emailService,
			authConfig,
		),
		Auth: application.NewAuthService(
			authRepo,
			tokenMaker,
			passwordService,
			emailService,
			authConfig,
		),
	}
}
