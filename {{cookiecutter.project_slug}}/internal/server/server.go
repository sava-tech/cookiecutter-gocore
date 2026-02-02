package server

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	// wsh "{{ cookiecutter.module_path }}/server/ws"

	"github.com/gin-gonic/gin"
	"{{ cookiecutter.module_path }}/pkg/token"
	"{{ cookiecutter.module_path }}/utils"
)

// Server serves HTTP request for {{ cookiecutter.project_slug}} service
type Server struct {
	config     utils.Config
	db         *pgxpool.Pool
	tokenMaker token.Maker
	router     *gin.Engine
	Services   *Services
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config utils.Config, dbConn *pgxpool.Pool) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	services := NewServices(dbConn)

	s := &Server{
		config:     config,
		db:         dbConn,
		tokenMaker: tokenMaker,
		Services:   services,
	}

	s.setupRouter(config)

	return s, nil
}

// Start runs the HTTP server on a specification address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// func (ws *Server) WSHandler() *wsh.WSHandler {
// 	return ws.wsHandler
// }

func errorResponse(err error) gin.H {
	// // Split the error message into individual validation errors
	// lines := strings.Split(err.Error(), "\n")

	// // Prepare a slice of maps to hold individual field errors
	// var details []map[string]string

	// for _, line := range lines {
	// 	// Example line: "Key: 'tansferRequest.FromAccountID' Error:Field validation for 'FromAccountID' failed on the 'required' tag"
	// 	parts := strings.SplitN(line, "Error:", 2)
	// 	if len(parts) == 2 {
	// 		fieldPart := strings.TrimSpace(parts[0])
	// 		errorMsg := strings.TrimSpace(parts[1])

	// 		details = append(details, map[string]string{
	// 			"key":   fieldPart,
	// 			"error": errorMsg,
	// 		})
	// 	}
	// }

	// return gin.H{
	// 	"errors": details,
	// }
	return gin.H{"error": err.Error()}

}
