package auth

import (
	"github.com/google/uuid"
)

type User struct {
	ID       int64
	Uuid     uuid.UUID
	Name     Name
	Email    EmailAddress
	Password HashedPassword
}

// type EmailMessage struct {
// 	To      []string
// 	Subject string
// 	Body    *string
// 	Data    interface{}
// }

// type VerificationCodeMail struct {
// 	Email   string
// 	Message string
// }
