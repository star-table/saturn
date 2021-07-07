package lark

type larkProxy struct {
	AppId  string
	Secret string
	Ticket string
}

func NewLarkProxy(appId, secret, ticket string) *larkProxy {
	return &larkProxy{
		AppId:  appId,
		Secret: secret,
		Ticket: ticket,
	}
}
