package lark

import (
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
	"github.com/galaxy-book/saturn/model/context"
	"github.com/galaxy-book/saturn/model/req"
	"github.com/galaxy-book/saturn/model/resp"
	"github.com/galaxy-book/saturn/util/queue"
	"log"
)

func (l *larkProxy) GetDeptIds(ctx *context.Context, req req.GetDeptIdsReq) resp.GetDeptIdsResp {
	client := sdk.Tenant{
		TenantAccessToken: ctx.TenantAccessToken,
	}

	deptIdContains := map[string]bool{}
	q := queue.New()
	q.Push("0")
	if req.ParentId != "" {
		q.Clear()
		q.Push(req.ParentId)
	}
	for {
		obj, err := q.Pop()
		if err != nil {
			break
		}
		parentId := obj.(string)

		hasMore := true
		pageToken := ""
		for hasMore {
			deptSimpleInfoResp, err := client.GetDepartmentSimpleListV2(parentId, pageToken, 100, false)
			if err != nil {
				log.Println(err)
				break
			}
			if deptSimpleInfoResp.Code != 0 {
				log.Println(deptSimpleInfoResp.Code, deptSimpleInfoResp.Msg)
				break
			}
			hasMore = deptSimpleInfoResp.Data.HasMore
			pageToken = deptSimpleInfoResp.Data.PageToken
			for _, deptInfo := range deptSimpleInfoResp.Data.DepartmentInfos {
				if !deptIdContains[deptInfo.Id] {
					deptIdContains[deptInfo.Id] = true
					if req.FetchChild {
						q.Push(deptInfo.Id)
					}
				}
			}
		}
	}
	deptIds := make([]string, 0)
	for k, _ := range deptIdContains {
		deptIds = append(deptIds, k)
	}
	return resp.GetDeptIdsResp{
		Resp: resp.SucResp(),
		Data: deptIds,
	}
}

func (l *larkProxy) GetDepts(ctx *context.Context, req req.GetDeptsReq) resp.GetDeptsResp {
	client := sdk.Tenant{
		TenantAccessToken: ctx.TenantAccessToken,
	}

	deptIdContains := map[string]bool{}
	depts := make([]resp.Dept, 0)
	q := queue.New()
	q.Push("0")
	if req.ParentId != "" {
		q.Clear()
		q.Push(req.ParentId)
	}
	for {
		obj, err := q.Pop()
		if err != nil {
			break
		}
		parentId := obj.(string)

		hasMore := true
		pageToken := ""
		for hasMore {
			deptSimpleInfoResp, err := client.GetDepartmentSimpleListV2(parentId, pageToken, 100, false)
			if err != nil {
				log.Println(err)
				break
			}
			if deptSimpleInfoResp.Code != 0 {
				log.Println(deptSimpleInfoResp.Code, deptSimpleInfoResp.Msg)
				break
			}
			hasMore = deptSimpleInfoResp.Data.HasMore
			pageToken = deptSimpleInfoResp.Data.PageToken
			for _, deptInfo := range deptSimpleInfoResp.Data.DepartmentInfos {
				if !deptIdContains[deptInfo.Id] {
					deptIdContains[deptInfo.Id] = true
					depts = append(depts, resp.Dept{
						Name:         deptInfo.Name,
						ID:           deptInfo.Id,
						ParentID:     deptInfo.ParentId,
						OpenID:       deptInfo.OpenDepartmentID,
						ParentOpenID: deptInfo.ParentOpenDepartmentID,
					})
					if req.FetchChild {
						q.Push(deptInfo.Id)
					}
				}
			}
		}
	}
	deptIds := make([]string, 0)
	for k, _ := range deptIdContains {
		deptIds = append(deptIds, k)
	}
	return resp.GetDeptsResp{
		Resp: resp.SucResp(),
		Data: resp.GetDeptsRespData{
			Depts: depts,
		},
	}
}

func (l *larkProxy) GetRootDept(ctx *context.Context) resp.GetRootDeptResp {
	return resp.GetRootDeptResp{
		Resp: resp.SucResp(),
		Data: resp.Dept{
			Name:         "飞书企业",
			ID:           "0",
			OpenID:       "0",
			ParentID:     "0",
			ParentOpenID: "0",
		},
	}
}
