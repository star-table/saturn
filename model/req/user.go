package req

type GetUsersReq struct {
	DepartmentID string `json:"departmentId"`
	FetchChild   bool   `json:"fetchChild"`
	Limit        int    `json:"limit"`
}
