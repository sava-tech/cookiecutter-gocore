// @title cnn-nigeria API
// @version 1.0
// @description This is the API documentation for blog.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@blog.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name api-access-key
// @description Use your API key to authenticate requests.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/michaelassa01/blog/docs"
	"github.com/michaelassa01/blog/internal/server"
	"github.com/michaelassa01/blog/utils"
)

var (
	conn *pgxpool.Pool
)

func main() {

	// Then start consuming
	config, err := utils.LoadConfig("")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}

	serverConfig, err := server.NewServer(config, conn)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	// Setup cron jobs
	// serverConfig.SetupCronJobs(store)
	err = serverConfig.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	defer conn.Close()
}
