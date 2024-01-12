package presenter

import (
	"context"
	"github/yogabagas/join-app/domain/model"
	"github/yogabagas/join-app/domain/service"
)

type ResourcesPresenterImpl struct{}

type ResourcesPresenter interface {
	GetResourcesByType(ctx context.Context, req []*model.ReadResourcesByTypeResp) ([]service.GetResourcesByTypeResp, error)
}

func NewResourcesPresenter() ResourcesPresenter {
	return &ResourcesPresenterImpl{}
}

func (rp *ResourcesPresenterImpl) GetResourcesByType(ctx context.Context, req []*model.ReadResourcesByTypeResp) (resp []service.GetResourcesByTypeResp, err error) {

	menuMap := make(map[string]service.GetResourcesByTypeResp)
	indexMap := make(map[string]int)

	if len(req) > 0 {
		for i, v := range req {

			indexMap[v.UID] = i

			res := service.GetResourcesByTypeResp{
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
						parent.Child = []service.GetResourcesByTypeResp{}
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
