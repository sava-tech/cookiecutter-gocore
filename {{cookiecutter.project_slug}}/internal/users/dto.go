package users

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"{{ cookiecutter.module_path }}/internal/users/domain"
)

type CreateUserReq struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Avatar      string `json:"avatar" binding:"required"`
	Password    string `json:"password" binding:"required,min=5"`
	Age         int32  `json:"age" binding:"required,gte=13"`
}

type CreateUserRes struct {
	ID          pgtype.UUID `json:"id"`
	Email       string      `json:"email"`
	Username    string      `json:"username"`
	PhoneNumber string      `json:"phone_number"`
	Avatar      string      `json:"avatar"`
	Age         int32       `json:"age"`
	UpdatedAt   time.Time   `json:"updated_at"`
	CreatedAt   time.Time   `json:"created_at"`
}

type GetUserReq struct {
	ID string `json:"id" binding:"required,uuid"`
}

type GetUserRes struct {
	ID          pgtype.UUID `json:"id"`
	Email       string      `json:"email"`
	Username    string      `json:"username"`
	PhoneNumber string      `json:"phone_number"`
	Avatar      string      `json:"avatar"`
	Age         int32       `json:"age"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type UpdateUserReq struct {
	Email       string `json:"email" binding:"omitempty,email"`
	Username    string `json:"username" binding:"omitempty"`
	PhoneNumber string `json:"phone_number" binding:"omitempty"`
	Avatar      string `json:"avatar" binding:"omitempty"`
	Age         int32  `json:"age" binding:"omitempty,gte=13"`
}

type UpdateUserRes struct {
	ID          pgtype.UUID `json:"id"`
	Email       string      `json:"email"`
	Username    string      `json:"username"`
	PhoneNumber string      `json:"phone_number"`
	Avatar      string      `json:"avatar"`
	Age         int32       `json:"age"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ListUsersReq struct {
	Limit  int32 `json:"limit" binding:"omitempty,min=1,max=100"`
	Offset int32 `json:"offset" binding:"omitempty,min=0"`
}

type ListUsersRes struct {
	Users []GetUserRes `json:"users"`
}

type DeleteUserReq struct {
	ID string `json:"id" binding:"required,uuid"`
}

func toUserResponse(u *domain.User) *CreateUserRes {
	return &CreateUserRes{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		PhoneNumber: u.PhoneNumber,
		Avatar:      u.Avatar,
		Age:         u.Age,
		UpdatedAt:   u.UpdatedAt,
		CreatedAt:   u.CreatedAt,
	}
}
