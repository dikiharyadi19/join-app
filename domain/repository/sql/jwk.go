package sql

import (
	"context"
	"database/sql"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/service/jwk/repository"
	"strings"
	"time"
)

const (
	insertJWK              = "INSERT INTO jwk (id, `key`, expired_at) VALUES (?,?,?)"
	updateJWK              = "UPDATE jwk SET `key` = ?, expired_at = ? WHERE id = ?"
	selectUnexpiredKeyByID = "SELECT id, `key`, expired_at FROM jwk WHERE id = ? AND expired_at > ?"
	selectUnexpiredKey     = "SELECT id, `key`, expired_at FROM jwk WHERE expired_at > ?"
)

type JWKRepositoryImpl struct {
	db DBExecutor
}

func NewJWKRepository(db DBExecutor) repository.JWKRepository {
	return &JWKRepositoryImpl{db: db}
}

func (jr *JWKRepositoryImpl) UpsertJWK(ctx context.Context, req *model.JWK) error {

	_, err := jr.db.ExecContext(ctx, insertJWK, req.ID, req.Key, req.ExpiredAt)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			_, err = jr.db.ExecContext(ctx, updateJWK, req.Key, req.ExpiredAt, req.ID)
			if err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

func (jr *JWKRepositoryImpl) ReadUnexpiredKeyByID(ctx context.Context, req *model.ReadUnexpiredKeyByIDReq) (resp *model.ReadUnexpiredKeyByIDResp, err error) {

	resp = &model.ReadUnexpiredKeyByIDResp{}

	now := time.Now().Unix()

	err = jr.db.QueryRowContext(ctx, selectUnexpiredKeyByID, req.KeyID, now).Scan(&resp.ID, &resp.Key, &resp.ExpiredAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return resp, nil

}

func (jr *JWKRepositoryImpl) ReadUnexpiredKeys(ctx context.Context) (resp []*model.ReadUnexpiredKeyResp, err error) {

	now := time.Now().UTC()

	rows, err := jr.db.QueryContext(ctx, selectUnexpiredKey, now)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		res := &model.ReadUnexpiredKeyResp{}
		err = rows.Scan(&res.ID, &res.Key, &res.ExpiredAt)
		if err != nil {
			return nil, err
		}
		resp = append(resp, res)
	}

	return
}
