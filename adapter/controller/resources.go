package controller

import (
	"context"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/service/resources/usecase"
)

type ResourcesControllerImpl struct {
	resourcesSvc usecase.ResourcesService
}

type ResourcesController interface {
	CreateResources(ctx context.Context, req service.CreateResourcesReq) error
	GetResourcesByType(ctx context.Context, req service.GetResourcesByTypeReq) ([]service.GetResourcesByTypeResp, error)
}

func NewResourcesController(resourcesSvc usecase.ResourcesService) ResourcesController {
	return &ResourcesControllerImpl{resourcesSvc: resourcesSvc}
}

func (rc *ResourcesControllerImpl) CreateResources(ctx context.Context, req service.CreateResourcesReq) error {
	return rc.resourcesSvc.CreateResources(ctx, req)
}

func (rc *ResourcesControllerImpl) GetResourcesByType(ctx context.Context, req service.GetResourcesByTypeReq) ([]service.GetResourcesByTypeResp, error) {
	return rc.resourcesSvc.GetResourcesByType(ctx, req)
}
