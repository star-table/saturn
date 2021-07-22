package wechat

import (
	c1 "context"
	"github.com/galaxy-book/saturn/model/context"
	"github.com/galaxy-book/saturn/model/req"
	"github.com/galaxy-book/saturn/model/resp"
	"github.com/galaxy-book/saturn/util/json"
	"github.com/galaxy-book/work-wechat"
	"strconv"
)

func (w *wechatProxy) GetUsers(ctx *context.Context, r req.GetUsersReq) resp.GetUsersResp {
	if r.DepartmentID == "" {
		r.DepartmentID = "1"
	}
	fetchChild := 0
	if r.FetchChild {
		fetchChild = 1
	}
	action := work.GetDeptMemberList(ctx.TenantAccessToken, r.DepartmentID, fetchChild)
	respBody, err := action.DoRequest(c1.Background())
	if err != nil {
		return resp.GetUsersResp{Resp: resp.ErrResp(err)}
	}
	deptMemberListResp := work.GetDeptMemberListResp{}
	json.FromJsonIgnoreError(string(respBody), &deptMemberListResp)
	if deptMemberListResp.ErrCode != 0 {
		return resp.GetUsersResp{Resp: resp.Resp{Code: deptMemberListResp.ErrCode, Msg: deptMemberListResp.ErrMsg}}
	}

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
	action := work.GetUserInfoAction(ctx.TenantAccessToken, id)
	respBody, err := action.DoRequest(c1.Background())
	if err != nil {
		return resp.GetUserResp{Resp: resp.ErrResp(err)}
	}
	userInfoDetailsResp := work.UserInfo{}
	json.FromJsonIgnoreError(string(respBody), &userInfoDetailsResp)
	if userInfoDetailsResp.ErrCode != 0 {
		return resp.GetUserResp{Resp: resp.Resp{Code: userInfoDetailsResp.ErrCode, Msg: userInfoDetailsResp.ErrMsg}}
	}
	deptIdList := make([]string, 0)
	for _, deptId := range userInfoDetailsResp.Department {
		deptIdList = append(deptIdList, strconv.Itoa(deptId))
	}
	return resp.GetUserResp{
		Resp: resp.SucResp(),
		Data: resp.User{
			OpenID:  userInfoDetailsResp.OpenUserid,
			UserID:  userInfoDetailsResp.Userid,
			UnionID: userInfoDetailsResp.OpenUserid,
			Name:    userInfoDetailsResp.Name,
			EnName:  userInfoDetailsResp.Name,
			Email:   userInfoDetailsResp.Email,
			Mobile:  userInfoDetailsResp.Mobile,
			IsAdmin: false,
			Avatar: resp.Avatar{
				Avatar72:     userInfoDetailsResp.Avatar,
				Avatar240:    userInfoDetailsResp.Avatar,
				Avatar640:    userInfoDetailsResp.Avatar,
				AvatarOrigin: userInfoDetailsResp.Avatar,
			},
			DepartmentIds: deptIdList,
		},
	}
}
