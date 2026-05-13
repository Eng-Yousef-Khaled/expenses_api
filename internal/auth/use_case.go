package users

import (
	"golang.org/x/crypto/bcrypt"
)

func HashingPassword(password *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*password = string(hash)
	return nil
}

func ValidateEnterPassword(hashingPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashingPassword), []byte(password))
	return err == nil
}
