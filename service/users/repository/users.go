package repository

import (
	"context"
	"github/yogabagas/join-app/domain/model"
)

type UsersRepository interface {
	CreateUsers(ctx context.Context, req *model.User) error
	ReadUserByEmail(ctx context.Context, req *model.ReadUserByEmailReq) (*model.ReadUserByEmailResp, error)
	ReadUsersWithPagination(ctx context.Context, req *model.ReadUsersWithPaginationReq) (*model.ReadUsersWithPaginationResp, error)
	CountUsers(ctx context.Context, req *model.CountUsersReq) (*model.CountUsersResp, error)
}
