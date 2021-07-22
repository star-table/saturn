package lark

import (
	"errors"
	"fmt"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
	"github.com/galaxy-book/saturn/model/resp"
	"github.com/galaxy-book/saturn/util/json"
)

func (l *larkProxy) GetTenantAccessToken(tenantKey string) resp.GetTenantAccessTokenResp {
	ticket, err := l.Ticket()
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	appTokenResp, err := sdk.GetAppAccessToken(l.AppId, l.Secret, ticket)
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

func (l *larkProxy) CodeLogin(tenantKey, code string) resp.CodeLoginResp {
	ticket, err := l.Ticket()
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	appTokenResp, err := sdk.GetAppAccessToken(l.AppId, l.Secret, ticket)
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	if appTokenResp.Code != 0 {
		return resp.CodeLoginResp{Resp: resp.Resp{Code: appTokenResp.Code, Msg: appTokenResp.Msg}}
	}
	respData := resp.CodeLoginRespData{}
	oauthResp, err1 := sdk.GetOauth2AccessToken(vo.OAuth2AccessTokenReqVo{
		Code:           code,
		AppId:          l.AppId,
		AppSecret:      l.Secret,
		AppAccessToken: appTokenResp.AppAccessToken,
		GrantType:      "authorization_code",
	})
	var respErr error
	var accessToken = ""
	if err1 != nil {
		respErr = err
	} else if oauthResp.Code != 0 {
		respErr = errors.New(fmt.Sprintln("GetOauth2AccessToken fail", json.ToJsonIgnoreError(oauthResp)))
	} else {
		respData.UserID = oauthResp.OpenId
		respData.OpenID = oauthResp.OpenId
		respData.UnionID = oauthResp.OpenId
		respData.TenantKey = oauthResp.TenantKey
		respData.Avatar = oauthResp.AvatarUrl
		respData.Name = oauthResp.Name
		accessToken = oauthResp.AccessToken
	}
	if respErr != nil {
		loginValidateResp, err3 := sdk.TokenLoginValidate(appTokenResp.AppAccessToken, code)
		if err3 != nil {
			respErr = err3
		} else if loginValidateResp.Code != 0 {
			respErr = errors.New(fmt.Sprintln("TokenLoginValidate fail", json.ToJsonIgnoreError(loginValidateResp)))
		} else {
			respData.UserID = loginValidateResp.Data.Uid
			respData.OpenID = loginValidateResp.Data.OpenId
			respData.UnionID = loginValidateResp.Data.UnionId
			respData.TenantKey = loginValidateResp.Data.TenantKey
			accessToken = loginValidateResp.Data.AccessToken
		}
	}
	if respErr != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(respErr)}
	}

	// 获取当前登录用户头像和名字
	userInfo, err2 := sdk.GetOAuth2UserInfo(accessToken)
	if err2 != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err2)}
	}
	respData.Name = userInfo.Name
	respData.Avatar = userInfo.AvatarUrl

	// 获取当前登录用户部门, 这个接口不需要通讯录权限
	tenantAccessTokenResp := l.GetTenantAccessToken(respData.TenantKey)
	if !tenantAccessTokenResp.Suc {
		return resp.CodeLoginResp{Resp: resp.Resp{Code: tenantAccessTokenResp.Code, Msg: tenantAccessTokenResp.Msg}}
	}

	client := sdk.Tenant{
		TenantAccessToken: tenantAccessTokenResp.Data.Token,
	}
	userBatchResp, err := client.GetUserBatchGetV2(nil, []string{respData.OpenID})
	if err == nil && userBatchResp.Code == 0 && len(userBatchResp.Data.Users) > 0 {
		respData.DeptIds = userBatchResp.Data.Users[0].Departments
	}

	return resp.CodeLoginResp{
		Resp: resp.SucResp(),
		Data: respData,
	}
}
