package auth

import "strings"

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
	return EmailAddress(raw), nil
}

type RawPassword string
type HashedPassword string

func NewRawPassword(raw string) (RawPassword, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", PasswordError("password cannot be empty")
	}
	return RawPassword(trimmed), nil
}

func NewHashedPassword(raw RawPassword) (HashedPassword, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return "", PasswordError("password cannot be empty")
	}
	return HashedPassword(trimmed), nil
}
