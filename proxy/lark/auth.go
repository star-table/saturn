package lark

import (
	"errors"
	"fmt"
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/util/json"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
)

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

func (l *larkProxy) CodeLogin(ctx *context.Context, code string) resp.CodeLoginResp {
	appTokenResp, err := sdk.GetAppAccessToken(l.AppId, l.Secret, l.Ticket)
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	if appTokenResp.Code != 0 {
		return resp.CodeLoginResp{Resp: resp.Resp{Code: appTokenResp.Code, Msg: appTokenResp.Msg}}
	}
	respData := resp.CodeLoginRespData{}
	loginValidateResp, err := sdk.TokenLoginValidate(appTokenResp.AppAccessToken, code)
	var respErr error
	if err != nil {
		respErr = err
	} else if loginValidateResp.Code != 0 {
		respErr = errors.New(fmt.Sprintln("TokenLoginValidate fail", json.ToJsonIgnoreError(loginValidateResp)))
	} else {
		respData.UserID = loginValidateResp.Data.Uid
		respData.OpenID = loginValidateResp.Data.OpenId
		respData.UnionID = loginValidateResp.Data.UnionId
		respData.TenantKey = loginValidateResp.Data.TenantKey
	}
	if respErr != nil {
		oauthResp, err1 := sdk.GetOauth2AccessToken(vo.OAuth2AccessTokenReqVo{
			Code:           code,
			AppId:          l.AppId,
			AppSecret:      l.Secret,
			AppAccessToken: appTokenResp.AppAccessToken,
			GrantType:      "authorization_code",
		})
		if err1 != nil {
			respErr = err
		} else if oauthResp.Code != 0 {
			respErr = errors.New(fmt.Sprintln("GetOauth2AccessToken fail", json.ToJsonIgnoreError(oauthResp)))
		} else {
			respData.UserID = oauthResp.OpenId
			respData.OpenID = oauthResp.OpenId
			respData.UnionID = oauthResp.OpenId
			respData.TenantKey = oauthResp.TenantKey
		}
	}
	if respErr != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(respErr)}
	}
	return resp.CodeLoginResp{
		Resp: resp.SucResp(),
		Data: respData,
	}
}
