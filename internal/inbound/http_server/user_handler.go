package httpserver

import (
	"log/slog"
	"net/http"

	authApp "github.com/eng-yousef-khaled/expenses_api/internal/application/auth"
	authCore "github.com/eng-yousef-khaled/expenses_api/internal/core/auth"
	"github.com/eng-yousef-khaled/expenses_api/internal/inbound/json"
	"github.com/jackc/pgx/v5/pgtype"
)

type ErrorResponse struct {
	Error string `json:"error"`
}
type UserHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
	CheckEnteredVerificationCode(w http.ResponseWriter, r *http.Request)
}
type handler struct {
	service authApp.UserService
}
type CreateUserRequest struct {
	Uuid     pgtype.UUID `json:"uuid"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
}

type CheckEnteredVerificationCodeRequest struct {
	UserId int64  `json:"user_id"`
	Code   string `json:"code"`
}
type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateHandler(ser authApp.UserService) UserHandler {
	return &handler{
		service: ser,
	}
}

// RegisterUser godoc
// @Summary      Register a new user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body      CreateUserRequest  true  "Payload"
// @Success      201     {object}  authCore.User
// @Router       /auth/register [post]
func (s *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var raw CreateUserRequest
	if err := json.Read(r, &raw); err != nil {
		slog.Log(r.Context(), slog.LevelError, "convert request ", "error", err)
		json.Write(w, http.StatusBadRequest, ErrorResponse{Error: authCore.InvalidRequestBody.Error()})
		return
	}
	email, mailErr := authCore.NewEmailAddress(raw.Email)
	if mailErr != nil {
		slog.Log(r.Context(), slog.LevelError, "email is not valid", "error", mailErr)
		json.Write(w, http.StatusUnprocessableEntity, ErrorResponse{Error: mailErr.Error()})
		return
	}
	pass, passErr := authCore.NewRawPassword(raw.Password)
	if passErr != nil {
		slog.Log(r.Context(), slog.LevelError, "password is not valid", "error", passErr)
		json.Write(w, http.StatusUnprocessableEntity, ErrorResponse{Error: passErr.Error()})
		return
	}
	name, nameErr := authCore.NewName(raw.Name)
	if nameErr != nil {
		slog.Log(r.Context(), slog.LevelError, "name is not valid", "error", nameErr)
		json.Write(w, http.StatusUnprocessableEntity, ErrorResponse{Error: nameErr.Error()})
		return
	}
	user := authApp.CreateUser{
		Uuid:     raw.Uuid.Bytes,
		Name:     name,
		Email:    email,
		Password: pass,
	}
	createdUser, createUserError := s.service.CreateUser(r.Context(), user)
	if createUserError != nil {
		if createUserError == authCore.Duplicate {
			slog.Log(r.Context(), slog.LevelError, "user Duplicate error while creation", "error", createUserError)
			json.Write(w, http.StatusConflict, ErrorResponse{Error: authCore.Duplicate.Error()})
			return
		}
		slog.Log(r.Context(), slog.LevelError, "unknown error, may from  server", "error", createUserError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, createdUser)
}

func (h *handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var raw LoginUserRequest
	if err := json.Read(r, &raw); err != nil {
		slog.Log(r.Context(), slog.LevelError, "convert request data", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// loginUser := authApp.LoginUser{
	// 	Email:    authCore.EmailAddress(raw.Email),
	// 	Password: authCore.RawPassword(raw.Password),
	// }
	// u, err := h.service.LoginUser(r.Context(), loginUser)
	// if err != nil {
	// 	slog.Log(r.Context(), slog.LevelError, "Failed to login", "error", err)
	// 	json.Write(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	// 	return
	// }
	// json.Write(w, http.StatusOK, u)
}

func (h *handler) CheckEnteredVerificationCode(w http.ResponseWriter, r *http.Request) {
	var raw CheckEnteredVerificationCodeRequest
	if err := json.Read(r, &raw); err != nil {
		slog.Log(r.Context(), slog.LevelError, "convert request data", "error", err)
		json.Write(w, http.StatusBadRequest, ErrorResponse{Error: authCore.InvalidRequestBody.Error()})
		return
	}
	code, err := authCore.NewVerificationCode(raw.Code)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, "incoming data is not correct", "error", err)
		json.Write(w, http.StatusUnprocessableEntity, ErrorResponse{Error: err.Error()})
		return
	}
	userId, uIdErr := authCore.NewUserId(raw.UserId)
	if uIdErr != nil {
		slog.Log(r.Context(), slog.LevelError, "incoming data is not correct", "error", uIdErr)
		json.Write(w, http.StatusUnprocessableEntity, ErrorResponse{Error: uIdErr.Error()})
		return
	}
	slog.Log(r.Context(), slog.LevelError, "User Verification Code received", "payload", raw)
	enterCode := authApp.NewEnterCodeRequest(code)
	u, uErr := h.service.CheckEnteredVerificationCode(r.Context(), userId, enterCode)
	if uErr != nil {
		slog.Log(r.Context(), slog.LevelError, "error while checking code", "error", uErr)
		json.Write(w, http.StatusUnprocessableEntity, ErrorResponse{Error: uErr.Error()})
		return
	}
	json.Write(w, http.StatusOK, u)
}
