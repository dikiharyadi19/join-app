package repository

import (
	"context"
	"github/yogabagas/join-app/domain/model"
)

type JWKRepository interface {
	UpsertJWK(context.Context, *model.JWK) error
	ReadUnexpiredKeyByID(context.Context, *model.ReadUnexpiredKeyByIDReq) (*model.ReadUnexpiredKeyByIDResp, error)
	ReadUnexpiredKeys(context.Context) ([]*model.ReadUnexpiredKeyResp, error)
}
