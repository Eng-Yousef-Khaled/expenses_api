package auth

import (
	"errors"
	"fmt"
)

type CreateUserError struct {
	Duplicate       string
	PasswordHashing string
	Unknown         string
}

var (
	InvalidEmailOrPasswordError = errors.New("can't Login, Invalid email or password")
	LoginPasswordHashingError   = errors.New("Failed To hashing password")
)

var (
	VerificationCodeError        = errors.New("code must be a 6-digit number")
	VerificationCodeExpiredError = errors.New("Verification code is expired")
	VerificationCodeExpireError  = errors.New("Verification code is before current time")
	VerificationCodeInvalidError = errors.New("Verification code invalid")
)

type NameError string

func (e NameError) Error() string {
	return string(e)
}

type EmailAddressError struct {
	InvalidEmailAddress string
}

func (e *EmailAddressError) Error() string {
	return fmt.Sprintf("%q is not a valid email address", e.InvalidEmailAddress)
}

type PasswordError string

func (e PasswordError) Error() string {
	return string(e)
}
