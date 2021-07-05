package lark

import (
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
)

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

func (l *larkProxy) GetTenantAccessToken(tenantKey string) resp.GetTenantAccessTokenResp {
	appTokenResp, err := sdk.GetAppAccessToken(l.AppId, l.Secret, l.Ticket)
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	if appTokenResp.Code != 0 {
		return resp.GetTenantAccessTokenResp{Resp: resp.Resp{Code: appTokenResp.Code, Msg: appTokenResp.Msg}}
	}
	tenantTokenResp, err := sdk.GetTenantAccessToken(appTokenResp.AppAccessToken, tenantKey)
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	if tenantTokenResp.Code != 0 {
		return resp.GetTenantAccessTokenResp{Resp: resp.Resp{Code: tenantTokenResp.Code, Msg: tenantTokenResp.Msg}}
	}
	return resp.GetTenantAccessTokenResp{
		Resp: resp.SucResp(),
		Data: resp.GetTenantAccessTokenRespData{
			Token:  tenantTokenResp.TenantAccessToken,
			Expire: tenantTokenResp.Expire,
		},
	}
}
