package welink

import (
	"gitea.bjx.cloud/allstar/welink"
	"github.com/galaxy-book/saturn/model/context"
	"github.com/galaxy-book/saturn/model/req"
	"github.com/galaxy-book/saturn/model/resp"
	"strconv"
)

func (w *welinkProxy) GetDeptIds(ctx *context.Context, r req.GetDeptIdsReq) resp.GetDeptIdsResp {
	deptsResp := w.GetDepts(ctx, req.GetDeptsReq{
		ParentId:   r.ParentId,
		FetchChild: r.FetchChild,
	})
	if !deptsResp.Suc {
		return resp.GetDeptIdsResp{Resp: deptsResp.Resp}
	}
	deptIds := make([]string, 0)
	for _, dept := range deptsResp.Data.Depts {
		deptIds = append(deptIds, dept.ID)
	}
	return resp.GetDeptIdsResp{
		Resp: resp.SucResp(),
		Data: deptIds,
	}
}

func (w *welinkProxy) GetDepts(ctx *context.Context, req req.GetDeptsReq) resp.GetDeptsResp {
	parentDeptIds := make([]string, 0)
	depts := make([]resp.Dept, 0)
	if req.ParentId == "" || req.ParentId == "0" {
		primaryDepts, err := w.GetSubDepts(ctx, []string{"0"}, false)
		if err != nil {
			return resp.GetDeptsResp{Resp: resp.ErrResp(err)}
		}
		for _, dept := range primaryDepts {
			parentDeptIds = append(parentDeptIds, dept.ID)
		}
		depts = append(depts, primaryDepts...)
	} else {
		parentDeptIds = append(parentDeptIds, req.ParentId)
	}
	subDepts, err := w.GetSubDepts(ctx, parentDeptIds, true)
	if err != nil {
		return resp.GetDeptsResp{Resp: resp.ErrResp(err)}
	}
	depts = append(depts, subDepts...)
	return resp.GetDeptsResp{
		Resp: resp.SucResp(),
		Data: resp.GetDeptsRespData{
			Depts: depts,
		},
	}
}

func (w *welinkProxy) GetSubDepts(ctx *context.Context, parentIds []string, fetchChild bool) ([]resp.Dept, error) {
	deptList := welink.ContactV2DepartmentsListRequest()
	deptList.SetAccessToken(ctx.TenantAccessToken)
	err := deptList.SetUrl("https://open.welink.huaweicloud.com/api/contact/v3/departments/list")
	if err != nil {
		return nil, err
	}
	deptIdContains := map[string]bool{}
	depts := make([]resp.Dept, 0)
	for _, parentId := range parentIds {
		offset := 1
		limit := 100
		for {
			deptList.DeptCode = parentId
			deptList.Limit = strconv.Itoa(limit)
			deptList.Offset = strconv.Itoa(offset)
			if fetchChild {
				deptList.RecursiveFlag = "1"
			} else {
				deptList.RecursiveFlag = "0"
			}
			err, list := deptList.GetResponse()
			if err != nil {
				return nil, err
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
	}
	return depts, nil
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
