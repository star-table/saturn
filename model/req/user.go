package req

type GetUsersReq struct {
	DepartmentID string `json:"departmentId"`
	PageToken    string `json:"pageToken"`
	PageSize     int    `json:"pageSize"`
}
