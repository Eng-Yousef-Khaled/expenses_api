package auth

import (
	"context"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/gocraft/work"
)

type UserService interface {
	CreateUser(ctx context.Context, u repo.User) (repo.User, error)
	ProccessSendMail(job *work.Job) error
}

type QueuePort interface {
	Enqueue(taskName string, payload map[string]any) error
}

type OVHTextMailer interface {
	SendMail(r Request) error
}
type JobHandler struct {
	Mail OVHTextMailer
}

type RequestType int

const (
	HTML RequestType = iota
	Text
)

type Request struct {
	To       []string
	Subject  string
	Body     *string
	TempFile *string
	Data     interface{}
	MailType RequestType
}
