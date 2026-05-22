package auth

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type CreateUserError struct {
	Duplicate       string
	PasswordHashing string
	Unknown         string
}
type CreateUser struct {
	Uuid     uuid.UUID
	Name     Name
	Email    EmailAddress
	Password Password
}
type User struct {
	ID       int64
	Uuid     uuid.UUID
	Name     string
	Email    string
	Password string
}
type EmailMessage struct {
	To      []string
	Subject string
	Body    *string
	Data    interface{}
}

type VerificationCodeMail struct {
	Email   string
	Message string
}
type NameError string

func (e NameError) Error() string {
	return string(e)
}

type Name string

func NewName(raw string) (Name, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", NameError("name cannot be empty")
	}
	return Name(trimmed), nil
}

type EmailAddress string

type EmailAddressError struct {
	InvalidEmailAddress string
}

func (e *EmailAddressError) Error() string {
	return fmt.Sprintf("%q is not a valid email address", e.InvalidEmailAddress)
}

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
	return EmailAddress(raw), nil
}

type PasswordError string

func (e PasswordError) Error() string {
	return string(e)
}

type Password string

func NewPassword(raw string) (Password, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", PasswordError("password cannot be empty")
	}
	return Password(trimmed), nil
}
