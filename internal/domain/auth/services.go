package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/eng-yousef-khaled/expenses_api/internal/adapters/queue"
	"github.com/gocraft/work"
	"golang.org/x/crypto/bcrypt"
)

type svc struct {
	repo  repo.Querier
	queue QueuePort
	mail  OVHTextMailer
}

func CreateService(repo repo.Querier, Queue queue.TaskQueue, mail OVHTextMailer) UserService {
	return &svc{
		repo:  repo,
		queue: Queue,
		mail:  mail,
	}
}

func (s *svc) CreateUser(ctx context.Context, user repo.User) (repo.User, error) {
	params := repo.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Uuid:     user.Uuid,
	}
	err := s.queue.Enqueue("send_welcome_email", map[string]any{
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

func (j *svc) ProccessSendMail(job *work.Job) error {

	email := job.ArgString("email")
	message := job.ArgString("message")
	data := map[string]any{"Name": email, "message": message}
	slog.Info("Processing email task after delay", "to", email)
	err := j.mail.SendMail(Request{
		To:       []string{email},
		Subject:  "Code Verification",
		Body:     &message,
		Data:     data,
		MailType: Text,
	})
	if err != nil {
		slog.Error("ProcessSendMail Failed", "error", err)
	}
	return err
}
