package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/eng-yousef-khaled/expenses_api/internal/adapters/queue"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, u repo.User) (repo.User, error)
}
type svc struct {
	repo  repo.Querier
	Queue queue.TaskQueue
}

func CreateService(repo repo.Querier, Queue queue.TaskQueue) UserService {
	return &svc{
		repo:  repo,
		Queue: Queue,
	}
}

func (s *svc) CreateUser(ctx context.Context, user repo.User) (repo.User, error) {
	params := repo.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Uuid:     user.Uuid,
	}
	err := s.Queue.Enqueue("send_welcome_email", map[string]any{
		"email":   user.Email,
		"message": fmt.Sprintf("Welcome dear: %s, your UUID is: %s", user.Name, user.Uuid),
	})
	if err != nil {
		log.Printf("Has an error while adding to Queue : %s \n", err)
	}
	hashErr := s.HashingPassword(&params.Password)
	if hashErr != nil {
		log.Println(hashErr)

		return repo.User{}, errors.New("Fail to encrypt password")
	}
	return s.repo.CreateUser(ctx, params)
}

func (s *svc) HashingPassword(password *string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	*password = string(hash)
	return nil
}

func (s *svc) ValidateEnterPassword(hashingPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashingPassword), []byte(password))
	return err == nil
}
