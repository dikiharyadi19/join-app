package controller

type AppController struct {
	UsersController     interface{ UsersController }
	ResourcesController interface{ ResourcesController }
	RolesController     interface{ RolesController }
	ModulesController   interface{ ModulesController }
}
