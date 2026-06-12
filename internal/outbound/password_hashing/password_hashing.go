package passwordhashing

import (
	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHash interface {
	HashingPassword(Password auth.Password) (auth.Password, error)
	ValidateEnterPassword(hashingPassword string, password string) bool
}

type bcryptPasswordHash struct {
	cost int
}

func CreateBcryptPasswordHash(cost int) BcryptPasswordHash {
	return &bcryptPasswordHash{
		cost: cost,
	}
}

func (s *bcryptPasswordHash) HashingPassword(Password auth.Password) (auth.Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(Password.Value), s.cost)
	if err != nil {
		return auth.Password{}, err
	}
	hashPassword, error := auth.NewHashedPassword(string(hash))
	if error != nil {
		return auth.Password{}, err
	}
	return hashPassword, nil
}

func (s *bcryptPasswordHash) ValidateEnterPassword(hashingPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashingPassword), []byte(password))
	return err == nil
}
