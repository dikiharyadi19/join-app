package repository

import (
	"context"
	"github/yogabagas/join-app/domain/model"
)

type AccessRepository interface {
	UpsertAccess(ctx context.Context, req []*model.Access) error
	ReadAccessByRoleUID(ctx context.Context, req *model.ReadAccessByRoleUIDReq) ([]*model.ReadAccessByRoleUIDResp, error)
}
