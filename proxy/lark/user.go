package lark

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
	"log"
)

func (l *larkProxy) GetUsers(ctx *context.Context, r req.GetUsersReq) resp.GetUsersResp {
	client := sdk.Tenant{
		TenantAccessToken: ctx.TenantAccessToken,
	}
	if r.DepartmentID == "" {
		allDeptIdsResp := l.GetDeptIds(ctx, req.GetDeptIdsReq{})
		if !allDeptIdsResp.Suc {
			return resp.GetUsersResp{Resp: allDeptIdsResp.Resp}
		}
		allDeptIdsResp.Data = append(allDeptIdsResp.Data, "0")
		userList := make([]resp.User, 0)
		userContains := map[string]bool{}
		for _, deptId := range allDeptIdsResp.Data {
			hasMore := true
			pageToken := ""
			for hasMore {
				deptUserListResp, err := client.GetUsersV3("", "", deptId, pageToken, 100)
				if err != nil {
					log.Println(err)
					break
				}
				if deptUserListResp.Code != 0 {
					log.Println(deptUserListResp.Code, deptUserListResp.Msg)
					break
				}
				pageToken = deptUserListResp.Data.PageToken
				hasMore = deptUserListResp.Data.HasMore

				respUserList := convertUsers(deptUserListResp.Data.Items)
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

	deptUserListResp, err := client.GetUsersV3("", "", r.DepartmentID, r.PageToken, r.PageSize)
	if err != nil {
		return resp.GetUsersResp{Resp: resp.ErrResp(err)}
	}
	if deptUserListResp.Code != 0 {
		return resp.GetUsersResp{Resp: resp.Resp{Code: deptUserListResp.Code, Msg: deptUserListResp.Msg}}
	}

	return resp.GetUsersResp{
		Resp: resp.SucResp(),
		Data: resp.GetUsersRespData{
			HasMore:   deptUserListResp.Data.HasMore,
			PageToken: deptUserListResp.Data.PageToken,
			Users:     convertUsers(deptUserListResp.Data.Items),
		},
	}
}

func convertUsers(users []vo.UserDetailInfoV3) []resp.User {
	respUsers := make([]resp.User, len(users))
	for i, user := range users {
		deptIdList := make([]string, 0)
		for _, deptId := range user.DepartmentIds {
			deptIdList = append(deptIdList, deptId)
		}
		respUsers[i] = resp.User{
			OpenID:  user.OpenId,
			UserID:  user.UserId,
			UnionID: user.UnionId,
			Name:    user.Name,
			EnName:  user.EnName,
			Email:   user.Email,
			Mobile:  user.Mobile,
			Avatar: resp.Avatar{
				Avatar72:     user.Avatar.Avatar72,
				Avatar240:    user.Avatar.Avatar240,
				Avatar640:    user.Avatar.Avatar640,
				AvatarOrigin: user.Avatar.AvatarOrigin,
			},
			DepartmentIds: deptIdList,
		}
	}
	return respUsers
}
