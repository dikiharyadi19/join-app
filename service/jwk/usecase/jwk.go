package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github/yogabagas/join-app/domain/repository/cache"
	"github/yogabagas/join-app/domain/repository/sql"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/service/jwk/presenter"
	"github/yogabagas/join-app/shared/util"
	"log"
	"strings"

	"github.com/go-jose/go-jose/v3"
)

type JWKServiceImpl struct {
	repo      sql.RepositoryRegistry
	cache     cache.Cache
	presenter presenter.JWKPresenter
}

type JWKService interface {
	VerifyJWT(ctx context.Context, req service.VerifyTokenReq) (resp service.VerifyTokenResp, err error)
}

func NewJWKService(repository sql.RepositoryRegistry, cache cache.Cache, presenter presenter.JWKPresenter) JWKService {
	return &JWKServiceImpl{
		repo:      repository,
		cache:     cache,
		presenter: presenter,
	}
}

func (js *JWKServiceImpl) VerifyJWT(ctx context.Context, req service.VerifyTokenReq) (resp service.VerifyTokenResp, err error) {

	jwkRepo := js.repo.JWKRepository()

	token, err := util.SplitBearer(req.Token)
	if err != nil {
		return resp, err
	}

	b, err := base64.RawStdEncoding.DecodeString(strings.Split(token, ".")[0])
	if err != nil {
		return resp, err
	}

	headerToken := make(map[string]string)
	if err := json.Unmarshal(b, &headerToken); err != nil {
		return resp, err
	}

	kid, ok := headerToken["kid"]
	if !ok {
		return resp, fmt.Errorf("token key ID is missing %s", kid)
	}

	object, err := jose.ParseSigned(token)
	if err != nil {
		return resp, err
	}

	keys, err := jwkRepo.ReadUnexpiredKeys(ctx)
	if err != nil {
		return
	}

	if len(keys) <= 0 {
		return resp, errors.New("key has expired")
	}

	var key interface{}
	for _, k := range keys {

		m := jose.JSONWebKey{}

		if err = json.Unmarshal(k.Key.([]byte), &m); err != nil {
			return
		}

		if m.KeyID == kid {
			key = m
		}

	}

	if key == nil {
		return resp, errors.New("key not found")
	}

	pb, err := object.Verify(key)
	if err != nil {
		log.Println("error verify object", err)
		return resp, err
	}

	payload := make(map[string]interface{})
	if err = json.Unmarshal(pb, &payload); err != nil {
		return resp, err
	}

	return js.presenter.VerifyJWT(ctx, payload)

}
