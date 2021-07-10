package resp

type GetDeptIdsResp struct {
	Resp

	Data []string `json:"data"`
}

type GetDeptsResp struct {
	Resp

	Data GetDeptsRespData `json:"data"`
}

type GetDeptsRespData struct {
	Depts []Dept `json:"depts"`
}

type GetRootDeptResp struct {
	Resp

	Data Dept `json:"data"`
}

type Dept struct {
	Name         string `json:"name"`
	ID           string `json:"id"`
	OpenID       string `json:"openId"`
	ParentID     string `json:"parentId"`
	ParentOpenID string `json:"parentOpenId"`
}
