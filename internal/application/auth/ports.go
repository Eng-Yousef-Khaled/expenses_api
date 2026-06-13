package auth

import (
	"context"

	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user CreateUser) (auth.User, error)
	LoginUser(ctx context.Context, email auth.EmailAddress) (auth.User, error)
	SaveVerification(ctx context.Context, mail auth.EmailAddress, user_id int64, code string) (auth.UserVerificationCode, error)
	CheckEnteredVerificationCode(ctx context.Context, userId int64, code EnterCodeRequest) (auth.User, error)
}

type VerificationNotifier interface {
	SendVerification(ctx context.Context, mail auth.EmailAddress, subject string, content string, code string, name string) error
}
type PasswordHash interface {
	HashingPassword(password auth.Password) (auth.Password, error)
}
type Job struct {
	Name    string
	Payload map[string]any
}
type JobPublisher interface {
	Publish(ctx context.Context, job Job) error
}
