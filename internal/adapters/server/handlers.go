package server

import (
	"log"
	"net/http"

	"github.com/eng-yousef-khaled/expenses_api/internal/domain/auth"
	"github.com/eng-yousef-khaled/expenses_api/internal/json"
)

type httpServer struct {
	userService auth.UserService
}

func NewHttpServer(ser auth.UserService) *httpServer {
	return &httpServer{
		userService: ser,
	}
}
func (h *httpServer) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var u auth.CreateUserRequest
	if err := json.Read(r, &u); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, createUserError := h.userService.CreateUser(r.Context(), u)
	if createUserError != nil {
		log.Println(createUserError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, user)
}
