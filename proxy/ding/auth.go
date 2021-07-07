package ding

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
)

func (d *dingProxy) GetTenantAccessToken(tenantKey string) resp.GetTenantAccessTokenResp {
	dingTalkSDK := &sdk.DingTalkSDK{
		SuiteKey:    d.SuiteKey,
		SuiteSecret: d.SuiteSecret,
		Token:       d.Token,
		AesKey:      d.AesKey,
		AppId:       d.AppId,
	}
	corp := dingTalkSDK.CreateCorp(tenantKey, d.Ticket)
	tokenInfo, err := corp.GetCorpToken()
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	if tokenInfo.ErrCode != 0 {
		return resp.GetTenantAccessTokenResp{Resp: resp.Resp{Code: tokenInfo.ErrCode, Msg: tokenInfo.ErrMsg}}
	}
	authInfo, err := corp.GetAuthInfo()
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	if authInfo.ErrCode != 0 {
		return resp.GetTenantAccessTokenResp{Resp: resp.Resp{Code: authInfo.ErrCode, Msg: authInfo.ErrMsg}}
	}
	agents := authInfo.AuthInfo.Agent
	var targetAgent *sdk.Agent = nil
	for _, agent := range agents {
		if agent.AppId == d.AppId {
			targetAgent = &agent
			break
		}
	}
	if targetAgent == nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.Resp{Code: -1, Msg: "target agent is null"}}
	}
	d.AgentId = targetAgent.AgentId
	return resp.GetTenantAccessTokenResp{
		Resp: resp.SucResp(),
		Data: resp.GetTenantAccessTokenRespData{
			Token:  tokenInfo.AccessToken,
			Expire: tokenInfo.ExpiresIn,
		},
	}
}

func (d *dingProxy) CodeLogin(ctx *context.Context, code string) resp.CodeLoginResp {
	client := &sdk.DingTalkClient{
		AccessToken: ctx.TenantAccessToken,
		AgentId:     d.AgentId,
	}
	userInfoResp, err := client.GetUserInfoFromThird(code)
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	if userInfoResp.ErrCode != 0 {
		return resp.CodeLoginResp{Resp: resp.Resp{Code: userInfoResp.ErrCode, Msg: userInfoResp.ErrMsg}}
	}
	userDetailResp, err := client.GetUserDetail(userInfoResp.UserId, nil)
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	if userDetailResp.ErrCode != 0 {
		return resp.CodeLoginResp{Resp: resp.Resp{Code: userDetailResp.ErrCode, Msg: userDetailResp.ErrMsg}}
	}
	return resp.CodeLoginResp{
		Resp: resp.SucResp(),
		Data: resp.CodeLoginRespData{
			UserID:  userInfoResp.UserId,
			UnionID: userDetailResp.UnionId,
			OpenID:  userDetailResp.UnionId,
		},
	}
}
