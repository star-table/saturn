package wechat

import "github.com/galaxy-book/saturn/proxy"

type wechatProxy struct {
	ProviderCorpID string
	ProviderSecret string
	SuiteID        string
	SuiteSecret    string
	AgentId        int64
	Ticket         proxy.Ticket
}

func NewWechatProxy(providerCorpID, providerSecret, suiteId, suiteSecret string, ticket proxy.Ticket) *wechatProxy {
	return &wechatProxy{
		ProviderSecret: providerSecret,
		ProviderCorpID: providerCorpID,
		SuiteID:        suiteId,
		SuiteSecret:    suiteSecret,
		Ticket:         ticket,
	}
}
