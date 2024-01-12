package handler

import (
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/transport/rest/handler/response"
	"net/http"
)

func (h *HandlerImpl) VerifyJWT(w http.ResponseWriter, r *http.Request) {

	res := response.NewJSONResponse()

	if r.Method != http.MethodGet {
		res.SetError(response.ErrMethodNotAllowed).Send(w)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		res.SetError(response.ErrBadRequest).SetMessage("token can't be empty").Send(w)
		return
	}

	req := service.VerifyTokenReq{
		Token: token,
	}

	resp, err := h.Controller.JWKController.VerifyJWT(r.Context(), req)
	if err != nil {
		res.SetError(response.ErrBadRequest).SetMessage(err.Error()).Send(w)
		return
	}

	res.APIStatusAccepted().SetData(resp).Send(w)
}
