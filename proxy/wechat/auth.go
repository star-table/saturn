package wechat

import (
	c1 "context"
	"github.com/galaxy-book/saturn/model/resp"
	"github.com/galaxy-book/saturn/util/json"
	"github.com/galaxy-book/work-wechat"
	"strings"
)

func (w *wechatProxy) createWechatSDK(tenantKey string) (*work.WorkWechat, error) {
	corpInfos := strings.Split(tenantKey, ":")
	ticket, err := w.Ticket()
	if err != nil {
		return nil, err
	}
	wechatSDK := work.NewWorkWechat(work.Config{
		SuiteID:        w.SuiteID,
		SuiteSecret:    w.SuiteSecret,
		SuiteTicket:    ticket,
		ProviderCorpID: w.ProviderCorpID,
		ProviderSecret: w.ProviderSecret,
		CorpId:         corpInfos[0],
		PermanentCode:  corpInfos[1],
	})
	return wechatSDK, nil
}

func (w *wechatProxy) GetTenantAccessToken(tenantKey string) resp.GetTenantAccessTokenResp {
	wechatSDK, err := w.createWechatSDK(tenantKey)
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	suiteAccessTokenResp, err := wechatSDK.GetSuiteAccessToken()
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}

	corpAccessTokenResp, err := wechatSDK.GetCorpAccessToken()
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	action := work.GetCorpAuthInfoAction(suiteAccessTokenResp.SuiteAccessToken, wechatSDK.CorpId, wechatSDK.PermanentCode)
	respBody, err := action.DoRequest(c1.Background())
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	corpAuthInfoResp := work.GetCorpAuthInfoResp{}
	json.FromJsonIgnoreError(string(respBody), corpAuthInfoResp)
	if len(corpAuthInfoResp.AuthInfo.Agent) > 0 {
		w.AgentId = int64(corpAuthInfoResp.AuthInfo.Agent[0].AgentID)
	}
	return resp.GetTenantAccessTokenResp{
		Resp: resp.SucResp(),
		Data: resp.GetTenantAccessTokenRespData{
			Token:  corpAccessTokenResp.AccessToken,
			Expire: corpAccessTokenResp.ExpiresIn,
		},
	}
}

func (w *wechatProxy) CodeLogin(tenantKey, code string) resp.CodeLoginResp {
	userInfo3rdResp, userDetail3rdResp, err := w.getUserInfo3rd(code)
	if err == nil {
		return resp.CodeLoginResp{
			Resp: resp.SucResp(),
			Data: resp.CodeLoginRespData{
				UserID:    userDetail3rdResp.UserId,
				UnionID:   userInfo3rdResp.OpenUserId,
				OpenID:    userInfo3rdResp.OpenUserId,
				Name:      userDetail3rdResp.Name,
				Avatar:    userDetail3rdResp.Avatar,
				TenantKey: userInfo3rdResp.CorpId,
			},
		}
	}
	loginUserInfoResp, err := w.getLoginUserInfo(code)
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	return resp.CodeLoginResp{
		Resp: resp.SucResp(),
		Data: resp.CodeLoginRespData{
			UserID:    loginUserInfoResp.UserInfo.Userid,
			UnionID:   loginUserInfoResp.UserInfo.OpenUserid,
			OpenID:    loginUserInfoResp.UserInfo.OpenUserid,
			Name:      loginUserInfoResp.UserInfo.Name,
			Avatar:    loginUserInfoResp.UserInfo.Avatar,
			TenantKey: userInfo3rdResp.CorpId,
		},
	}
}

func (w *wechatProxy) getUserInfo3rd(code string) (*work.GetUserInfo3rdResp, *work.GetUserDetail3rdResp, error) {
	wechatSDK, err := w.createWechatSDK(":")
	if err != nil {
		return nil, nil, err
	}
	suiteAccessTokenResp, err := wechatSDK.GetSuiteAccessToken()
	if err != nil {
		return nil, nil, err
	}

	action := work.GetUserInfo3RD(suiteAccessTokenResp.SuiteAccessToken, code)
	respBody, err := action.DoRequest(c1.Background())
	if err != nil {
		return nil, nil, err
	}
	userInfo3rdResp := work.GetUserInfo3rdResp{}
	json.FromJsonIgnoreError(string(respBody), &userInfo3rdResp)

	action = work.GetUserDetail3RD(suiteAccessTokenResp.SuiteAccessToken, userInfo3rdResp.UserTicket)
	respBody, err = action.DoRequest(c1.Background())
	if err != nil {
		return nil, nil, err
	}
	userDetail3rdResp := work.GetUserDetail3rdResp{}
	json.FromJsonIgnoreError(string(respBody), &userDetail3rdResp)

	return &userInfo3rdResp, &userDetail3rdResp, nil
}

func (w *wechatProxy) getLoginUserInfo(code string) (*work.GetLoginInfoResp, error) {
	wechatSDK, err := w.createWechatSDK(":")
	if err != nil {
		return nil, err
	}
	providerAccessTokenResp, err := wechatSDK.GetProviderAccessToken()
	if err != nil {
		return nil, err
	}
	action := work.GetLoginInfo(providerAccessTokenResp.ProviderAccessToken, code)
	respBody, err := action.DoRequest(c1.Background())
	if err != nil {
		return nil, err
	}
	loginUserInfoResp := work.GetLoginInfoResp{}
	json.FromJsonIgnoreError(string(respBody), &loginUserInfoResp)
	return &loginUserInfoResp, nil
}
