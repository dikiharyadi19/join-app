package presenter

import (
	"context"
	"errors"
	"github/yogabagas/join-app/domain/service"
	"time"
)

type JWKPresenterImpl struct{}

type JWKPresenter interface {
	VerifyJWT(ctx context.Context, payload map[string]interface{}) (service.VerifyTokenResp, error)
}

func NewJWKPresenter() JWKPresenter {
	return &JWKPresenterImpl{}
}

func (jp *JWKPresenterImpl) VerifyJWT(ctx context.Context, payload map[string]interface{}) (resp service.VerifyTokenResp, err error) {

	sub, ok := payload["sub"].(string)
	if !ok {
		return resp, errors.New("subject is nil")
	}

	exp, ok := payload["exp"].(float64)
	if !ok {
		return resp, errors.New("expired unset")
	}

	roleUID, ok := payload["role_uid"].(string)
	if !ok {
		return resp, errors.New("role uid is undefined")
	}

	lat, ok := payload["last_active"].(float64)
	if !ok {
		return resp, errors.New("invalid token")
	}

	if time.Now().After(time.Unix(int64(exp), 0)) {
		return resp, errors.New("token expired")
	}

	return service.VerifyTokenResp{
		Valid:      true,
		UserUID:    sub,
		RoleUID:    roleUID,
		LastActive: time.Unix(int64(lat), 0).UTC(),
		ExpiredAt:  time.Unix(int64(exp), 0).UTC(),
	}, nil

}
