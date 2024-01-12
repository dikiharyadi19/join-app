package handler

import (
	"encoding/json"
	"errors"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/transport/rest/handler/response"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateResources handler
// @Summary Create New Resources
// @Description New Resources Registration
// @Tags Resources
// @Produce json
// @Security ApiKeyAuth
// @Param users body service.CreateResourcesReq true "Request Create Resources"
// @Success 200 {object} response.JSONResponse().APIStatusCreated()
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/resources [POST]
func (h *HandlerImpl) CreateResources(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodPost {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	var req service.CreateResourcesReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(w)
		return
	}

	err := h.Controller.ResourcesController.CreateResources(r.Context(), req)
	if err != nil {
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(w)
		return
	}

	res.APIStatusCreated().Send(w)

}

func (h *HandlerImpl) GetResourcesByType(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodGet {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	vars := mux.Vars(r)
	t, ok := vars["type"]
	if !ok {
		res.SetError(response.ErrBadRequest).SetMessage(errors.New("resources type is missing").Error()).Send(w)
		return
	}

	ty, err := strconv.Atoi(t)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(w)
		return
	}

	req := service.GetResourcesByTypeReq{
		Type: ty,
	}

	resp, err := h.Controller.ResourcesController.GetResourcesByType(r.Context(), req)
	if err != nil {
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(w)
		return
	}

	res.SetData(resp).Send(w)
}
