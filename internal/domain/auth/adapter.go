package auth

import (
	"context"
	"log/slog"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/gocraft/work"
)

func toDomainUser(raw repo.User) User {
	return User{
		ID:       raw.ID,
		Uuid:     raw.Uuid,
		Name:     raw.Name,
		Email:    raw.Email,
		Password: raw.Password,
	}
}
func toSqlcParams(raw CreateUserRequest) repo.CreateUserParams {
	return repo.CreateUserParams{
		Uuid:     raw.Uuid,
		Name:     raw.Name,
		Email:    raw.Email,
		Password: raw.Password,
	}
}

type PostgresUserRepository interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (User, error)
}

func CreateRepo(q *repo.Queries) PostgresUserRepository {
	return &postgresUserRepository{
		q: q,
	}
}

type postgresUserRepository struct {
	q *repo.Queries
}

func (p *postgresUserRepository) CreateUser(ctx context.Context, user CreateUserRequest) (User, error) {
	user1 := toSqlcParams(user)
	createdUser, err := p.q.CreateUser(ctx, user1)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Error while creating user in adapter", "error", err)
	}
	return toDomainUser(createdUser), err
}

type JobPublisher interface {
	Publish(ctx context.Context, job Job) error
}

type Job struct {
	Name    string
	Payload map[string]any
}
type JobHandler struct {
	Service UserService
}

func (j JobHandler) ProccessSendMail(job *work.Job) error {

	email := job.ArgString("email")
	message := job.ArgString("message")
	// data := map[string]any{"Name": email, "message": message}
	// slog.Info("Processing email task after delay", "to", email)
	// err := j.Mail.Send(Request{
	// 	To:       []string{email},
	// 	Subject:  "Code Verification",
	// 	Body:     &message,
	// 	Data:     data,
	// 	MailType: Text,
	// })
	err := j.Service.SendVerificationMail(context.Background(), VerificationCodeMail{
		Message: message,
		Email:   email,
	})
	if err != nil {
		slog.Error("ProcessSendMail Failed", "error", err)
	}
	return err
}
