package usecase

import (
	"context"
	"fmt"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/domain/repository/cache"
	"github/yogabagas/join-app/domain/repository/sql"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/service/resources/presenter"
	"github/yogabagas/join-app/shared/util"
)

type ResourcesServiceImpl struct {
	cache     cache.Cache
	repo      sql.RepositoryRegistry
	presenter presenter.ResourcesPresenter
}

type ResourcesService interface {
	CreateResources(ctx context.Context, req service.CreateResourcesReq) error
	GetResourcesByType(ctx context.Context, req service.GetResourcesByTypeReq) ([]service.GetResourcesByTypeResp, error)
}

func NewResourcesService(cache cache.Cache, repository sql.RepositoryRegistry, presenter presenter.ResourcesPresenter) ResourcesService {
	return &ResourcesServiceImpl{
		cache:     cache,
		repo:      repository,
		presenter: presenter}
}

func (rs *ResourcesServiceImpl) CreateResources(ctx context.Context, req service.CreateResourcesReq) error {

	resourcesRepo := rs.repo.ResourcesRepository()

	uID := util.NewULIDGenerate()

	return resourcesRepo.CreateResources(ctx, &model.Resource{
		UID:       uID,
		Name:      req.Name,
		Type:      req.Type,
		Action:    req.Action,
		ParentUID: req.ParentUID,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.CreatedBy,
	})
}

func (rs *ResourcesServiceImpl) GetResourcesByType(ctx context.Context, req service.GetResourcesByTypeReq) (resp []service.GetResourcesByTypeResp, err error) {

	resourcesRepo := rs.repo.ResourcesRepository()

	keyCache := fmt.Sprintf("resources::type:%d", req.Type)
	err = rs.cache.GetObject(ctx, keyCache, &resp)
	if err == nil {
		return
	}

	res, err := resourcesRepo.ReadResourcesByType(ctx, &model.ReadResourcesByTypeReq{
		Type: req.Type,
	})
	if err != nil {
		return nil, err
	}

	resp, err = rs.presenter.GetResourcesByType(ctx, res)
	if err != nil {
		return resp, err
	}

	go rs.cache.Set(ctx, keyCache, resp, 0)

	return resp, nil
}
