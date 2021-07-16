package ding

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/util/json"
	"github.com/galaxy-book/work-wechat"
	"strings"
)

func (w *wechatProxy) SendMsg(ctx *context.Context, req req.SendMsgReq) resp.SendMsgResp {
	msg := work.SendMsgReq{}
	json.FromJsonIgnoreError(json.ToJsonIgnoreError(req.Msg), &msg)
	msg.MsgType = req.MsgType

	if len(req.OpenIds) > 0 {
		msg.ToUser = strings.Join(req.OpenIds, "|")
	}
	if len(req.DeptIds) > 0 {
		msg.ToParty = strings.Join(req.DeptIds, "|")
	}
	if len(req.ChatIds) > 0 {
		msg.ToTag = strings.Join(req.ChatIds, "|")
	}
	action := work.SendMsg(ctx.TenantAccessToken, msg)
	respBody, err := action.GetRequestBody()
	if err != nil {
		return resp.SendMsgResp{Resp: resp.ErrResp(err)}
	}
	sendMsgResp := work.SendMsgResp{}
	json.FromJsonIgnoreError(string(respBody), &sendMsgResp)
	if sendMsgResp.ErrCode != 0 {
		return resp.SendMsgResp{Resp: resp.Resp{Code: sendMsgResp.ErrCode, Msg: sendMsgResp.ErrMsg}}
	}
	return resp.SendMsgResp{Resp: resp.SucResp()}
}
