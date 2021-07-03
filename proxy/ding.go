package proxy

import (
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
)

type dingProxy struct {
	AppId       int64
	SuiteKey    string
	SuiteSecret string
	Ticket      string
	Token       string
	AesKey      string
}

func NewDingProxy(appId int64, suiteKey, suiteSecret, ticket, token, aesKey string) *dingProxy {
	return &dingProxy{
		AppId:       appId,
		SuiteKey:    suiteKey,
		SuiteSecret: suiteSecret,
		Ticket:      ticket,
		Token:       token,
		AesKey:      aesKey,
	}
}

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
	return resp.GetTenantAccessTokenResp{
		Resp: resp.SucResp(),

		Token:  tokenInfo.AccessToken,
		Expire: tokenInfo.ExpiresIn,
	}
}
