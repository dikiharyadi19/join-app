package usecase

import (
	"context"
	"fmt"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/domain/repository/cache"
	"github/yogabagas/join-app/domain/repository/sql"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/service/access/presenter"
	"github/yogabagas/join-app/shared/util"
)

type AccessServiceImpl struct {
	repo      sql.RepositoryRegistry
	cache     cache.Cache
	presenter presenter.AccessPresenter
}

type AccessService interface {
	UpsertAccess(ctx context.Context, req service.UpsertAccessReq) error
	GetAccessByRoleUID(ctx context.Context, req service.GetAccessByRoleUIDReq) ([]service.GetAccessByRoleUIDResp, error)
}

func NewAccessService(repository sql.RepositoryRegistry, cache cache.Cache, presenter presenter.AccessPresenter) AccessService {
	return &AccessServiceImpl{
		repo:      repository,
		cache:     cache,
		presenter: presenter}
}

func (as *AccessServiceImpl) UpsertAccess(ctx context.Context, req service.UpsertAccessReq) error {

	accessRepo := as.repo.AccessRepository()

	accessReqs := []*model.Access{}

	if len(req.ResourceUID) > 0 {

		for _, resourceID := range req.ResourceUID {
			accessReq := &model.Access{
				UID:         util.NewULIDGenerate(),
				RoleUID:     req.RoleUID,
				ResourceUID: resourceID,
				CreatedBy:   req.CreatedBy,
				UpdatedBy:   req.UpdatedBy,
			}
			accessReqs = append(accessReqs, accessReq)
		}
	}

	return accessRepo.UpsertAccess(ctx, accessReqs)
}

func (as *AccessServiceImpl) GetAccessByRoleUID(ctx context.Context, req service.GetAccessByRoleUIDReq) (resp []service.GetAccessByRoleUIDResp, err error) {

	accessRepo := as.repo.AccessRepository()

	keyCache := fmt.Sprintf("resources::role-uid:%s:type:%d", req.RoleUID, req.Type)
	err = as.cache.GetObject(ctx, keyCache, &resp)
	if err == nil {
		return
	}

	res, err := accessRepo.ReadAccessByRoleUID(ctx, &model.ReadAccessByRoleUIDReq{
		RoleUID: req.RoleUID,
		Type:    req.Type,
	})
	if err != nil {
		return nil, err
	}

	resp, err = as.presenter.GetAccessByRoleUID(ctx, res)
	if err != nil {
		return resp, err
	}

	go as.cache.Set(ctx, keyCache, resp, 0)

	return resp, nil

}
