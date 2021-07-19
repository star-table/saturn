package welink

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/welink"
	"github.com/galaxy-book/feishu-sdk-golang/core/model/vo"
	"log"
	"strconv"
)

func (w *welinkProxy) GetUsers(ctx *context.Context, r req.GetUsersReq) resp.GetUsersResp {
	deptIds := make([]string, 0)
	if r.DepartmentID == "" {
		deptIds = append(deptIds, "0")
	} else {
		deptIds = append(deptIds, r.DepartmentID)
	}
	if r.FetchChild {
		allDeptIdsResp := w.GetDeptIds(ctx, req.GetDeptIdsReq{
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

	usersReq := welink.ContactV1UserUsersRequest()
	usersReq.SetAccessToken(ctx.TenantAccessToken)
	err := usersReq.SetUrl("https://open.welink.huaweicloud.com/api/contact/v2/user/users")
	if err != nil {
		return resp.GetUsersResp{Resp: resp.ErrResp(err)}
	}
	for _, deptId := range deptIds {
		pageNo := 1
		pageSize := 50
		for limit > len(userList) {
			usersReq.DeptCode = deptId
			usersReq.PageNo = strconv.Itoa(pageNo)
			usersReq.PageSize = strconv.Itoa(pageSize)
			err, usersResp := usersReq.GetResponse()
			if err != nil {
				log.Println(err)
				break
			}
			if usersResp.Code != "0" {
				return resp.GetUsersResp{Resp: resp.Resp{Code: -1, Msg: usersResp.Message}}
			}
			for _, respUser := range *usersResp.Data {
				if !userContains[respUser.UserId] {
					userList = append(userList, resp.User{
						OpenID:        respUser.UserId,
						UserID:        respUser.UserId,
						UnionID:       respUser.UserId,
						Name:          respUser.UserNameCn,
						EnName:        respUser.UserNameEn,
						Email:         respUser.UserEmail,
						Mobile:        respUser.MobileNumber,
						IsAdmin:       true,
						DepartmentIds: []string{respUser.DeptCode},
					})
					userContains[respUser.UserId] = true
				}
			}
			if len(*usersResp.Data) < limit {
				break
			}
			pageNo++
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

func (w *welinkProxy) GetUser(ctx *context.Context, id string) resp.GetUserResp {
	userInfoReq := welink.ContactV1UsersRequest()
	userInfoReq.SetAccessToken(ctx.TenantAccessToken)
	userInfoReq.UserId = id
	err := userInfoReq.SetUrl("https://open.welink.huaweicloud.com/api/contact/v3/users/simple")
	if err != nil {
		return resp.GetUserResp{Resp: resp.ErrResp(err)}
	}
	err, respUser := userInfoReq.GetResponse()
	if err != nil {
		return resp.GetUserResp{Resp: resp.ErrResp(err)}
	}
	if respUser.Code != "0" {
		return resp.GetUserResp{Resp: resp.Resp{Code: -1, Msg: respUser.Message}}
	}
	return resp.GetUserResp{
		Resp: resp.SucResp(),
		Data: resp.User{
			OpenID:        respUser.UserId,
			UserID:        respUser.UserId,
			UnionID:       respUser.UserId,
			Name:          respUser.UserNameCn,
			EnName:        respUser.UserNameEn,
			Email:         respUser.UserEmail,
			Mobile:        respUser.MobileNumber,
			IsAdmin:       true,
			DepartmentIds: []string{respUser.DeptCode},
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
