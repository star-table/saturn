package welink

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/util/queue"
	"gitea.bjx.cloud/allstar/welink"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
	"log"
	"strconv"
)

func (w *welinkProxy) GetDeptIds(ctx *context.Context, req req.GetDeptIdsReq) resp.GetDeptIdsResp {
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

func (w *welinkProxy) GetDepts(ctx *context.Context, req req.GetDeptsReq) resp.GetDeptsResp {
	deptList := welink.ContactV2DepartmentsListRequest()
	deptList.SetAccessToken(ctx.TenantAccessToken)
	err := deptList.SetUrl("https://open.welink.huaweicloud.com/api/contact/v3/departments/list")
	if err != nil {
		return resp.GetDeptsResp{Resp: resp.ErrResp(err)}
	}
	if req.ParentId == "" {
		deptList.DeptCode = "0"
	} else {
		deptList.DeptCode = req.ParentId
	}
	if req.FetchChild {
		deptList.RecursiveFlag = "1"
	} else {
		deptList.RecursiveFlag = "0"
	}

	depts := make([]resp.Dept, 0)
	deptIdContains := map[string]bool{}
	offset := 1
	limit := 100
	for {
		deptList.Limit = strconv.Itoa(limit)
		deptList.Offset = strconv.Itoa(offset)
		err, list := deptList.GetResponse()
		if err != nil {
			return resp.GetDeptsResp{Resp: resp.ErrResp(err)}
		}
		for _, welinkDept := range list.DepartmentInfo {
			if !deptIdContains[welinkDept.DeptCode] {
				deptIdContains[welinkDept.DeptCode] = true
				depts = append(depts, resp.Dept{
					Name:         welinkDept.DeptNameCn,
					ID:           welinkDept.DeptCode,
					ParentID:     welinkDept.FatherCode,
					OpenID:       welinkDept.DeptCode,
					ParentOpenID: welinkDept.FatherCode,
				})
			}
		}

		if len(list.DepartmentInfo) < limit {
			break
		}
		offset++
	}
	return resp.GetDeptsResp{
		Resp: resp.SucResp(),
		Data: resp.GetDeptsRespData{
			Depts: depts,
		},
	}
}

func (w *welinkProxy) GetRootDept(ctx *context.Context) resp.GetRootDeptResp {
	return resp.GetRootDeptResp{
		Resp: resp.SucResp(),
		Data: resp.Dept{
			Name:         "welink企业",
			ID:           "0",
			OpenID:       "0",
			ParentID:     "0",
			ParentOpenID: "0",
		},
	}
}
