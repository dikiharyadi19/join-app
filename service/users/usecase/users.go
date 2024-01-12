package usecase

import (
	"context"
	"errors"
	"github/yogabagas/join-app/config"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/domain/repository/cache"
	"github/yogabagas/join-app/domain/repository/sql"
	"github/yogabagas/join-app/domain/service"
	"github/yogabagas/join-app/service/users/presenter"
	"github/yogabagas/join-app/shared/util"
	"log"

	"github/yogabagas/join-app/shared/constant"
	"time"
)

type UsersServiceImpl struct {
	repo      sql.RepositoryRegistry
	cache     cache.Cache
	presenter presenter.UsersPresenter
}

type UsersService interface {
	CreateUsers(ctx context.Context, req service.CreateUsersReq) error
	GetUsersWithPagination(ctx context.Context, req service.GetUsersWithPaginationReq) (service.GetUsersWithPaginationResp, error)
}

func NewUsersService(repository sql.RepositoryRegistry, cache cache.Cache, presenter presenter.UsersPresenter) UsersService {
	return &UsersServiceImpl{
		repo:      repository,
		cache:     cache,
		presenter: presenter,
	}
}

func (us *UsersServiceImpl) CreateUsers(ctx context.Context, req service.CreateUsersReq) error {

	authzRepo := us.repo.AuthzRepository()
	rolesRepo := us.repo.RolesRepository()

	role, err := rolesRepo.ReadRolesByID(ctx, &model.ReadRolesByIDReq{
		ID: req.RoleID,
	})
	if err != nil {
		return err
	} else if role == nil {
		return errors.New("role is not found")
	}

	pwd, err := util.Hash(config.GlobalCfg.PasswordAlg, req.Password)
	if err != nil {
		return err
	}

	hbd, err := time.Parse(time.DateOnly, req.Birthdate)
	if err != nil {
		return err
	}

	userUID := util.NewULIDGenerate()

	var InTransaction = func(rr sql.RepositoryRegistry) (out interface{}, err error) {

		usersRepo := rr.UsersRepository()
		userCredentialsRepo := rr.UserCredentialsRepository()

		err = usersRepo.CreateUsers(ctx, &model.User{
			UID:         userUID,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Birthdate:   hbd,
			Email:       req.Email,
			Gender:      req.Gender,
			Country:     req.Country,
			Description: req.Bio,
			CreatedBy:   userUID,
			UpdatedBy:   userUID,
		})
		if err != nil {
			log.Println("error insert users", err)
			return nil, err
		}

		err = userCredentialsRepo.InsertCredential(ctx, &model.UserCredential{
			UserUID:  userUID,
			Username: req.Username,
			Password: util.Base64(pwd),
		})
		if err != nil {
			log.Println("error insert credentials", err)
			return nil, err
		}

		return nil, nil
	}

	_, err = us.repo.DoInTransaction(ctx, InTransaction)
	if err != nil {
		log.Println("error do in transaction", err)
		return err
	}

	authzUID := util.NewULIDGenerate()

	err = authzRepo.CreateAuthz(ctx, &model.Authz{
		UID:       authzUID,
		UserUID:   userUID,
		RoleUID:   role.UID,
		CreatedBy: userUID,
		UpdatedBy: userUID,
	})
	if err != nil {
		log.Println("error insert authz", err)
		return err
	}
	return nil
}

func (us *UsersServiceImpl) GetUsersWithPagination(ctx context.Context, req service.GetUsersWithPaginationReq) (resp service.GetUsersWithPaginationResp, err error) {

	usersRepo := us.repo.UsersRepository()

	users, err := usersRepo.ReadUsersWithPagination(ctx, &model.ReadUsersWithPaginationReq{
		Fullname: req.Fullname,
		Limit:    req.Limit,
		Offset:   util.PageToOffset(req.Limit, req.Page),
	})
	if err != nil {
		return
	}

	count, err := usersRepo.CountUsers(ctx, &model.CountUsersReq{
		IsDeleted: constant.False.Int(),
	})
	if err != nil {
		return
	}

	return us.presenter.GetUsersWithPagination(ctx, req, users, count)
}
