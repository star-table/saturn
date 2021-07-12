package resp

type GetUsersResp struct {
	Resp

	Data GetUsersRespData `json:"data"`
}

type GetUsersRespData struct {
	Users []User `json:"users"`
}

type GetUserResp struct {
	Resp

	Data User `json:"data"`
}

type User struct {
	OpenID        string   `json:"openId"`
	UserID        string   `json:"userId"`
	UnionID       string   `json:"unionId"`
	Name          string   `json:"name"`
	EnName        string   `json:"enName"`
	Email         string   `json:"email"`
	Mobile        string   `json:"mobile"`
	Avatar        Avatar   `json:"avatar"`
	IsAdmin       bool     `json:"isAdmin"`
	DepartmentIds []string `json:"departmentIds"`
}

type Avatar struct {
	Avatar72     string `json:"avatar72"`
	Avatar240    string `json:"avatar240"`
	Avatar640    string `json:"avatar640"`
	AvatarOrigin string `json:"avatarOrigin"`
}
