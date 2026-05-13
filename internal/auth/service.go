package users

import (
	"context"
	"fmt"
	"log"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/eng-yousef-khaled/expenses_api/internal/adapters/queue"
)

type Service interface {
	CreateUser(ctx context.Context, u repo.User) (repo.User, error)
}
type svc struct {
	repo  repo.Querier
	Queue queue.TaskQueue
}

func CreateService(repo repo.Querier, Queue queue.TaskQueue) Service {
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
	return s.repo.CreateUser(ctx, params)
}
