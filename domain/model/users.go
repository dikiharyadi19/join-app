package model

import (
	"time"

	"github.com/go-jose/go-jose/v3"
)

type User struct {
	ID          int
	UID         string
	FirstName   string
	LastName    string
	Email       string
	Birthdate   time.Time
	Description string
	Gender      int
	Country     string
	Photo       string
	IsDeleted   bool
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}

type ReadUsersWithPaginationReq struct {
	Fullname string
	Limit    int
	Offset   int
}

type ReadUsersWithPaginationResp struct {
	Users   []UserWithRole
	PerPage int
}

type UserWithRole struct {
	ID        int
	UID       string
	FirstName string
	LastName  string
	Email     string
	Birthdate time.Time
	Username  string
	Password  string
	IsDeleted bool
	CreatedBy string
	CreatedAt time.Time
	UpdatedBy string
	UpdatedAt time.Time
	RoleName  string
}

type CountUsersReq struct {
	IsDeleted int
}

type CountUsersResp struct {
	Total int
}

type GenerateAccessTokenReq struct {
	KeyID      string
	UserUID    string
	RoleUID    string
	LastActive int64
	ExpiredAt  int
	Signer     jose.Signer
}

type GenerateAccessTokenResp struct {
	Token string
}

type GenerateRefreshTokenReq struct {
	KeyID     string
	UserUID   string
	ExpiredAt int
	Signer    jose.Signer
}

type GenerateRefreshTokenResp struct {
	Token string
}

type ReadUserByEmailReq struct {
	Email string
}

type ReadUserByEmailResp struct {
	UserUID    string
	RoleUID    string
	RoleName   string
	LastActive time.Time
}
