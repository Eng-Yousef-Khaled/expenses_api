package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
)

type UserService interface {
	CreateUser(ctx context.Context, u CreateUser) (auth.User, error)
	CheckEnteredVerificationCode(ctx context.Context, userId int64, code EnterCodeRequest) (auth.User, error)
	LoginUser(ctx context.Context, u LoginUser) (auth.User, error)
}
type svc struct {
	users     UserRepository
	publisher JobPublisher
	hasher    PasswordHash
}

func NewService(users UserRepository,
	publisher JobPublisher,
	hasher PasswordHash) UserService {
	return &svc{
		users:     users,
		publisher: publisher,
		hasher:    hasher,
	}
}

func (s *svc) CreateUser(ctx context.Context, input CreateUser) (auth.User, error) {

	password, hasherErr := s.hasher.HashingPassword(input.Password)
	if hasherErr != nil {
		slog.Log(ctx, slog.LevelError, "hashing error proccess is failed", "error", hasherErr)
		return auth.User{}, auth.ServerError
	}
	user := NewCreateUser(input, password)
	u, uError := s.users.CreateUser(ctx, user)
	if uError != nil {
		slog.Log(ctx, slog.LevelError, "Has an error while Create User ", "error", uError)
		return u, uError
	}
	code, codeError := GenerateVerificationCode()
	if codeError != nil {
		return auth.User{}, auth.VerificationCodeCantBeGeneratedError
	}
	verCode, verErr := s.users.SaveVerification(ctx, u.Email, u.ID, code)
	if verErr != nil {
		slog.Info("error in proccess of send verification code in services", "error", verErr)
		return auth.User{}, verErr
	}
	expiresIn := int(time.Until(time.Time(verCode.ExpiresAt)).Minutes())
	content := fmt.Sprintf("Welcome to Expenses App, this code is valid for %v Min", expiresIn)
	pubErr := s.publisher.Publish(ctx, Job{
		Name: "send_verification_mail",
		Payload: map[string]any{
			"email":   u.Email,
			"content": content,
			"subject": "User Verification Code",
			"code":    string(verCode.VerificationCode),
			"name":    string(u.Name),
		},
	})
	if pubErr != nil {
		slog.Error("Has an error while adding to Queue ", "error", verErr)
		return auth.User{}, auth.ServerError
	}
	return u, nil
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
	if user.Password.Value != password.Value {
		slog.Info("password not matching", "password", password)
		return auth.User{}, auth.InvalidEmailOrPasswordError
	}
	return user, nil
}

func (s *svc) CheckEnteredVerificationCode(ctx context.Context, userId int64, code EnterCodeRequest) (auth.User, error) {
	u, err := s.users.CheckEnteredVerificationCode(ctx, userId, code)
	if err != nil {
		return auth.User{}, err
	}
	return u, nil
}
func GenerateVerificationCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", n.Int64())

	return code, nil
}
