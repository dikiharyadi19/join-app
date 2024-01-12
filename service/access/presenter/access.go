package presenter

import (
	"context"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/domain/service"
)

type AccessPresenterImpl struct{}

type AccessPresenter interface {
	GetAccessByRoleUID(ctx context.Context, req []*model.ReadAccessByRoleUIDResp) ([]service.GetAccessByRoleUIDResp, error)
}

func NewAccessPresenter() AccessPresenter {
	return &AccessPresenterImpl{}
}

func (ap *AccessPresenterImpl) GetAccessByRoleUID(ctx context.Context, req []*model.ReadAccessByRoleUIDResp) (resp []service.GetAccessByRoleUIDResp, err error) {
	menuMap := make(map[string]service.GetAccessByRoleUIDResp)
	indexMap := make(map[string]int)

	if len(req) > 0 {
		for i, v := range req {

			indexMap[v.UID] = i

			res := service.GetAccessByRoleUIDResp{
				UID:       v.UID,
				Name:      v.Name,
				Type:      v.Type,
				Action:    v.Action,
				ParentUID: v.ParentUID.String,
				Level:     v.Level,
			}

			if v.ParentUID.String == "" {
				resp = append(resp, res)
			} else {

				parent := menuMap[v.ParentUID.String]

				if parent.UID != "" {
					if parent.Child == nil {
						parent.Child = []service.GetAccessByRoleUIDResp{}
					}
					parent.Child = append(parent.Child, res)
				}

				resp[indexMap[v.ParentUID.String]].Child = append(resp[indexMap[v.ParentUID.String]].Child, res)

			}

			menuMap[v.UID] = res
		}
	}
	return
}
