package v1

import (
	"github/yogabagas/join-app/transport/rest/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func NewUsersV1(h handler.HandlerImpl, r *mux.Router) {
	r.HandleFunc("/users", h.CreateUsers).Methods(http.MethodPost)
	r.HandleFunc("/users", h.GetUsersWithPagination).Methods(http.MethodGet)
}
