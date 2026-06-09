package auth

import (
	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Uuid     uuid.UUID
	Name     auth.Name
	Email    auth.EmailAddress
	Password auth.RawPassword
}

type CreateUser struct {
	Uuid     uuid.UUID
	Name     auth.Name
	Email    auth.EmailAddress
	Password auth.HashedPassword
}

func NewCreateUser(input CreateUserRequest, password auth.HashedPassword) CreateUser {
	return CreateUser{
		Uuid:     input.Uuid,
		Name:     input.Name,
		Email:    input.Email,
		Password: password,
	}
}

type LoginUser struct {
	Email    auth.EmailAddress
	Password auth.RawPassword
}
