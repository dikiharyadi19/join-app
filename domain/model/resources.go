package model

import (
	"database/sql"
	"time"
)

type Resource struct {
	ID        int
	UID       string
	Name      string
	ParentUID string
	Type      int
	Action    string
	IsDeleted bool
	CreatedBy string
	CreatedAt time.Time
	UpdatedBy string
	UpdatedAt time.Time
}

type ReadResourcesByTypeReq struct {
	Type int
}

type ReadResourcesByTypeResp struct {
	UID       string
	Name      string
	Type      int
	Action    string
	ParentUID sql.NullString
	Level     int
}
