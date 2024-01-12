package controller

import (
	"context"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/service/jwk/usecase"
)

type JWKControllerImpl struct {
	jwkSvc usecase.JWKService
}

type JWKController interface {
	VerifyJWT(ctx context.Context, req service.VerifyTokenReq) (resp service.VerifyTokenResp, err error)
}

func NewJWKController(jwkSvc usecase.JWKService) JWKController {
	return &JWKControllerImpl{
		jwkSvc: jwkSvc,
	}
}

func (jc *JWKControllerImpl) VerifyJWT(ctx context.Context, req service.VerifyTokenReq) (resp service.VerifyTokenResp, err error) {
	return jc.jwkSvc.VerifyJWT(ctx, req)
}
