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
	userList := make([]resp.User, 0)
	userContains := map[string]bool{}
	for _, deptIdStr := range deptIds {
		hasMore := true
		var cursor int64 = 0
		deptId, _ := strconv.ParseInt(deptIdStr, 10, 64)
		for hasMore {
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
	}
	return resp.GetUsersResp{
		Resp: resp.SucResp(),
		Data: resp.GetUsersRespData{
			Users: userList,
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
