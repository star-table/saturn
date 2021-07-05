package lark

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/util/queue"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
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
