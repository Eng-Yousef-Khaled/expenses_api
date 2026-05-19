package auth

import "github.com/jackc/pgx/v5/pgtype"

type CreateUserRequest struct {
	ID       int64       `json:"id"`
	Uuid     pgtype.UUID `json:"uuid"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}

type User struct {
	ID       int64       `json:"id"`
	Uuid     pgtype.UUID `json:"uuid"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}
