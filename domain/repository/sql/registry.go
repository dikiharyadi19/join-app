package sql

import (
	"context"
	"database/sql"
	accessRepo "github/yogabagas/join-app/service/access/repository"
	authzRepo "github/yogabagas/join-app/service/authz/repository"
	jwkRepo "github/yogabagas/join-app/service/jwk/repository"
	modulesRepo "github/yogabagas/join-app/service/modules/repository"
	resourcesRepo "github/yogabagas/join-app/service/resources/repository"
	rolesRepo "github/yogabagas/join-app/service/roles/repository"
	userCredentialsRepo "github/yogabagas/join-app/service/userCredentials/repository"
	usersRepo "github/yogabagas/join-app/service/users/repository"
)

type InTransaction func(RepositoryRegistry) (interface{}, error)

type RepositoryRegistryImpl struct {
	db         *sql.DB
	tx         *sql.Tx
	dbExecutor DBExecutor
}

func (r RepositoryRegistryImpl) AccessRepository() accessRepo.AccessRepository {
	if r.dbExecutor != nil {
		return NewAccessRepository(r.dbExecutor)
	}
	return NewAccessRepository(r.db)
}

func (r RepositoryRegistryImpl) AuthzRepository() authzRepo.AuthzRepository {
	if r.dbExecutor != nil {
		return NewAuthzRepository(r.dbExecutor)
	}
	return NewAuthzRepository(r.db)
}

func (r RepositoryRegistryImpl) JWKRepository() jwkRepo.JWKRepository {
	if r.dbExecutor != nil {
		return NewJWKRepository(r.dbExecutor)
	}
	return NewJWKRepository(r.db)
}

func (r RepositoryRegistryImpl) RolesRepository() rolesRepo.RolesRepository {
	if r.dbExecutor != nil {
		return NewRolesRepository(r.dbExecutor)
	}
	return NewRolesRepository(r.db)
}

func (r RepositoryRegistryImpl) ResourcesRepository() resourcesRepo.ResourcesRepository {
	if r.dbExecutor != nil {
		return NewResourcesRepository(r.dbExecutor)
	}
	return NewResourcesRepository(r.db)
}

func (r RepositoryRegistryImpl) UserCredentialsRepository() userCredentialsRepo.UserCredentialsRepository {
	if r.dbExecutor != nil {
		return NewUserCredentialsRepository(r.dbExecutor)
	}
	return NewUserCredentialsRepository(r.db)
}

func (r RepositoryRegistryImpl) UsersRepository() usersRepo.UsersRepository {
	if r.dbExecutor != nil {
		return NewUsersRepository(r.dbExecutor)
	}
	return NewUsersRepository(r.db)
}

func (r RepositoryRegistryImpl) ModulesRepository() modulesRepo.ModulesRepository {
	if r.dbExecutor != nil {
		return NewModulesRepository(r.dbExecutor)
	}
	return NewModulesRepository(r.db)
}

func (r RepositoryRegistryImpl) DoInTransaction(ctx context.Context, txFunc InTransaction) (out interface{}, err error) {
	var tx *sql.Tx

	registry := r

	if r.dbExecutor == nil {
		tx, err = r.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}

		defer func() {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p)
			} else if err != nil {
				rErr := tx.Rollback()
				if rErr != nil {
					err = rErr
				}
			} else {
				err = tx.Commit()
			}
		}()
		registry = RepositoryRegistryImpl{
			db:         r.db,
			dbExecutor: tx,
		}
	}
	out, err = txFunc(registry)
	return
}

type RepositoryRegistry interface {
	AccessRepository() accessRepo.AccessRepository
	AuthzRepository() authzRepo.AuthzRepository
	JWKRepository() jwkRepo.JWKRepository
	RolesRepository() rolesRepo.RolesRepository
	ResourcesRepository() resourcesRepo.ResourcesRepository
	ModulesRepository() modulesRepo.ModulesRepository
	UserCredentialsRepository() userCredentialsRepo.UserCredentialsRepository
	UsersRepository() usersRepo.UsersRepository

	DoInTransaction(ctx context.Context, txFunc InTransaction) (out interface{}, err error)
}

func NewRepositoryRegistry(db *sql.DB) *RepositoryRegistryImpl {
	return &RepositoryRegistryImpl{db: db}
}
