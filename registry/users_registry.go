package registry

import (
	"github/yogabagas/join-app/adapter/controller"
	"github/yogabagas/join-app/service/users/presenter"
	"github/yogabagas/join-app/service/users/usecase"
)

func (m *module) NewUsersPresenter() presenter.UsersPresenter {
	return presenter.NewUsersPresenter()
}

func (m *module) NewUsersRegistry() usecase.UsersService {
	return usecase.NewUsersService(
		m.NewRepositoryRegistry(),
		m.NewCacheRegistry(),
		m.NewUsersPresenter())
}

func (m *module) NewUsersController() controller.UsersController {
	return controller.NewUsersController(m.NewUsersRegistry())
}
