package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
)

type UserService interface {
	CreateUser(ctx context.Context, u CreateUser) (auth.User, *auth.CreateUserError)
	LoginUser(ctx context.Context, u LoginUser) (auth.User, error)
	SendVerification(ctx context.Context, mail auth.EmailAddress, message string) error
}
type svc struct {
	users     UserRepository
	publisher JobPublisher
	hasher    PasswordHash
	mailer    VerificationNotifier
}

func NewService(users UserRepository,
	publisher JobPublisher,
	hasher PasswordHash, mailer VerificationNotifier) UserService {
	return &svc{
		users:     users,
		publisher: publisher,
		hasher:    hasher,
		mailer:    mailer,
	}
}

func (s *svc) CreateUser(ctx context.Context, input CreateUser) (auth.User, *auth.CreateUserError) {

	password, hasherErr := s.hasher.HashingPassword(input.Password)
	if hasherErr != nil {
		slog.Log(ctx, slog.LevelError, "hashing error proccess is failed", "error", hasherErr)
		return auth.User{}, &auth.CreateUserError{PasswordHashing: string(input.Password)}
	}
	user := auth.User{
		Uuid:     input.Uuid,
		Name:     input.Name,
		Email:    input.Email,
		Password: password,
	}
	u, uError := s.users.CreateUser(ctx, user)
	if uError != nil {
		slog.Log(ctx, slog.LevelError, "Has an error while Create User ", "error", uError)
		return u, uError
	}
	err := s.publisher.Publish(ctx, Job{
		Name: "send_verification_mail",
		Payload: map[string]any{
			"email":   input.Email,
			"message": fmt.Sprintf("Welcome dear: %s, your UUID is: %s", input.Name, input.Uuid),
		},
	})
	if err != nil {
		slog.Error("Has an error while adding to Queue ", "error", err)
	}
	return u, nil
}
func (s *svc) SendVerification(ctx context.Context, mail auth.EmailAddress, message string) error {
	slog.Info("Processing email task", "to", mail)
	return s.mailer.SendVerification(ctx, mail, message)
}

func (s *svc) LoginUser(ctx context.Context, u LoginUser) (auth.User, error) {
	user, err := s.users.LoginUser(ctx, u.Email)

	if err != nil {
		slog.Info("User fail to login with this", "with this mail:", u.Email)
		return auth.User{}, err
	}
	password, hasherErr := s.hasher.HashingPassword(u.Password)
	if hasherErr != nil {
		slog.Log(ctx, slog.LevelError, "hashing error proccess is failed", "error", hasherErr)
		return auth.User{}, auth.LoginPasswordHashingError
	}
	if user.Password != password {
		slog.Info("password not matching", "password", password)
		return auth.User{}, auth.InvalidEmailOrPasswordError
	}
	return user, nil
}
