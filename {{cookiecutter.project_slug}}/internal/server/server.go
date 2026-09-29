package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"{{ cookiecutter.module_path }}/internal/socialauth"
	"{{ cookiecutter.module_path }}/internal/users"
	"{{ cookiecutter.module_path }}/pkg/token"
	"{{ cookiecutter.module_path }}/utils"
)

// Server serves HTTP request for {{ cookiecutter.project_name }} service
type Server struct {
	config             utils.Config
	db                 *pgxpool.Pool
	tokenMaker         token.Maker
	router             *gin.Engine
	UserServices       *users.Services
	SocialAuthServices *socialauth.Services
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config utils.Config, dbConn *pgxpool.Pool) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	// imported services
	userServices := users.NewServices(dbConn, tokenMaker, config)
	socialauthServices := socialauth.NewServices(dbConn, tokenMaker, config)

	s := &Server{
		config:       config,
		db:           dbConn,
		tokenMaker:   tokenMaker,
		UserServices: userServices,
		SocialAuthServices: socialauthServices,
	}

	s.setupRouter(config)

	return s, nil
}

// Start runs the HTTP server on a specification address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
