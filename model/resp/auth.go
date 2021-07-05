package resp

type GetTenantAccessTokenResp struct {
	Resp

	Data GetTenantAccessTokenRespData `json:"data"`
}

type GetTenantAccessTokenRespData struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}
