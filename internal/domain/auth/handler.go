package auth

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/eng-yousef-khaled/expenses_api/internal/json"
)

type ErrorResponse struct {
	Error string `json:"error"`
}
type UserHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
}
type handler struct {
	service UserService
}

func CreateHandler(ser UserService) UserHandler {
	return &handler{
		service: ser,
	}
}
func (s *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var raw CreateUserRequest
	if err := json.Read(r, &raw); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email, mailErr := NewEmailAddress(raw.Email)
	if mailErr != nil {
		slog.Log(r.Context(), slog.LevelError, "email is not valid", "error", mailErr)
		json.Write(w, http.StatusBadRequest, ErrorResponse{Error: mailErr.Error()})
		return
	}
	pass, passErr := NewPassword(raw.Password)
	if passErr != nil {
		slog.Log(r.Context(), slog.LevelError, "password is not valid", "error", passErr)
		json.Write(w, http.StatusBadRequest, ErrorResponse{Error: passErr.Error()})
		return
	}
	name, nameErr := NewName(raw.Name)
	if nameErr != nil {
		slog.Log(r.Context(), slog.LevelError, "name is not valid", "error", nameErr)
		json.Write(w, http.StatusBadRequest, ErrorResponse{Error: nameErr.Error()})
		return
	}
	user := CreateUser{
		Uuid:     raw.Uuid.Bytes,
		Name:     name,
		Email:    email,
		Password: pass,
	}
	createdUser, createUserError := s.service.CreateUser(r.Context(), user)
	if createUserError != nil {
		if createUserError.Duplicate != "" {
			slog.Log(r.Context(), slog.LevelError, "user Duplicate error while creation", "error", createUserError)
			json.Write(w, http.StatusConflict, ErrorResponse{Error: "This uuid or email already in use please try different one"})
			return
		}
		slog.Log(r.Context(), slog.LevelError, "unknown error, may from  server", "error", createUserError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, createdUser)
}
