package registry

import (
	"github/yogabagas/join-app/adapter/controller"
	"github/yogabagas/join-app/service/resources/presenter"
	"github/yogabagas/join-app/service/resources/usecase"
)

func (m *module) NewResourcesPresenter() presenter.ResourcesPresenter {
	return presenter.NewResourcesPresenter()
}

func (m *module) NewResourcesRegistry() usecase.ResourcesService {
	return usecase.NewResourcesService(
		m.NewCacheRegistry(),
		m.NewRepositoryRegistry(),
		m.NewResourcesPresenter())
}

func (m *module) NewResourcesController() controller.ResourcesController {
	return controller.NewResourcesController(m.NewResourcesRegistry())
}
