package sql

import (
	"context"
	"database/sql"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/service/resources/repository"
	"strings"
)

const (
	insertResources                = `INSERT INTO resources (uid, name, parent_uid, type, action, created_by, updated_by) VALUES (?,?,?,?,?,?,?)`
	selectResourcesHierarchyByType = `WITH RECURSIVE menu_hierarchy AS (
		SELECT uid, name, action, type, parent_uid, 1 as level FROM resources WHERE parent_uid IS NULL
		UNION ALL
		SELECT m.uid, m.name, m.action, m.type, m.parent_uid, mh.level + 1 FROM resources m
		JOIN menu_hierarchy mh ON m.parent_uid = mh.uid)
	  	SELECT uid, name, action, type, parent_uid, level
	  	FROM menu_hierarchy WHERE type = ?`
)

type ResourcesRepositoryImpl struct {
	db DBExecutor
}

func NewResourcesRepository(db DBExecutor) repository.ResourcesRepository {
	return &ResourcesRepositoryImpl{db: db}
}

func (rr *ResourcesRepositoryImpl) CreateResources(ctx context.Context, req *model.Resource) error {

	parentUID := sql.NullString{}
	if req.ParentUID != "" {
		parentUID.String = req.ParentUID
		parentUID.Valid = true
	}

	_, err := rr.db.ExecContext(ctx, insertResources, req.UID, req.Name, parentUID, req.Type, req.Action,
		req.CreatedBy, req.UpdatedBy)
	if err != nil && !strings.Contains(err.Error(), "duplicate") {
		return err
	}

	return nil
}

func (rr *ResourcesRepositoryImpl) ReadResourcesByType(ctx context.Context, req *model.ReadResourcesByTypeReq) (resp []*model.ReadResourcesByTypeResp, err error) {

	rows, err := rr.db.QueryContext(ctx, selectResourcesHierarchyByType, req.Type)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	for rows.Next() {
		res := &model.ReadResourcesByTypeResp{}

		err = rows.Scan(&res.UID, &res.Name, &res.Action, &res.Type, &res.ParentUID, &res.Level)
		if err != nil {
			return nil, err
		}
		resp = append(resp, res)
	}

	return
}
