package lark

import "gitea.bjx.cloud/allstar/saturn/proxy"

type larkProxy struct {
	AppId  string
	Secret string
	Ticket proxy.Ticket
}

func NewLarkProxy(appId, secret string, ticket proxy.Ticket) *larkProxy {
	return &larkProxy{
		AppId:  appId,
		Secret: secret,
		Ticket: ticket,
	}
}
