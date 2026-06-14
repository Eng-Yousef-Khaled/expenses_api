package auth

import (
	"strings"
	"time"
)

const EXPIRES_MINUTES = 30

type Name string

func NewName(raw string) (Name, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", NameError("name cannot be empty")
	}
	return Name(trimmed), nil
}

type EmailAddress string

func validateEmailAddress(email string) error {
	if email == "" {
		return &EmailAddressError{InvalidEmailAddress: email}
	}

	// TODO: Unimplemented example (e.g., regex check)

	return nil
}

func NewEmailAddress(raw string) (EmailAddress, error) {
	trimmed := strings.TrimSpace(raw)
	if err := validateEmailAddress(trimmed); err != nil {
		return "", err
	}
	if !strings.Contains(raw, "@") {
		return "", InvalidEmailAddress
	}
	return EmailAddress(raw), nil
}

type Password struct {
	Value  string
	Hashed bool
}
type RawPassword string
type HashedPassword string

func NewRawPassword(raw string) (Password, error) {
	trimmed := strings.TrimSpace(raw)

	if trimmed == "" {
		return Password{}, PasswordError
	}

	return Password{
		Value:  trimmed,
		Hashed: false,
	}, nil
}

func NewHashedPassword(hash string) (Password, error) {
	trimmed := strings.TrimSpace(hash)
	if trimmed == "" {
		return Password{}, PasswordError
	}
	return Password{
		Value:  trimmed,
		Hashed: true,
	}, nil
}

type VerificationCode string

func NewVerificationCode(raw string) (VerificationCode, error) {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) != 6 {
		return "", VerificationCodeError
	}
	return VerificationCode(raw), nil
}

type VerificationCodeExpire time.Time

func NewVerificationCodeExpired(raw time.Time) (VerificationCodeExpire, error) {
	if raw.Before(time.Now()) {
		return VerificationCodeExpire(time.Now()), VerificationCodeExpireError
	}
	return VerificationCodeExpire(raw), nil
}

type UserId int64

func NewUserId(raw int64) (UserId, error) {
	if raw == 0 {
		return 0, InvalidUserId
	}
	return UserId(raw), nil
}
