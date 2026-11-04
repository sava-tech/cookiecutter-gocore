package main

import (
	"log"

	"{{ cookiecutter.module_path }}/internal/reports"
	"{{ cookiecutter.module_path }}/pkg/token"
	"{{ cookiecutter.module_path }}/utils"
)

func main() {
	config, err := utils.LoadConfig("")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	if config.UsersServiceURL == "" || config.ApiAccessKey == "" {
		log.Fatal("USERS_SERVICE_URL and API_ACCESS_KEY must be set")
	}

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	usersClient := reports.NewUsersClient(config.UsersServiceURL, config.ApiAccessKey)
	router := reports.NewRouter(tokenMaker, usersClient)
	if err := router.Run(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
