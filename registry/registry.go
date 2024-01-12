package registry

import (
	"database/sql"
	"github/yogabagas/join-app/adapter/controller"
	"github/yogabagas/join-app/domain/repository/cache"
	repo "github/yogabagas/join-app/domain/repository/sql"

	"github.com/go-redis/redis/v8"
)

type module struct {
	sqlDB *sql.DB
	cache *redis.Client
	ns    string
}

type Registry interface {
	NewAppController() controller.AppController
}

type Option func(*module)

func NewRegistry(opts ...Option) Registry {
	m := &module{}

	for _, o := range opts {
		o(m)
	}

	return m
}

func NewSQLConn(db *sql.DB) Option {
	return func(m *module) {
		m.sqlDB = db
	}
}

func NewCache(redisClient *redis.Client) Option {
	return func(m *module) {
		m.cache = redisClient
	}
}

func (m *module) NewRepositoryRegistry() repo.RepositoryRegistry {
	return repo.NewRepositoryRegistry(m.sqlDB)
}

func (m *module) NewCacheRegistry() cache.Cache {
	return cache.NewCacheRepository(m.cache, m.ns)
}

func (m *module) NewAppController() controller.AppController {
	return controller.AppController{
		AccessController:    m.NewAccessController(),
		AuthzController:     m.NewAuthzController(),
		JWKController:       m.NewJWKController(),
		ResourcesController: m.NewResourcesController(),
		ModulesController:   m.NewCoursesController(),
		RolesController:     m.NewRolesController(),
		UsersController:     m.NewUsersController(),
	}
}
