package controller

type AppController struct {
	AccessController    interface{ AccessController }
	AuthzController     interface{ AuthzController }
	JWKController       interface{ JWKController }
	UsersController     interface{ UsersController }
	ResourcesController interface{ ResourcesController }
	RolesController     interface{ RolesController }
	ModulesController   interface{ ModulesController }
}
