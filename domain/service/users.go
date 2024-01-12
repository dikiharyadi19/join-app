package service

type CreateUsersReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Birthdate string `json:"birthdate"`
	Gender    int    `json:"gender"`
	Country   string `json:"country"`
	Email     string `json:"email"`
	Photo     string `json:"photo"`
	Bio       string `json:"bio"`
	RoleID    int    `json:"role_id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedBy string `json:"-"`
}

type LogoutReq struct {
	UserUID string `json:"user_uid"`
}

type GetUsersWithPaginationReq struct {
	Fullname string `json:"name"`
	Limit    int    `json:"limit"`
	Page     int    `json:"page"`
}

type GetUsersWithPaginationResp struct {
	Users      []UserResp `json:"users"`
	Pagination Pagination `json:"pagination"`
}

type UserResp struct {
	Fullname  string `json:"name"`
	Username  string `json:"username"`
	Birthdate string `json:"birthdate"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type Pagination struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}
