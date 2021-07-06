package ding

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
	"log"
)

var supportedBatchSendMsgType = map[string]bool{
	"text":       true,
	"image":      true,
	"share_chat": true,
}

func (d *dingProxy) SendMsg(ctx *context.Context, req req.SendMsgReq) resp.SendMsgResp {
	client := sdk.Tenant{
		TenantAccessToken: ctx.TenantAccessToken,
	}

	if supportedBatchSendMsgType[req.MsgType] {
		sendBatchMsg := vo.BatchMsgVo{
			DepartmentIds: req.DeptIds,
			OpenIds:       req.OpenIds,
			UserIds:       req.UserIds,
			MsgType:       req.MsgType,
			Content: &vo.MsgContent{
				Text:        req.Msg.Text,
				ImageKey:    req.Msg.ImageKey,
				ShareChatId: req.Msg.ShareChatID,
			},
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
		card := &vo.Card{
			Config: &vo.CardConfig{
				EnableForward: true,
			},
			CardLink: &vo.CardElementUrl{
				Url:        req.Msg.Url.Url,
				AndroidUrl: req.Msg.Url.Android,
				IosUrl:     req.Msg.Url.IOS,
				PcUrl:      req.Msg.Url.PC,
			},
			Header: &vo.CardHeader{
				Title: &vo.CardHeaderTitle{
					Tag:     "plain_text",
					Content: req.Msg.Title,
				},
				Template: req.Msg.Color,
			},
		}
		elements := make([]interface{}, 0)
		for _, e := range req.Msg.Elements {
			ele := convertElements(e)
			if ele != nil {
				elements = append(elements, ele)
			}
		}
		sendInteractiveMsg(client, "openId", req.OpenIds, card)
		sendInteractiveMsg(client, "userId", req.UserIds, card)
		sendInteractiveMsg(client, "chatId", req.ChatIds, card)
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

func convertElements(e req.Element) interface{} {
	if e.Type == "div" {
		fields := make([]vo.CardElementField, 0)
		for _, f := range e.ContentFields {
			fields = append(fields, vo.CardElementField{
				IsShort: f.IsShort,
				Text: vo.CardElementText{
					Tag:     "lark_md",
					Content: f.Text,
				},
			})
		}
		return vo.CardElementContentModule{
			Tag: "div",
			Text: &vo.CardElementText{
				Tag:     "lark_md",
				Content: e.ContentText,
			},
			Fields: fields,
		}
	} else if e.Type == "hr" {
		return vo.CardElementBrModule{
			Tag: "hr",
		}
	} else if e.Type == "img" {
		return vo.CardElementImageModule{
			Tag:    "img",
			ImgKey: e.ImgKey,
			Alt: vo.CardElementText{
				Tag:     "lark_md",
				Content: e.ImgAlt,
			},
			Title: &vo.CardElementText{
				Tag:     "lark_md",
				Content: e.ImgTitle,
			},
			CustomWidth: e.ImgWidth,
			Mode:        e.ImgMode,
			Preview:     e.ImgPreview,
		}
	} else if e.Type == "action" {
		actions := make([]interface{}, 0)
		for _, a := range e.Actions {
			ele := convertElements(a)
			if ele != nil {
				actions = append(actions, ele)
			}
		}
		return vo.CardElementActionModule{
			Tag:     "action",
			Layout:  e.ActionLayout,
			Actions: actions,
		}
	} else if e.Type == "button" {
		return vo.ActionButton{
			Tag: "button",
			Text: vo.CardElementText{
				Tag:     "lark_md",
				Content: e.ButtonText,
			},
			Url: e.ButtonUrl.Url,
			MultiUrl: &vo.CardElementUrl{
				Url:        e.ButtonUrl.Url,
				AndroidUrl: e.ButtonUrl.Android,
				IosUrl:     e.ButtonUrl.IOS,
				PcUrl:      e.ButtonUrl.PC,
			},
			Type:  e.ButtonType,
			Value: e.ButtonValue,
			Confirm: &vo.CardElementConfirm{
				Title: &vo.CardHeaderTitle{
					Tag:     "plain_text",
					Content: e.ButtonConfirm.Title,
				},
				Text: &vo.CardElementText{
					Tag:     "lark_md",
					Content: e.ButtonConfirm.Text,
				},
			},
		}
	} else if e.Type == "select_static" || e.Type == "select_person" {
		options := make([]vo.CardElementOption, 0)
		for _, o := range e.SelectOptions {
			options = append(options, vo.CardElementOption{
				Text: &vo.CardElementText{
					Tag:     "lark_md",
					Content: o.Text,
				},
				Value: o.Value,
			})
		}
		return vo.ActionSelectMenu{
			Tag: e.Type,
			Placeholder: &vo.CardElementText{
				Tag:     "lark_md",
				Content: e.SelectPlaceholder,
			},
			InitialOption: e.SelectDefaultValue,
			Options:       options,
			Value:         e.SelectValue,
			Confirm: &vo.CardElementConfirm{
				Title: &vo.CardHeaderTitle{
					Tag:     "lark_md",
					Content: e.SelectConfirm.Title,
				},
				Text: &vo.CardElementText{
					Tag:     "lark_md",
					Content: e.SelectConfirm.Text,
				},
			},
		}
	} else if e.Type == "picker_datetime" {
		return vo.ActionDatePicker{
			Tag:             e.Type,
			InitialDatetime: e.DateDefaultValue,
			Placeholder: &vo.CardElementText{
				Tag:     "lark_md",
				Content: e.DatePlaceholder,
			},
			Value: e.DateValue,
			Confirm: &vo.CardElementConfirm{
				Title: &vo.CardHeaderTitle{
					Tag:     "lark_md",
					Content: e.DateConfirm.Title,
				},
				Text: &vo.CardElementText{
					Tag:     "lark_md",
					Content: e.DateConfirm.Text,
				},
			},
		}
	} else if e.Type == "note" {
		elements := make([]interface{}, 0)
		for _, e := range e.Notes {
			ele := convertElements(e)
			if ele != nil {
				elements = append(elements, ele)
			}
		}
		return vo.CardElementNote{
			Tag:      e.Type,
			Elements: elements,
		}
	}
	log.Println("Warning, not supported type ", e.Type)
	return nil
}
