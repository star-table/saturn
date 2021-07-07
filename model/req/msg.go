package req

type SendMsgReq struct {
	OpenIds []string    `json:"openIds"`
	UserIds []string    `json:"userIds"`
	DeptIds []string    `json:"deptIds"`
	ChatIds []string    `json:"chatIds"`
	MsgType string      `json:"msgType"`
	Msg     interface{} `json:"msg"`
}
