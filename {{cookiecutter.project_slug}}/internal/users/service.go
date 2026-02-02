package users

import (
	"errors"

	"{{ cookiecutter.module_path }}/internal/users/domain"
)

// Service provides application-level logic for posts.
type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateUser(u *domain.User) (*domain.User, error) {
	if u.Age < 13 {
		return nil, errors.New("user must be at least 13 years old")
	}
	if u.Email == "" {
		return nil, errors.New("email is required")
	}
	if u.Username == "" {
		return nil, errors.New("username is required")
	}
	return s.repo.Create(u)
}

func (s *Service) GetUserByID(id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("user ID is required")
	}
	return s.repo.GetByID(id)
}

func (s *Service) UpdateUser(u *domain.User) (*domain.User, error) {
	return s.repo.Update(u)
}

func (s *Service) DeleteUser(id string) error {
	if id == "" {
		return  errors.New("user ID is required")
	}
	return s.repo.Delete(id)
}

func (s *Service) ListUsers(limit, offset int32) ([]*domain.User, error) {
	return s.repo.List(limit, offset)
}
