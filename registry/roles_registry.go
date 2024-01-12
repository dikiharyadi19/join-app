package registry

import (
	"github/yogabagas/join-app/adapter/controller"
	"github/yogabagas/join-app/service/roles/usecase"
)

func (m *module) NewRolesRegistry() usecase.RolesService {
	return usecase.NewRolesService(m.NewRepositoryRegistry())
}

func (m *module) NewRolesController() controller.RolesController {
	return controller.NewRolesController(m.NewRolesRegistry())
}
