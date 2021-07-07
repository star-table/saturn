package lark

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/util/json"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
)

var supportedBatchSendMsgType = map[string]bool{
	"text":       true,
	"image":      true,
	"share_chat": true,
}

func (l *larkProxy) SendMsg(ctx *context.Context, req req.SendMsgReq) resp.SendMsgResp {
	client := sdk.Tenant{
		TenantAccessToken: ctx.TenantAccessToken,
	}

	if supportedBatchSendMsgType[req.MsgType] {
		content := vo.MsgContent{}
		json.FromJsonIgnoreError(json.ToJsonIgnoreError(req.Msg), &content)
		sendBatchMsg := vo.BatchMsgVo{
			DepartmentIds: req.DeptIds,
			OpenIds:       req.OpenIds,
			UserIds:       req.UserIds,
			MsgType:       req.MsgType,
			Content:       &content,
		}
		msgSendResp, err := client.SendMessageBatch(sendBatchMsg)
		if err != nil {
			return resp.SendMsgResp{Resp: resp.ErrResp(err)}
		}
		if msgSendResp.Code != 0 {
			return resp.SendMsgResp{Resp: resp.Resp{Code: msgSendResp.Code, Msg: msgSendResp.Msg}}
		}
		for _, chatId := range req.ChatIds {
			sendMsg := vo.MsgVo{
				ChatId:  chatId,
				Content: sendBatchMsg.Content,
			}
			_, _ = client.SendMessage(sendMsg)
		}
	} else {
		card := vo.Card{}
		json.FromJsonIgnoreError(json.ToJsonIgnoreError(req.Msg), &card)
		sendInteractiveMsg(client, "openId", req.OpenIds, &card)
		sendInteractiveMsg(client, "userId", req.UserIds, &card)
		sendInteractiveMsg(client, "chatId", req.ChatIds, &card)
	}
	return resp.SendMsgResp{Resp: resp.SucResp()}
}

func sendInteractiveMsg(client sdk.Tenant, t string, ids []string, card *vo.Card) {
	msgVo := vo.MsgVo{
		MsgType: "interactive",
		Card:    card,
	}
	for _, id := range ids {
		if t == "openId" {
			msgVo.OpenId = id
		} else if t == "userId" {
			msgVo.UserId = id
		} else if t == "chatId" {
			msgVo.ChatId = id
		}
		_, _ = client.SendMessage(msgVo)
	}
}
