package auth

import (
	"context"

	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user auth.User) (auth.User, *auth.CreateUserError)
	LoginUser(ctx context.Context, email auth.EmailAddress) (auth.User, error)
	SendVerification(ctx context.Context, user_id int64, code int32) (auth.UserVerificationCode, error)
}

type VerificationNotifier interface {
	SendVerification(ctx context.Context, mail auth.EmailAddress, message string) error
}
type PasswordHash interface {
	HashingPassword(password auth.RawPassword) (auth.HashedPassword, error)
}
type Job struct {
	Name    string
	Payload map[string]any
}
type JobPublisher interface {
	Publish(ctx context.Context, job Job) error
}
