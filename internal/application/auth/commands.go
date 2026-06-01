package auth

import (
	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"github.com/google/uuid"
)

type CreateUser struct {
	Uuid     uuid.UUID
	Name     auth.Name
	Email    auth.EmailAddress
	Password auth.RawPassword
}
type LoginUser struct {
	Email    auth.EmailAddress
	Password auth.RawPassword
}
