package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             int64
	Uuid           uuid.UUID
	Name           Name
	Email          EmailAddress
	Password       Password
	IsVerification bool
}

type UserVerificationCode struct {
	ID               int64
	UserID           UserId
	VerificationCode VerificationCode
	ExpiresAt        VerificationCodeExpire
	CreatedAt        time.Time
}

type UserSession struct {
	Id           int64
	UserId       UserId
	RefreshToken RefreshToken
	AccessToken  AccessToken
	CreatedAt    time.Time
	ExpiredAt    time.Time
	IsActive     bool
}
