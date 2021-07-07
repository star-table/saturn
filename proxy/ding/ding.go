package ding

import "gitea.bjx.cloud/allstar/saturn/proxy"

type dingProxy struct {
	AppId       int64
	SuiteKey    string
	SuiteSecret string
	Token       string
	AesKey      string
	AgentId     int64
	Ticket      proxy.Ticket
}

func NewDingProxy(appId int64, suiteKey, suiteSecret, token, aesKey string, ticket proxy.Ticket) *dingProxy {
	return &dingProxy{
		AppId:       appId,
		SuiteKey:    suiteKey,
		SuiteSecret: suiteSecret,
		Ticket:      ticket,
		Token:       token,
		AesKey:      aesKey,
	}
}
