package ding

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/util/queue"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"log"
	"strconv"
)

func (d *dingProxy) GetDeptIds(ctx *context.Context, req req.GetDeptIdsReq) resp.GetDeptIdsResp {
	client := &sdk.DingTalkClient{
		AccessToken: ctx.TenantAccessToken,
		AgentId:     d.AgentId,
	}
	deptIdContains := map[string]bool{}
	q := queue.New()
	q.Push("1")
	fetchChild := true
	if req.ParentId != "" {
		q.Clear()
		q.Push(req.ParentId)
		fetchChild = false
	}
	for {
		obj, err := q.Pop()
		if err != nil {
			break
		}
		parentId := obj.(string)
		subIdsResp, err := client.GetSubDept(parentId)
		if err != nil {
			log.Println(err)
			continue
		}
		if subIdsResp.ErrCode != 0 {
			log.Println(subIdsResp.ErrCode, subIdsResp.ErrMsg)
			continue
		}
		for _, subId := range subIdsResp.SubDeptIdList {
			deptStrId := strconv.FormatInt(subId, 10)
			if !deptIdContains[deptStrId] {
				deptIdContains[deptStrId] = true
				if fetchChild {
					q.Push(deptStrId)
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
