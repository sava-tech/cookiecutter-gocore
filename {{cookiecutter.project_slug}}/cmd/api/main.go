// @title {{ cookiecutter.project_name }} API
// @version 1.0
// @description This is the API documentation for {{ cookiecutter.project_name }}.
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name api-access-key
// @description Use your API key to authenticate requests.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use `Bearer <token>` for authentication.

package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	// _ "{{ cookiecutter.module_path }}/docs"
	"{{ cookiecutter.module_path }}/internal/server"
	"{{ cookiecutter.module_path }}/utils"
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
