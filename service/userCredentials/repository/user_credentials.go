package repository

import (
	"context"
	"github/yogabagas/join-app/domain/model"
)

type UserCredentialsRepository interface {
	InsertCredential(ctx context.Context, req *model.UserCredential) error
	ReadCredentialsByUserUIDAndPassword(ctx context.Context, req *model.ReadCredentialsByUserUIDAndPasswordReq) (resp *model.ReadCredentialsByUserUIDAndPasswordResp, err error)
}
