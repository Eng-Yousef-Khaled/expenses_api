package userrepo

import (
	"context"
	"errors"
	"log/slog"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/eng-yousef-khaled/expenses_api/internal/application/auth"
	authCore "github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserRequest struct {
	ID       int64       `json:"id"`
	Uuid     pgtype.UUID `json:"uuid"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}
type DBUser struct {
	ID       int64       `json:"id"`
	Uuid     pgtype.UUID `json:"uuid"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}

func convertPgTypeUUIDToGoogle(raw pgtype.UUID) uuid.UUID {
	if !raw.Valid {
		return uuid.Nil
	}

	return uuid.UUID(raw.Bytes)
}
func convertGoogleUUIDToPgType(raw uuid.UUID) pgtype.UUID {
	if raw == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}

	return pgtype.UUID{
		Bytes: raw,
		Valid: true,
	}
}

func toDomainUser(raw repo.User) authCore.User {
	return authCore.User{
		ID:       raw.ID,
		Uuid:     convertPgTypeUUIDToGoogle(raw.Uuid),
		Name:     authCore.Name(raw.Name),
		Email:    authCore.EmailAddress(raw.Email),
		Password: authCore.HashedPassword(raw.Password),
	}
}
func toSqlcParams(raw authCore.User) repo.CreateUserParams {
	return repo.CreateUserParams{
		Uuid:     convertGoogleUUIDToPgType(raw.Uuid),
		Name:     string(raw.Name),
		Email:    string(raw.Email),
		Password: string(raw.Password),
	}
}

type PostgresUserRepository interface {
	CreateUser(ctx context.Context, user authCore.User) (authCore.User, *authCore.CreateUserError)
}

func NewRepo(q *repo.Queries) PostgresUserRepository {
	return &postgresUserRepository{
		q: q,
	}
}

type postgresUserRepository struct {
	q *repo.Queries
}

func (p *postgresUserRepository) CreateUser(ctx context.Context, user authCore.User) (authCore.User, *authCore.CreateUserError) {
	user1 := toSqlcParams(user)
	createdUser, err := p.q.CreateUser(ctx, user1)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Error while creating user in adapter", "error", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return authCore.User{}, &authCore.CreateUserError{Duplicate: "a user with this UUID already exists"}
			}
		}
		return authCore.User{}, &authCore.CreateUserError{
			Unknown: "failed to register user",
		}
	}
	return toDomainUser(createdUser), nil
}

type JobHandler struct {
	Service auth.UserService
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
	err := j.Service.SendVerification(context.Background(), authCore.EmailAddress(email), message)
	if err != nil {
		slog.Error("ProcessSendMail Failed", "error", err)
	}
	return err
}
