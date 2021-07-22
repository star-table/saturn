package lark

import (
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"github.com/galaxy-book/feishu-sdk-golang/sdk"
	"github.com/galaxy-book/saturn/model/context"
	"github.com/galaxy-book/saturn/model/req"
	"github.com/galaxy-book/saturn/model/resp"
	"log"
)

func (l *larkProxy) GetUsers(ctx *context.Context, r req.GetUsersReq) resp.GetUsersResp {
	client := sdk.Tenant{
		TenantAccessToken: ctx.TenantAccessToken,
	}
	deptIds := make([]string, 0)
	if r.DepartmentID == "" {
		deptIds = append(deptIds, "0")
	} else {
		deptIds = append(deptIds, r.DepartmentID)
	}
	if r.FetchChild {
		allDeptIdsResp := l.GetDeptIds(ctx, req.GetDeptIdsReq{
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
	for _, deptId := range deptIds {
		hasMore := true
		pageToken := ""
		for hasMore && (limit > len(userList)) {
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

func (l *larkProxy) GetUser(ctx *context.Context, id string) resp.GetUserResp {
	client := sdk.Tenant{
		TenantAccessToken: ctx.TenantAccessToken,
	}
	userBatchResp, err := client.GetUserBatchGetV2(nil, []string{id})
	if err != nil {
		return resp.GetUserResp{Resp: resp.ErrResp(err)}
	}
	if userBatchResp.Code != 0 {
		return resp.GetUserResp{Resp: resp.Resp{Code: userBatchResp.Code, Msg: userBatchResp.Msg}}
	}
	if len(userBatchResp.Data.Users) == 0 {
		return resp.GetUserResp{Resp: resp.Resp{Code: -1, Msg: "user not exist"}}
	}
	user := userBatchResp.Data.Users[0]
	return resp.GetUserResp{
		Resp: resp.SucResp(),
		Data: resp.User{
			OpenID:  user.OpenId,
			UserID:  user.UserId,
			UnionID: user.UnionId,
			Name:    user.Name,
			EnName:  user.EnName,
			Email:   user.Email,
			Mobile:  user.Mobile,
			IsAdmin: user.IsTenantManager,
			Avatar: resp.Avatar{
				Avatar72:     user.Avatar.Avatar72,
				Avatar240:    user.Avatar.Avatar240,
				Avatar640:    user.Avatar.Avatar640,
				AvatarOrigin: user.Avatar.AvatarOrigin,
			},
			DepartmentIds: user.Departments,
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
			IsAdmin: user.IsTenantManager,
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
