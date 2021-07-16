package ding

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/util/json"
	"github.com/LLLjjjjjj/work-wechat"
	"strconv"
)

func (w *wechatProxy) GetDeptIds(ctx *context.Context, req req.GetDeptIdsReq) resp.GetDeptIdsResp {
	action := work.GetDeptList(ctx.TenantAccessToken, req.ParentId)
	respBody, err := action.GetRequestBody()
	if err != nil {
		return resp.GetDeptIdsResp{Resp: resp.ErrResp(err)}
	}
	deptListResp := work.GetDeptListResp{}
	json.FromJsonIgnoreError(string(respBody), &deptListResp)
	if deptListResp.ErrCode != 0 {
		return resp.GetDeptIdsResp{Resp: resp.Resp{Code: deptListResp.ErrCode, Msg: deptListResp.ErrMsg}}
	}
	parentId := 0
	if req.ParentId != "" {
		parentId, _ = strconv.Atoi(req.ParentId)
	}
	deptIds := make([]string, 0)
	for _, dept := range deptListResp.Department {
		if dept.ID != parentId {
			deptIds = append(deptIds, strconv.Itoa(dept.ID))
		}
	}
	return resp.GetDeptIdsResp{
		Resp: resp.SucResp(),
		Data: deptIds,
	}
}

func (w *wechatProxy) GetDepts(ctx *context.Context, req req.GetDeptsReq) resp.GetDeptsResp {
	action := work.GetDeptList(ctx.TenantAccessToken, req.ParentId)
	respBody, err := action.GetRequestBody()
	if err != nil {
		return resp.GetDeptsResp{Resp: resp.ErrResp(err)}
	}
	deptListResp := work.GetDeptListResp{}
	json.FromJsonIgnoreError(string(respBody), &deptListResp)
	if deptListResp.ErrCode != 0 {
		return resp.GetDeptsResp{Resp: resp.Resp{Code: deptListResp.ErrCode, Msg: deptListResp.ErrMsg}}
	}
	parentId := 0
	if req.ParentId != "" {
		parentId, _ = strconv.Atoi(req.ParentId)
	}
	depts := make([]resp.Dept, 0)
	for _, dept := range deptListResp.Department {
		if dept.ID != parentId {
			deptId := strconv.Itoa(dept.ID)
			deptParentId := strconv.Itoa(dept.ParentId)
			depts = append(depts, resp.Dept{
				Name:         dept.Name,
				ID:           deptId,
				OpenID:       deptId,
				ParentID:     deptParentId,
				ParentOpenID: deptParentId,
			})
		}
	}
	return resp.GetDeptsResp{
		Resp: resp.SucResp(),
		Data: resp.GetDeptsRespData{
			Depts: depts,
		},
	}
}

func (w *wechatProxy) GetRootDept(ctx *context.Context) resp.GetRootDeptResp {
	return resp.GetRootDeptResp{
		Resp: resp.SucResp(),
		Data: resp.Dept{
			Name:         "微信企业",
			ID:           "0",
			OpenID:       "0",
			ParentID:     "0",
			ParentOpenID: "0",
		},
	}
}
