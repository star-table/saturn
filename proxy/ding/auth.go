package ding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/util/http"
	"gitea.bjx.cloud/allstar/saturn/util/json"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"net/url"
	"strconv"
	"time"
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
	timestamp := time.Now().Nanosecond() / 1e6
	queries := map[string]interface{}{
		"accessKey": d.SuiteKey,
		"timestamp": timestamp,
		"signature": signature(timestamp, d.SuiteSecret),
	}
	requestBody := map[string]interface{}{
		"tmp_auth_code": code,
	}
	userInfoByCodeRespData, err := http.Post("https://oapi.dingtalk.com/sns/getuserinfo_bycode", queries, json.ToJsonIgnoreError(requestBody))
	if err != nil {
		return resp.CodeLoginResp{Resp: resp.ErrResp(err)}
	}
	userInfoByCodeResp := sdk.GetUserInfoByCodeResp{}
	json.FromJsonIgnoreError(userInfoByCodeRespData, &userInfoByCodeResp)
	if userInfoByCodeResp.ErrCode != 0 {
		return resp.CodeLoginResp{Resp: resp.Resp{Code: userInfoByCodeResp.ErrCode, Msg: userInfoByCodeResp.ErrMsg}}
	}

	respData := resp.CodeLoginRespData{
		UserID:  userInfoByCodeResp.UserInfo.OpenId,
		UnionID: userInfoByCodeResp.UserInfo.UnionId,
		OpenID:  userInfoByCodeResp.UserInfo.OpenId,
		Name:    userInfoByCodeResp.UserInfo.Nick,
	}

	if tenantKey != "" {
		tenantAccessTokenResp := d.GetTenantAccessToken(tenantKey)
		if tenantAccessTokenResp.Suc {
			client := &sdk.DingTalkClient{
				AccessToken: tenantAccessTokenResp.Data.Token,
				AgentId:     d.AgentId,
			}
			// 需要测试是否需要转换
			userIdResp, err := client.GetUserIdByUnionId(userInfoByCodeResp.UserInfo.UnionId)
			if err == nil && userIdResp.ErrCode == 0 {
				userDetailResp, err := client.GetUserDetail(userIdResp.UserId, nil)
				if err == nil && userDetailResp.ErrCode == 0 {
					deptIds := make([]string, 0)
					for _, deptId := range userDetailResp.Department {
						deptIds = append(deptIds, strconv.FormatInt(deptId, 10))
					}
					respData.Name = userDetailResp.Name
					respData.Avatar = userDetailResp.Avatar
					respData.DeptIds = deptIds
				}
			}
		}
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
