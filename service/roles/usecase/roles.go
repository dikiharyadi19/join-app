package usecase

import (
	"context"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/domain/repository/sql"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/shared/util"
)

type RolesServiceImpl struct {
	repo sql.RepositoryRegistry
}

type RolesService interface {
	CreateRoles(ctx context.Context, req service.CreateRolesReq) error
}

func NewRolesService(repository sql.RepositoryRegistry) RolesService {
	return &RolesServiceImpl{repo: repository}
}

func (rs *RolesServiceImpl) CreateRoles(ctx context.Context, req service.CreateRolesReq) error {

	rolesRepo := rs.repo.RolesRepository()

	uID := util.NewULIDGenerate()

	return rolesRepo.CreateRoles(ctx, &model.Role{
		UID:       uID,
		Name:      req.Name,
		CreatedBy: req.CreatedBy,
		UpdatedBy: req.CreatedBy,
	})
}
