package v1

import (
	"github/yogabagas/join-app/transport/rest/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func NewModulesV1(h handler.HandlerImpl, r *mux.Router) {
	r.HandleFunc("/modules", h.CreateModules).Methods(http.MethodPost)
	r.HandleFunc("/modules", h.GetModulesWithPagination).Methods(http.MethodGet)
	r.HandleFunc("/modules/download/file", h.GetModulesWithPagination).Methods(http.MethodGet)
	r.HandleFunc("/modules/{uid}", h.UpdateCourses).Methods(http.MethodPost)
	r.HandleFunc("/modules/{uid}", h.DeleteCourse).Methods(http.MethodDelete)
}
