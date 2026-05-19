package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
)

type UserService interface {
	CreateUser(ctx context.Context, u CreateUserRequest) (User, error)
	SendVerificationMail(ctx context.Context, mailContent VerificationCodeMail) error
}
type svc struct {
	users     UserRepository
	publisher JobPublisher
	hasher    PasswordHash
	mailer    Mailer
}

func CreateService(users UserRepository,
	publisher JobPublisher,
	hasher PasswordHash, mailer Mailer) UserService {
	return &svc{
		users:     users,
		publisher: publisher,
		hasher:    hasher,
		mailer:    mailer,
	}
}

func (s *svc) CreateUser(ctx context.Context, user CreateUserRequest) (User, error) {

	hashErr := s.hasher.HashingPassword(&user.Password)
	if hashErr != nil {
		log.Println(hashErr)
		return User{}, errors.New("Fail to encrypt password")
	}
	u, uError := s.users.CreateUser(ctx, user)
	if uError != nil {
		log.Printf("Has an error while Create User : %s \n", uError)
	}
	err := s.publisher.Publish(ctx, Job{
		Name: "send_welcome_email",
		Payload: map[string]any{
			"email":   user.Email,
			"message": fmt.Sprintf("Welcome dear: %s, your UUID is: %s", user.Name, user.Uuid),
		},
	})
	if err != nil {
		log.Printf("Has an error while adding to Queue : %s \n", err)
	}
	return u, err
}
func (s *svc) SendVerificationMail(ctx context.Context, mailContent VerificationCodeMail) error {
	slog.Info("Processing email task after delay", "to", mailContent.Email)

	slog.Info("Processing email task after delay", "to", mailContent.Email)
	data := map[string]any{"Name": mailContent.Email, "message": mailContent.Message}
	return s.mailer.Send(ctx, Request{
		To:       []string{mailContent.Email},
		Subject:  "Code Verification",
		Body:     &mailContent.Message,
		Data:     data,
		MailType: Text,
	})
}
