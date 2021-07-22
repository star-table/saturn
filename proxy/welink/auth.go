package welink

import (
	"gitea.bjx.cloud/allstar/welink"
	"github.com/galaxy-book/saturn/model/resp"
)

func (w *welinkProxy) getAppAccessToken() (string, error) {
	req := welink.AuthV2TicketsRequest()
	req.ClientSecret = w.ClientSecret
	req.ClientId = w.ClientID
	err := req.SetUrl("https://open.welink.huaweicloud.com/api/auth/v2/tickets")
	if err != nil {
		return "", nil
	}
	err, accessToken := req.GetResponse()
	if err != nil {
		return "", nil
	}
	return accessToken.Token, nil
}

func (w *welinkProxy) GetTenantAccessToken(tenantKey string) resp.GetTenantAccessTokenResp {
	req := welink.AuthV2TicketsRequest()
	req.ClientSecret = w.ClientSecret
	req.ClientId = w.ClientID
	req.TenantId = tenantKey
	err := req.SetUrl("https://open.welink.huaweicloud.com/api/auth/v2/tickets")
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	err, accessToken := req.GetResponse()
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	return resp.GetTenantAccessTokenResp{
		Resp: resp.SucResp(),
		Data: resp.GetTenantAccessTokenRespData{
			Token:  accessToken.Token,
			Expire: int64(accessToken.ExpiresIn),
		},
	}
}

func (w *welinkProxy) CodeLogin(tenantKey, code string) resp.CodeLoginResp {
	appAccessToken, err := w.getAppAccessToken()
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}

	req := welink.AuthV2UseridRequest()
	req.Code = code
	req.SetAccessToken(appAccessToken)
	err = req.SetUrl("https://open.welink.huaweicloud.com/api/auth/v2/userid")
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	err, response := req.GetResponse()
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}

	tenantAccessTokenResp := w.GetTenantAccessToken(response.TenantId)
	if !tenantAccessTokenResp.Suc {
		return resp.CodeLoginResp{Resp: tenantAccessTokenResp.Resp}
	}

	userInfoReq := welink.ContactV1UsersRequest()
	userInfoReq.SetAccessToken(tenantAccessTokenResp.Data.Token)
	userInfoReq.UserId = response.UserId
	err = userInfoReq.SetUrl("https://open.welink.huaweicloud.com/api/contact/v3/users/simple")
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	err, userDetailInfo := userInfoReq.GetResponse()
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	return resp.CodeLoginResp{
		Resp: resp.SucResp(),
		Data: resp.CodeLoginRespData{
			UserID:    userDetailInfo.UserId,
			UnionID:   userDetailInfo.UserId,
			OpenID:    userDetailInfo.UserId,
			Name:      userDetailInfo.UserNameCn,
			TenantKey: response.TenantId,
		},
	}
}
