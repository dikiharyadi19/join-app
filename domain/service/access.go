package service

import "time"

type UpsertAccessReq struct {
	RoleUID     string   `json:"role_uid"`
	ResourceUID []string `json:"resources_uid"`
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}

type GetAccessByRoleUIDReq struct {
	RoleUID string
	Type    int
}

type GetAccessByRoleUIDResp struct {
	UID       string                   `json:"uid"`
	Name      string                   `json:"name"`
	Type      int                      `json:"type"`
	Action    string                   `json:"action"`
	ParentUID string                   `json:"parent_id,omitempty"`
	Level     int                      `json:"level"`
	Child     []GetAccessByRoleUIDResp `json:"child,omitempty"`
}
