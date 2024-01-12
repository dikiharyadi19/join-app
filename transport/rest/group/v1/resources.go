package v1

import (
	"github/yogabagas/join-app/transport/rest/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func NewResourcesV1(h handler.HandlerImpl, r *mux.Router) {
	r.HandleFunc("/resources", h.CreateResources).Methods(http.MethodPost)
	r.HandleFunc("/resources/{type}", h.GetResourcesByType).Methods(http.MethodGet)
}
