package handler

import (
	"encoding/json"
	"errors"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/shared/constant"
	"github/yogabagas/join-app/transport/rest/handler/response"
	"net/http"

	"github.com/gorilla/mux"
)

// UpsertAccess handler
// @Summary UpsertAccess
// @Description UpsertAccess for update and insert existing/new access
// @Tags Access
// @Produce json
// @Security ApiKeyAuth
// @Param access body service.UpsertAccessReq true "Request Upsert Access"
// @Success 200 {object} response.JSONResponse().APIStatusCreated()
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/access [PUT]
func (h *HandlerImpl) UpsertAccess(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodPut {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	var req service.UpsertAccessReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(w)
		return
	}

	err := h.Controller.AccessController.UpsertAccess(r.Context(), req)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(w)
		return
	}

	res.APIStatusCreated().Send(w)
}

// GetAccessByRoleUID handler
// @Summary GetAccessByRoleUID
// @Description GetAccessByRoleUID for get access by role uid
// @Tags Access
// @Produce json
// @Security ApiKeyAuth
// @Param type path string true "resource type"
// @Success 200 {object} response.JSONResponse{data=[]service.GetAccessByRoleUIDResp}
// @Failure 400 {object} response.JSONResponse
// @Failure 500 {object} response.JSONResponse
// @Router /v1/access/{type} [GET]
func (h *HandlerImpl) GetAccessByRoleUID(w http.ResponseWriter, r *http.Request) {
	res := response.NewJSONResponse()

	if r.Method != http.MethodGet {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	claims := r.Context().Value(constant.Claim).(service.JWTClaims)

	vars := mux.Vars(r)
	t, ok := vars["type"]
	if !ok {
		res.SetError(response.ErrBadRequest).SetMessage(errors.New("resources type is missing").Error()).Send(w)
		return
	}

	req := service.GetAccessByRoleUIDReq{
		RoleUID: claims.RoleUID,
		Type:    constant.ResourceTypeAtoi(t).Int(),
	}

	resp, err := h.Controller.AccessController.GetAccessByRoleUID(r.Context(), req)
	if err != nil {
		res.SetError(response.ErrInternalServerError).SetMessage(err.Error()).Send(w)
		return
	}

	res.SetData(resp).Send(w)
}
