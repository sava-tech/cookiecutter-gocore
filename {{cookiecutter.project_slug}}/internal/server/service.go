package server

import (
	"github.com/jackc/pgx/v5/pgxpool"
	users "{{ cookiecutter.module_path }}/internal/users"
)

type Services struct {
	User *users.Service
}

func NewServices(dbConn *pgxpool.Pool) *Services {
	return &Services{
		User: users.NewService(users.NewPgRepo(dbConn)),
	}
}
