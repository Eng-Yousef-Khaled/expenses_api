package passwordhashing

import (
	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHash interface {
	HashingPassword(Password auth.RawPassword) (auth.HashedPassword, error)
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

func (s *bcryptPasswordHash) HashingPassword(Password auth.RawPassword) (auth.HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(Password), s.cost)
	if err != nil {
		return "", err
	}
	return auth.HashedPassword(hash), nil
}

func (s *bcryptPasswordHash) ValidateEnterPassword(hashingPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashingPassword), []byte(password))
	return err == nil
}
