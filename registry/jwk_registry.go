package registry

import (
	"github/yogabagas/join-app/adapter/controller"
	"github/yogabagas/join-app/service/jwk/presenter"
	"github/yogabagas/join-app/service/jwk/usecase"
)

func (m *module) NewJWKPresenter() presenter.JWKPresenter {
	return presenter.NewJWKPresenter()
}

func (m *module) NewJWKRegistry() usecase.JWKService {
	return usecase.NewJWKService(
		m.NewRepositoryRegistry(),
		m.NewCacheRegistry(),
		m.NewJWKPresenter(),
	)
}

func (m *module) NewJWKController() controller.JWKController {
	return controller.NewJWKController(m.NewJWKRegistry())
}
