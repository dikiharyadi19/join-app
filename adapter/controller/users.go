package controller

import (
	"context"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/service/users/usecase"
)

type UsersControllerImpl struct {
	usersSvc usecase.UsersService
}

type UsersController interface {
	CreateUsers(ctx context.Context, req service.CreateUsersReq) error
	GetUsersWithPagination(ctx context.Context, req service.GetUsersWithPaginationReq) (service.GetUsersWithPaginationResp, error)
}

func NewUsersController(userSvc usecase.UsersService) UsersController {
	return &UsersControllerImpl{usersSvc: userSvc}
}

func (uc *UsersControllerImpl) CreateUsers(ctx context.Context, req service.CreateUsersReq) error {
	return uc.usersSvc.CreateUsers(ctx, req)
}

func (uc *UsersControllerImpl) GetUsersWithPagination(ctx context.Context, req service.GetUsersWithPaginationReq) (service.GetUsersWithPaginationResp, error) {
	return uc.usersSvc.GetUsersWithPagination(ctx, req)
}
