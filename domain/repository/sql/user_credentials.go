package sql

import (
	"context"
	"database/sql"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/service/userCredentials/repository"
)

const (
	insertUserCredentials = `INSERT INTO user_credentials (user_uid, username, password) 
	VALUES (?,?,?)`
	selectCredentialsByUserUIDAndPassword = `SELECT (1) FROM user_credentials WHERE user_uid = ? AND password = ?`
)

type UserCredentialsRepositoryImpl struct {
	db DBExecutor
}

func NewUserCredentialsRepository(db DBExecutor) repository.UserCredentialsRepository {
	return &UserCredentialsRepositoryImpl{db: db}
}

func (uc *UserCredentialsRepositoryImpl) InsertCredential(ctx context.Context, req *model.UserCredential) error {
	_, err := uc.db.ExecContext(ctx, insertUserCredentials, req.UserUID, req.Username, req.Password)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserCredentialsRepositoryImpl) ReadCredentialsByUserUIDAndPassword(ctx context.Context, req *model.ReadCredentialsByUserUIDAndPasswordReq) (resp *model.ReadCredentialsByUserUIDAndPasswordResp, err error) {

	var valid int
	err = uc.db.QueryRowContext(ctx, selectCredentialsByUserUIDAndPassword, req.UserUID, req.Password).
		Scan(&valid)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &model.ReadCredentialsByUserUIDAndPasswordResp{
		Valid: valid > 0,
	}, nil
}
