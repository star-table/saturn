package ding

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/util/json"
	"github.com/LLLjjjjjj/work-wechat"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"log"
	"strconv"
)

func (w *wechatProxy) GetUsers(ctx *context.Context, r req.GetUsersReq) resp.GetUsersResp {
	if r.DepartmentID == "" {
		r.DepartmentID = "0"
	}
	fetchChild := 0
	if r.FetchChild {
		fetchChild = 1
	}
	action := work.GetDeptMemberList(ctx.TenantAccessToken, r.DepartmentID, fetchChild)
	respBody, err := action.GetRequestBody()
	if err != nil {
		return resp.GetUsersResp{Resp: resp.ErrResp(err)}
	}
	deptMemberListResp := work.GetDeptMemberListResp{}
	json.FromJsonIgnoreError(string(respBody), &deptMemberListResp)

	respUsers := make([]resp.User, 0)
	for _, deptMember := range deptMemberListResp.UserList {
		deptIdList := make([]string, 0)
		for _, deptId := range deptMember.Department {
			deptIdList = append(deptIdList, strconv.Itoa(deptId))
		}
		respUsers = append(respUsers, resp.User{
			OpenID:  deptMember.OpenUserid,
			UserID:  deptMember.Userid,
			UnionID: deptMember.OpenUserid,
			Name:    deptMember.Name,
			EnName:  deptMember.Name,
			Email:   deptMember.Email,
			Mobile:  deptMember.Mobile,
			IsAdmin: false,
			Avatar: resp.Avatar{
				Avatar72:     deptMember.Avatar,
				Avatar240:    deptMember.Avatar,
				Avatar640:    deptMember.Avatar,
				AvatarOrigin: deptMember.Avatar,
			},
			DepartmentIds: deptIdList,
		})
	}
	return resp.GetUsersResp{
		Resp: resp.SucResp(),
		Data: resp.GetUsersRespData{
			Users: respUsers,
		},
	}
}

func (w *wechatProxy) GetUser(ctx *context.Context, id string) resp.GetUserResp {
	client := &sdk.DingTalkClient{
		AccessToken: ctx.TenantAccessToken,
		AgentId:     d.AgentId,
	}
	userDetailResp, err := client.GetUserDetail(id, nil)
	if err != nil {
		return resp.GetUserResp{Resp: resp.ErrResp(err)}
	}
	if userDetailResp.ErrCode != 0 {
		return resp.GetUserResp{Resp: resp.Resp{Code: userDetailResp.ErrCode, Msg: userDetailResp.ErrMsg}}
	}
	user := userDetailResp.UserList
	deptIdList := make([]string, 0)
	for _, deptId := range user.Department {
		deptIdList = append(deptIdList, strconv.FormatInt(deptId, 10))
	}
	return resp.GetUserResp{
		Resp: resp.SucResp(),
		Data: resp.User{
			OpenID:  user.UnionId,
			UserID:  user.UserId,
			UnionID: user.UnionId,
			Name:    user.Name,
			EnName:  user.Name,
			IsAdmin: user.IsAdmin,
			Avatar: resp.Avatar{
				Avatar72:     user.Avatar,
				Avatar240:    user.Avatar,
				Avatar640:    user.Avatar,
				AvatarOrigin: user.Avatar,
			},
			DepartmentIds: deptIdList,
		},
	}
}
