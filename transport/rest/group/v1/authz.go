package v1

import (
	"github/yogabagas/join-app/transport/rest/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAuthzV1(h handler.HandlerImpl, r *mux.Router) {
	r.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", h.Logout).Methods(http.MethodDelete)
}
