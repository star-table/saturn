package proxy

import (
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
)

type Proxy interface {
	// GetTenantAccessToken 获取企业认证Token
	GetTenantAccessToken(tenantKey string) resp.GetTenantAccessTokenResp
	// GetUsers 获取用户列表，部门ID未指定时查询所有用户
	GetUsers(ctx *context.Context, req req.GetUsersReq) resp.GetUsersResp
	// GetDeptIds 获取部门ID列表，不包含顶级部门及父部门
	GetDeptIds(ctx *context.Context, req req.GetDeptIdsReq) resp.GetDeptIdsResp
	// GetDepts 获取部门列表
	GetDepts(ctx *context.Context, req req.GetDeptsReq) resp.GetDeptsResp
}
