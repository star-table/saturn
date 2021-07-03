package resp

type Resp struct {
	Code int    `json:"code"` // proxy resp code
	Suc  bool   `json:"suc"`
	Msg  string `json:"msg"`
}

func ErrResp(err error) Resp {
	return Resp{
		Code: -1,
		Suc:  false,
		Msg:  err.Error(),
	}
}

func SucResp() Resp {
	return Resp{
		Code: 0,
		Suc:  true,
		Msg:  "Success",
	}
}

type GetAppAccessTokenResp struct {
	Resp

	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type GetTenantAccessTokenResp struct {
	Resp

	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}
