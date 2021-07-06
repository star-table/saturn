package req

type GetDeptIdsReq struct {
	// 上级部门id，非必填
	ParentId string `json:"parentId"`
	// 是否加载子部门
	FetchChild bool `json:"fetchChild"`
}

type GetDeptsReq struct {
	// 上级部门id，非必填
	ParentId string `json:"parentId"`
	// 是否加载子部门
	FetchChild bool `json:"fetchChild"`
}
