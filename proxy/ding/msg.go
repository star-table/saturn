package ding

import (
	"github.com/galaxy-book/saturn/model/context"
	"github.com/galaxy-book/saturn/model/req"
	"github.com/galaxy-book/saturn/model/resp"
	"github.com/galaxy-book/saturn/util/json"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
)

func (d *dingProxy) SendMsg(ctx *context.Context, req req.SendMsgReq) resp.SendMsgResp {
	client := &sdk.DingTalkClient{
		AccessToken: ctx.TenantAccessToken,
		AgentId:     d.AgentId,
	}
	msg := sdk.WorkNoticeMsg{}
	json.FromJsonIgnoreError(json.ToJsonIgnoreError(req.Msg), &msg)
	msg.MsgType = req.MsgType
	for _, userId := range req.UserIds {
		_, _ = client.SendWorkNotice(&userId, nil, false, msg)
	}
	for _, openId := range req.OpenIds {
		_, _ = client.SendWorkNotice(&openId, nil, false, msg)
	}
	for _, deptId := range req.DeptIds {
		_, _ = client.SendWorkNotice(nil, &deptId, false, msg)
	}
	return resp.SendMsgResp{Resp: resp.SucResp()}
}
