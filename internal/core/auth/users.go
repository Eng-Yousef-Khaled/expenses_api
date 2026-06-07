package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       int64
	Uuid     uuid.UUID
	Name     Name
	Email    EmailAddress
	Password HashedPassword
}

type UserVerificationCode struct {
	ID               int64
	UserID           int64
	VerificationCode VerificationCode
	ExpiresAt        VerificationCodeExpire
	CreatedAt        time.Time
}
