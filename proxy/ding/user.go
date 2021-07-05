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
	if r.DepartmentID == "" { // 获取所有用户
		allDeptIdsResp := d.GetDeptIds(ctx, req.GetDeptIdsReq{
			FetchChild: true,
		})
		if !allDeptIdsResp.Suc {
			return resp.GetUsersResp{Resp: allDeptIdsResp.Resp}
		}
		allDeptIdsResp.Data = append(allDeptIdsResp.Data, "1")
		userList := make([]resp.User, 0)
		userContains := map[string]bool{}
		for _, deptIdStr := range allDeptIdsResp.Data {
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
				HasMore:   false,
				PageToken: "",
				Users:     userList,
			},
		}
	}

	deptId, err := strconv.ParseInt(r.DepartmentID, 10, 64)
	if err != nil {
		return resp.GetUsersResp{Resp: resp.ErrResp(err)}
	}
	cursor, err := strconv.ParseInt(r.PageToken, 10, 64)
	if err != nil {
		return resp.GetUsersResp{Resp: resp.ErrResp(err)}
	}
	deptUserListResp, err := client.GetDeptUserListV2(deptId, cursor, r.PageSize, "", nil, "")
	if err != nil {
		return resp.GetUsersResp{Resp: resp.ErrResp(err)}
	}
	if deptUserListResp.ErrCode != 0 {
		return resp.GetUsersResp{Resp: resp.Resp{Code: deptUserListResp.ErrCode, Msg: deptUserListResp.ErrMsg}}
	}

	return resp.GetUsersResp{
		Resp: resp.SucResp(),
		Data: resp.GetUsersRespData{
			HasMore:   deptUserListResp.Result.HasMore,
			PageToken: strconv.FormatInt(deptUserListResp.Result.NextCursor, 10),
			Users:     convertUsers(deptUserListResp.Result.List),
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
