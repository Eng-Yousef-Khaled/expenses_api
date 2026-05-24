package auth

import "fmt"

type CreateUserError struct {
	Duplicate       string
	PasswordHashing string
	Unknown         string
}

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
