package auth

import (
	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"github.com/google/uuid"
)

type CreateUser struct {
	Uuid     uuid.UUID
	Name     auth.Name
	Email    auth.EmailAddress
	Password auth.Password
}

type EnterCodeRequest auth.VerificationCode

func NewEnterCodeRequest(code auth.VerificationCode) EnterCodeRequest {
	return EnterCodeRequest(code)
}

func NewCreateUser(input CreateUser, password auth.Password) CreateUser {
	return CreateUser{
		Uuid:     input.Uuid,
		Name:     input.Name,
		Email:    input.Email,
		Password: password,
	}
}

type LoginUser struct {
	Email    auth.EmailAddress
	Password auth.Password
}
