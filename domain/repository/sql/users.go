package sql

import (
	"context"
	"database/sql"
	"errors"

	"fmt"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/service/users/repository"
	"strings"
)

const (
	insertUsers = `INSERT INTO users (uid, first_name, last_name, email, birthdate, description, gender, country, photo, created_by, updated_by) 
	VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	selectUsersByEmail = `SELECT u.uid, a.role_uid, r.name as role_name, a.last_active FROM users u JOIN authz a ON u.uid = a.user_uid 
	JOIN roles r ON a.role_uid = r.uid WHERE u.email = ? ORDER BY r.id ASC LIMIT 1`
	selectUsersWithPagination = `SELECT u.uid, u.first_name, u.last_name, u.email, u.birthdate, u.username, u.created_at, 
	(SELECT COUNT(*) from users us WHERE us.id = u.id) as per_page, r.name as role_name FROM users u JOIN authz a ON u.uid = a.user_uid 
	JOIN roles r ON a.role_uid = r.uid %s`
	selectCountUsers = `SELECT COUNT(*) FROM users WHERE is_deleted = ?`
)

type UsersRepositoryImpl struct {
	db DBExecutor
}

func NewUsersRepository(db DBExecutor) repository.UsersRepository {
	return &UsersRepositoryImpl{db: db}
}

func (ur *UsersRepositoryImpl) CreateUsers(ctx context.Context, req *model.User) error {

	_, err := ur.db.ExecContext(ctx, insertUsers, req.UID, req.FirstName, req.LastName, req.Email, req.Birthdate,
		req.Description, req.Gender, req.Country, req.Photo, req.CreatedBy, req.UpdatedBy)
	if err != nil && !strings.Contains(err.Error(), "duplicate") {
		return err
	}

	return nil
}

func (ur *UsersRepositoryImpl) ReadUserByEmail(ctx context.Context, req *model.ReadUserByEmailReq) (resp *model.ReadUserByEmailResp, err error) {
	resp = &model.ReadUserByEmailResp{}

	err = ur.db.QueryRowContext(ctx, selectUsersByEmail, req.Email).
		Scan(&resp.UserUID, &resp.RoleUID, &resp.RoleName, &resp.LastActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found, please check the credential")
		}
		return nil, err
	}

	return resp, nil
}

func (ur *UsersRepositoryImpl) ReadUsersWithPagination(ctx context.Context, req *model.ReadUsersWithPaginationReq) (resp *model.ReadUsersWithPaginationResp, err error) {

	cond := fmt.Sprintf("WHERE MATCH (u.first_name, u.last_name) AGAINST ('%s*' IN BOOLEAN MODE) LIMIT %d OFFSET %d", req.Fullname, req.Limit, req.Offset)

	if req.Fullname == "" {
		cond = fmt.Sprintf("LIMIT %d OFFSET %d", req.Limit, req.Offset)
	}

	q := fmt.Sprintf(selectUsersWithPagination, cond)

	rows, err := ur.db.QueryContext(ctx, q)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	resp = &model.ReadUsersWithPaginationResp{}

	for rows.Next() {
		user := model.UserWithRole{}
		var perPage int

		err = rows.Scan(&user.UID, &user.FirstName, &user.LastName, &user.Email, &user.Birthdate, &user.Username,
			&user.CreatedAt, &perPage, &user.RoleName)
		if err != nil {
			return nil, err
		}

		resp.PerPage += perPage
		resp.Users = append(resp.Users, user)
	}
	return resp, nil
}

func (ur *UsersRepositoryImpl) CountUsers(ctx context.Context, req *model.CountUsersReq) (resp *model.CountUsersResp, err error) {

	resp = &model.CountUsersResp{}

	err = ur.db.QueryRowContext(ctx, selectCountUsers, req.IsDeleted).Scan(&resp.Total)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return
}
