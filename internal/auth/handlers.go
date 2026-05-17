package auth

import (
	"log"
	"net/http"

	repo "github.com/eng-yousef-khaled/expenses_api/internal/adapters/postgresql/sqlc"
	"github.com/eng-yousef-khaled/expenses_api/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(ser Service) *handler {
	return &handler{
		service: ser,
	}
}
func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var u repo.User
	if err := json.Read(r, &u); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, createUserError := h.service.CreateUser(r.Context(), u)
	if createUserError != nil {
		log.Println(createUserError)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, user)
}
