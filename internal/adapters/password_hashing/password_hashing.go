package passwordhashing

import (
	"github.com/eng-yousef-khaled/expenses_api/internal/domain/auth"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHash interface {
	HashingPassword(Password *auth.Password) error
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

func (s *bcryptPasswordHash) HashingPassword(password *auth.Password) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*password = auth.Password(hash)
	return nil
}

func (s *bcryptPasswordHash) ValidateEnterPassword(hashingPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashingPassword), []byte(password))
	return err == nil
}
