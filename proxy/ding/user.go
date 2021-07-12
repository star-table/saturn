package ding

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
	"log"
	"strconv"
)

func (d *dingProxy) GetUsers(ctx *context.Context, r req.GetUsersReq) resp.GetUsersResp {
	client := &sdk.DingTalkClient{
		AccessToken: ctx.TenantAccessToken,
		AgentId:     d.AgentId,
	}
	deptIds := make([]string, 0)
	if r.DepartmentID == "" {
		deptIds = append(deptIds, "1")
	} else {
		deptIds = append(deptIds, r.DepartmentID)
	}
	if r.FetchChild {
		allDeptIdsResp := d.GetDeptIds(ctx, req.GetDeptIdsReq{
			ParentId:   deptIds[0],
			FetchChild: r.FetchChild,
		})
		if !allDeptIdsResp.Suc {
			return resp.GetUsersResp{Resp: allDeptIdsResp.Resp}
		}
		deptIds = append(deptIds, allDeptIdsResp.Data...)
	}
	limit := 9999999
	if r.Limit > 0 {
		limit = r.Limit
	}
	userList := make([]resp.User, 0)
	userContains := map[string]bool{}
	for _, deptIdStr := range deptIds {
		hasMore := true
		var cursor int64 = 0
		deptId, _ := strconv.ParseInt(deptIdStr, 10, 64)
		for hasMore && (limit > len(userList)) {
			deptUserListResp, err := client.GetDeptUserListV2(deptId, cursor, 100, "", nil, "")
			if err != nil {
				log.Println(err)
				break
			}
			if deptUserListResp.ErrCode != 0 {
				log.Println(deptUserListResp.ErrCode, deptUserListResp.ErrMsg)
				break
			}
			cursor = deptUserListResp.Result.NextCursor
			hasMore = deptUserListResp.Result.HasMore

			respUserList := convertUsers(deptUserListResp.Result.List)
			for _, respUser := range respUserList {
				if !userContains[respUser.OpenID] {
					userList = append(userList, respUser)
					userContains[respUser.OpenID] = true
				}
			}
		}
		if limit <= len(userList) {
			break
		}
	}
	return resp.GetUsersResp{
		Resp: resp.SucResp(),
		Data: resp.GetUsersRespData{
			Users: userList,
		},
	}
}

func (d *dingProxy) GetUser(ctx *context.Context, id string) resp.GetUserResp {
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

func convertUsers(users []sdk.UserDetailInfoV2) []resp.User {
	respUsers := make([]resp.User, len(users))
	for i, user := range users {
		deptIdList := make([]string, 0)
		for _, deptId := range user.DeptIdList {
			deptIdList = append(deptIdList, strconv.FormatInt(deptId, 10))
		}
		respUsers[i] = resp.User{
			OpenID:  user.UnionID,
			UserID:  user.UserID,
			UnionID: user.UnionID,
			Name:    user.Name,
			EnName:  user.Name,
			Email:   user.Email,
			Mobile:  user.Mobile,
			IsAdmin: user.Admin,
			Avatar: resp.Avatar{
				Avatar72:     user.Avatar,
				Avatar240:    user.Avatar,
				Avatar640:    user.Avatar,
				AvatarOrigin: user.Avatar,
			},
			DepartmentIds: deptIdList,
		}
	}
	return respUsers
}
