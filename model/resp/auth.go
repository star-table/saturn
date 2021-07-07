package resp

type GetTenantAccessTokenResp struct {
	Resp

	Data GetTenantAccessTokenRespData `json:"data"`
}

type GetTenantAccessTokenRespData struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type CodeLoginResp struct {
	Resp

	Data CodeLoginRespData `json:"data"`
}

type CodeLoginRespData struct {
	UserID    string   `json:"userId"`
	UnionID   string   `json:"unionId"`
	OpenID    string   `json:"openId"`
	Name      string   `json:"name"`
	Avatar    string   `json:"avatar"`
	TenantKey string   `json:"tenantKey"`
	DeptIds   []string `json:"deptIds"`
}
