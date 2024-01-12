package service

type CreateResourcesReq struct {
	Name      string `json:"name"`
	Type      int    `json:"type"`
	Action    string `json:"action"`
	ParentUID string `json:"parent_uid"`
	CreatedBy string `json:"-"`
}

type GetResourcesByTypeReq struct {
	Type int `json:"type"`
}

type GetResourcesByTypeResp struct {
	UID       string                   `json:"uid"`
	Name      string                   `json:"name"`
	Type      int                      `json:"type"`
	Action    string                   `json:"action"`
	ParentUID string                   `json:"parent_id,omitempty"`
	Level     int                      `json:"level"`
	Child     []GetResourcesByTypeResp `json:"child,omitempty"`
}
