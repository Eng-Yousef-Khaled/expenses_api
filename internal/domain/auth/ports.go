package auth

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user CreateUser) (User, *CreateUserError)
}

type Mailer interface {
	Send(ctx context.Context, msg EmailMessage) error
}
type PasswordHash interface {
	HashingPassword(password *Password) error
}
