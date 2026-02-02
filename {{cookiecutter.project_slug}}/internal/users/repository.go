package users

import "{{ cookiecutter.module_path }}/internal/users/domain"


type Repository interface {
	Create(u *domain.User) (*domain.User, error)
	GetByID(id string) (*domain.User, error)
	Update(u *domain.User) (*domain.User, error)
	Delete(id string) error
	List(limit, offset int32) ([]*domain.User, error)
}
