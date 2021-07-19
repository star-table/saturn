package ding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"net/url"
	"strconv"
)

func (d *dingProxy) GetTenantAccessToken(tenantKey string) resp.GetTenantAccessTokenResp {
	dingTalkSDK := &sdk.DingTalkSDK{
		SuiteKey:    d.SuiteKey,
		SuiteSecret: d.SuiteSecret,
		Token:       d.Token,
		AesKey:      d.AesKey,
		AppId:       d.AppId,
	}
	ticket, err := d.Ticket()
	if err != nil {
		return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(err)}
	}
	corp := dingTalkSDK.CreateCorp(tenantKey, ticket)
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

func (d *dingProxy) CodeLogin(tenantKey, code string) resp.CodeLoginResp {
	tenantAccessTokenResp := d.GetTenantAccessToken(tenantKey)
	if !tenantAccessTokenResp.Suc {
		return resp.CodeLoginResp{Resp: tenantAccessTokenResp.Resp}
	}
	client := &sdk.DingTalkClient{
		AccessToken: tenantAccessTokenResp.Data.Token,
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

	deptIds := make([]string, 0)
	for _, deptId := range userDetailResp.Department {
		deptIds = append(deptIds, strconv.FormatInt(deptId, 10))
	}
	respData := resp.CodeLoginRespData{
		UserID:  userInfoResp.UserId,
		UnionID: userDetailResp.UnionId,
		OpenID:  userDetailResp.UnionId,
		Name:    userDetailResp.Name,
		Avatar:  userDetailResp.Avatar,
		DeptIds: deptIds,
	}
	return resp.CodeLoginResp{
		Resp: resp.SucResp(),
		Data: respData,
	}
}

func signature(timestamp int, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(fmt.Sprintf("%v", timestamp)))
	return url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}
