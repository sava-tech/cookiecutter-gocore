// @title {{ cookiecutter.project_name }} API
// @version 1.0
// @description This is the API documentation for {{ cookiecutter.project_name }}.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@{{ cookiecutter.project_name }}.com

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
	_ "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/docs"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/server"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/utils"
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
