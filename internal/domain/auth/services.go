package auth

import (
	"context"
	"fmt"
	"log"
	"log/slog"
)

type UserService interface {
	CreateUser(ctx context.Context, u CreateUser) (User, *CreateUserError)
	SendVerificationMail(ctx context.Context, mailContent VerificationCodeMail) error
}
type svc struct {
	users     UserRepository
	publisher JobPublisher
	hasher    PasswordHash
	mailer    Mailer
}

func NewService(users UserRepository,
	publisher JobPublisher,
	hasher PasswordHash, mailer Mailer) UserService {
	return &svc{
		users:     users,
		publisher: publisher,
		hasher:    hasher,
		mailer:    mailer,
	}
}

func (s *svc) CreateUser(ctx context.Context, user CreateUser) (User, *CreateUserError) {

	password, hasherErr := s.hasher.HashingPassword(user.Password)
	if hasherErr != nil {
		slog.Log(ctx, slog.LevelError, "hashing error proccess is failed", "error", hasherErr)
		return User{}, &CreateUserError{PasswordHashing: string(user.Password)}
	}
	user.Password = password
	u, uError := s.users.CreateUser(ctx, user)
	if uError != nil {
		slog.Log(ctx, slog.LevelError, "Has an error while Create User ", "error", uError)
		return u, uError
	}
	err := s.publisher.Publish(ctx, Job{
		Name: "send_verification_mail",
		Payload: map[string]any{
			"email":   user.Email,
			"message": fmt.Sprintf("Welcome dear: %s, your UUID is: %s", user.Name, user.Uuid),
		},
	})
	if err != nil {
		log.Printf("Has an error while adding to Queue : %s \n", err)
	}
	return u, nil
}
func (s *svc) SendVerificationMail(ctx context.Context, mailContent VerificationCodeMail) error {
	slog.Info("Processing email task after delay", "to", mailContent.Email)

	slog.Info("Processing email task after delay", "to", mailContent.Email)
	data := map[string]any{"Name": mailContent.Email, "message": mailContent.Message}
	return s.mailer.Send(ctx, EmailMessage{
		To:      []string{mailContent.Email},
		Subject: "Code Verification",
		Body:    &mailContent.Message,
		Data:    data,
	})
}
