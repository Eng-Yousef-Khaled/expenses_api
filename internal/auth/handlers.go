package users

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
	// mailReq := mail.NewRequest([]string{u.Email}, fmt.Sprintf("مرحبا بك عزيزي %s", u.Email))
	// res := h.mailService.SendMail(*mailReq)
	// if res {
	// 	log.Printf("Send Mail succeeded")
	// } else {
	// 	log.Printf("Send mail Failed")
	// }
	err := HashingPassword(&u.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.service.CreateUser(r.Context(), u)
	json.Write(w, http.StatusCreated, u)
}
