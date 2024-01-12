package v1

import (
	"github/yogabagas/join-app/transport/rest/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func NewAccessV1(h handler.HandlerImpl, r *mux.Router) {
	r.HandleFunc("/access", h.UpsertAccess).Methods(http.MethodPut)
	r.HandleFunc("/access/{type}", h.GetAccessByRoleUID).Methods(http.MethodGet)
}
