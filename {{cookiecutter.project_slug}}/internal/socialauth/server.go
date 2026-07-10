package socialauth

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/socialauth/application"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/socialauth/infrastructure"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/socialauth/repository"
	userApplication "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/application"
	userRepository "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/repository"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/emailer"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/password"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/token"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/utils"
)

type Services struct {
	User       *userApplication.UserService
	Auth       *userApplication.AuthService
	SocialAuth *application.SocialAuthService
}

func NewServices(
	dbConn *pgxpool.Pool,
	tokenMaker token.Maker,
	config utils.Config,
) *Services {
	sessionSecret := config.SessionSecret
	if sessionSecret == "" {
		sessionSecret = "12345678901234567890123456789012"
		log.Printf("SESSION_SECRET not set, using fallback - NOT SECURE FOR PRODUCTION")
	}
	infrastructure.SetupSession(sessionSecret)
	infrastructure.SetupGoth(infrastructure.GothProviderConfig{})

	userRepo := userRepository.NewPgRepo(dbConn)
	authRepo := userRepository.NewPgRepo(dbConn)
	socialRepo := repository.NewPgRepo(dbConn)

	// infrastructure services
	passwordService := password.NewPasswordService()
	emailService, err := emailer.NewMailer(config)
	if err != nil {
		fmt.Println("email service not handled or config")
		panic(err)
	}
	authConfig := &application.Config{
		AccessTokenDuration:  config.AccessTokenDuration,
		RefreshTokenDuration: config.RefreshTokenDuration,
		CallbackURL:          config.SocialCallbackURL,
	}

	return &Services{
		User: userApplication.NewUserService(
			userRepo,
			tokenMaker,
			passwordService,
			emailService,
			&userApplication.Config{},
		),
		Auth: userApplication.NewAuthService(
			authRepo,
			tokenMaker,
			passwordService,
			emailService,
			&userApplication.Config{},
		),
		SocialAuth: application.NewSocialAuthService(
			socialRepo,
			userRepo,
			tokenMaker,
			authConfig,
		),
	}
}
