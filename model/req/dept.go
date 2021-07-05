package req

type GetDeptIdsReq struct {
	// 上级部门id，非必填
	ParentId string `json:"parentId"`
}
