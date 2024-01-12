package controller

import (
	"context"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/service/access/usecase"
)

type AccessControllerImpl struct {
	accessSvc usecase.AccessService
}

type AccessController interface {
	UpsertAccess(ctx context.Context, req service.UpsertAccessReq) error
	GetAccessByRoleUID(ctx context.Context, req service.GetAccessByRoleUIDReq) ([]service.GetAccessByRoleUIDResp, error)
}

func NewAccessController(accessSvc usecase.AccessService) AccessController {
	return &AccessControllerImpl{accessSvc: accessSvc}
}

func (ac *AccessControllerImpl) UpsertAccess(ctx context.Context, req service.UpsertAccessReq) error {
	return ac.accessSvc.UpsertAccess(ctx, req)
}

func (ac *AccessControllerImpl) GetAccessByRoleUID(ctx context.Context, req service.GetAccessByRoleUIDReq) ([]service.GetAccessByRoleUIDResp, error) {
	return ac.accessSvc.GetAccessByRoleUID(ctx, req)
}
