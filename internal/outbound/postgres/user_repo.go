package userrepo

import (
	"context"
	"errors"
	"log/slog"
	"time"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/eng-yousef-khaled/expenses_api/internal/application/auth"
	authCore "github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"github.com/gocraft/work"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

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
		ID:             raw.ID,
		Uuid:           convertPgTypeUUIDToGoogle(raw.Uuid),
		Name:           authCore.Name(raw.Name),
		Email:          authCore.EmailAddress(raw.Email),
		Password:       authCore.Password{Value: raw.Password, Hashed: true},
		IsVerification: raw.IsVerification,
	}
}
func toSqlcParams(raw auth.CreateUser) repo.CreateUserParams {
	return repo.CreateUserParams{
		Uuid:     convertGoogleUUIDToPgType(raw.Uuid),
		Name:     string(raw.Name),
		Email:    string(raw.Email),
		Password: string(raw.Password.Value),
	}
}

func VerificationCodeToSqlcParams(raw authCore.UserVerificationCode) repo.CreateVerificationCodeParams {
	return repo.CreateVerificationCodeParams{
		Code:      string(raw.VerificationCode),
		UsersID:   pgtype.Int8{Int64: raw.UserID, Valid: true},
		ExpiresAt: pgtype.Timestamptz{Time: time.Time(raw.ExpiresAt), Valid: true},
	}
}

func toDomainUserVerificationCode(raw repo.UserVerificationCode) authCore.UserVerificationCode {
	return authCore.UserVerificationCode{
		ID:               raw.ID,
		UserID:           raw.UsersID.Int64,
		VerificationCode: authCore.VerificationCode(raw.Code),
		ExpiresAt:        authCore.VerificationCodeExpire(raw.ExpiresAt.Time),
		CreatedAt:        raw.CreatedAt.Time,
	}
}

type PostgresUserRepository interface {
	CreateUser(ctx context.Context, user auth.CreateUser) (authCore.User, error)
	LoginUser(ctx context.Context, email authCore.EmailAddress) (authCore.User, error)
	SaveVerification(ctx context.Context, mail authCore.EmailAddress, user_id int64, code string) (authCore.UserVerificationCode, error)
	CheckEnteredVerificationCode(ctx context.Context, userId int64, code auth.EnterCodeRequest) (authCore.User, error)
}

func NewRepo(q *repo.Queries) PostgresUserRepository {
	return &postgresUserRepository{
		q: q,
	}
}

type postgresUserRepository struct {
	q *repo.Queries
}

func (p *postgresUserRepository) CreateUser(ctx context.Context, user auth.CreateUser) (authCore.User, error) {
	user1 := toSqlcParams(user)
	createdUser, err := p.q.CreateUser(ctx, user1)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Error while creating user in adapter", "error", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return authCore.User{}, authCore.Duplicate
			}
		}
		return authCore.User{}, authCore.ServerError
	}
	return toDomainUser(createdUser), nil
}

func (p *postgresUserRepository) LoginUser(ctx context.Context, email authCore.EmailAddress) (authCore.User, error) {
	createdUser, err := p.q.LoginUser(ctx, string(email))
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Can't login has an", "error", err)
		return authCore.User{}, authCore.InvalidEmailOrPasswordError
	}
	return toDomainUser(createdUser), nil
}

func (p *postgresUserRepository) SaveVerification(ctx context.Context, mail authCore.EmailAddress, user_id int64, code string) (authCore.UserVerificationCode, error) {
	parm := VerificationCodeToSqlcParams(authCore.UserVerificationCode{UserID: user_id, VerificationCode: authCore.VerificationCode(code), ExpiresAt: authCore.VerificationCodeExpire(time.Now().Add(authCore.EXPIRES_MINUTES * time.Minute))})
	user_ver, err := p.q.CreateVerificationCode(ctx, parm)
	if err != nil {
		// TODO: readable error message
		// Like not correct code or user is not even register!!
		// in handler also should check user is null of not
		slog.Log(ctx, slog.LevelError, "User Verification Code error", "error", err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			slog.Log(ctx, slog.LevelError, "User Verification Code error", "error", pgErr.Code)
			// if pgErr.Code == "23505" {
			// 	return authCore.User{}, &authCore.CreateUserError{Duplicate: "a user with this UUID already exists"}
			// }
		}
		return authCore.UserVerificationCode{}, authCore.VerificationCodeError
	}

	return toDomainUserVerificationCode(user_ver), nil

}

func (p *postgresUserRepository) CheckEnteredVerificationCode(ctx context.Context, userId int64, code auth.EnterCodeRequest) (authCore.User, error) {
	param := repo.CheckUserVerificationCodeParams{
		ID:   userId,
		Code: string(code),
	}
	user, err := p.q.CheckUserVerificationCode(ctx, param)
	if err != nil {
		return authCore.User{}, err
	}
	return toDomainUser(user), nil
}

type JobHandler struct {
	VerificationNotifier auth.VerificationNotifier
}

func (j JobHandler) ProccessSendMail(job *work.Job) error {
	email := job.ArgString("email")
	content := job.ArgString("content")
	subject := job.ArgString("subject")
	code := job.ArgString("code")
	name := job.ArgString("name")

	slog.Info("ProcessSendMail vars", "email", email, "message", content)
	err := j.VerificationNotifier.SendVerification(context.Background(), authCore.EmailAddress(email), subject, content, code, name)

	if err != nil {
		slog.Error("ProcessSendMail Failed", "error", err)
	}
	return err
}
