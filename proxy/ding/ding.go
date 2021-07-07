package ding

type dingProxy struct {
	AppId       int64
	SuiteKey    string
	SuiteSecret string
	Ticket      string
	Token       string
	AesKey      string
	AgentId     int64
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
