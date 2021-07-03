package proxy

import (
	"gitea.bjx.cloud/allstar/saturn/model/resp"
)

type Proxy interface {
	GetTenantAccessToken(tenantKey string) resp.GetTenantAccessTokenResp
}
