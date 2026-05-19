package auth

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user CreateUserRequest) (User, error)
}

type Mailer interface {
	Send(ctx context.Context, msg Request) error
}
type PasswordHash interface {
	HashingPassword(password *string) error
}
