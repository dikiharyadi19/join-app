package controller

import (
	"context"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/service/modules/usecase"
	"github/yogabagas/join-app/shared/util"
)

type ModulesControllerImpl struct {
	modulesSvc usecase.ModulesService
}

type ModulesController interface {
	CreateModules(ctx context.Context, req service.CreateModulesReq, userData *util.UserData) error
}

func NewModulesController(modulesService usecase.ModulesService) ModulesController {
	return &ModulesControllerImpl{modulesSvc: modulesService}
}

func (cs *ModulesControllerImpl) CreateModules(ctx context.Context, req service.CreateModulesReq, userData *util.UserData) error {
	return cs.modulesSvc.CreateModules(ctx, req, userData)
}
