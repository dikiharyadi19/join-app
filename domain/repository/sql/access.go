package sql

import (
	"context"
	"database/sql"
	"fmt"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/service/access/repository"
	"strings"
)

const (
	insertAccess             = `INSERT INTO access (uid, role_uid, resource_uid, created_by, updated_by) VALUES %s ON DUPLICATE KEY UPDATE is_deleted = false`
	updateAccess             = `UPDATE access SET is_deleted = TRUE WHERE role_uid = ? AND resource_uid NOT IN (?)`
	selectResourcesByRoleUID = `WITH RECURSIVE menu_hierarchy AS (
		SELECT uid, name, action, type, parent_uid, 1 as level FROM resources WHERE parent_uid IS NULL
		UNION ALL
		SELECT m.uid, m.name, m.action, m.type, m.parent_uid, mh.level + 1 FROM resources m
		JOIN menu_hierarchy mh ON m.parent_uid = mh.uid)
	  	SELECT mh.uid, a.role_uid, mh.name, mh.action, mh.type, mh.parent_uid, mh.level
	  	FROM menu_hierarchy mh JOIN access a ON mh.uid = a.resource_uid WHERE role_uid = ? AND type = ?`
)

type AccessRepositoryImpl struct {
	db DBExecutor
}

func NewAccessRepository(db DBExecutor) repository.AccessRepository {
	return &AccessRepositoryImpl{db: db}
}

func (ar *AccessRepositoryImpl) UpsertAccess(ctx context.Context, req []*model.Access) error {
	var (
		values       []string
		args         []interface{}
		roleUID      string
		resourcesUID []interface{}
	)

	if len(req) > 0 {
		for _, v := range req {
			values = append(values, "(?,?,?,?,?)")
			args = append(args, v.UID, v.RoleUID, v.ResourceUID, v.CreatedBy, v.UpdatedBy)
			roleUID = v.RoleUID
			resourcesUID = append(resourcesUID, v.ResourceUID)
		}

		q := fmt.Sprintf(insertAccess, strings.Join(values, ", "))

		_, err := ar.db.ExecContext(ctx, q, args...)
		if err != nil {
			return err
		}

		placeHolder := strings.Join(strings.Split(strings.Repeat("?", len(req)), ""), ", ")
		updateQuery := strings.Replace(updateAccess, "(?)", "("+placeHolder+")", 1)

		_, err = ar.db.ExecContext(ctx, updateQuery, append([]interface{}{roleUID}, resourcesUID...)...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ar *AccessRepositoryImpl) ReadAccessByRoleUID(ctx context.Context, req *model.ReadAccessByRoleUIDReq) (resp []*model.ReadAccessByRoleUIDResp, err error) {

	rows, err := ar.db.QueryContext(ctx, selectResourcesByRoleUID, req.RoleUID, req.Type)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	for rows.Next() {
		res := &model.ReadAccessByRoleUIDResp{}

		err = rows.Scan(&res.UID, &res.RoleUID, &res.Name, &res.Action, &res.Type, &res.ParentUID, &res.Level)
		if err != nil {
			return nil, err
		}
		resp = append(resp, res)
	}

	return
}
