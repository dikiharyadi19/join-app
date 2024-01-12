package model

import (
	"database/sql"
	"time"
)

type Access struct {
	ID          int
	UID         string
	RoleUID     string
	ResourceUID string
	IsDeleted   bool
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}

type ReadAccessByRoleUIDReq struct {
	RoleUID string
	Type    int
}

type ReadAccessByRoleUIDResp struct {
	UID       string
	RoleUID   string
	Name      string
	Type      int
	Action    string
	ParentUID sql.NullString
	Level     int
}
