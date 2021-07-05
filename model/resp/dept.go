package resp

type GetDeptIdsResp struct {
	Resp

	Data []string `json:"data"`
}
