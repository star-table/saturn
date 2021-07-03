package saturn

import (
	"errors"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/proxy"
)

var (
	PlatformNotRegistryError = errors.New("platform not registry. ")
)

type sdk struct {
	platforms map[string]proxy.Proxy
}

func New() *sdk {
	s := &sdk{platforms: map[string]proxy.Proxy{}}
	return s
}

func (s *sdk) RegistryPlatform(platform string, p proxy.Proxy) {
	s.platforms[platform] = p
}

func (s *sdk) GetTenantAccessToken(platform string, tenantKey string) resp.GetTenantAccessTokenResp {
	if platform, ok := s.platforms[platform]; ok {
		return platform.GetTenantAccessToken(tenantKey)
	}
	return resp.GetTenantAccessTokenResp{Resp: resp.ErrResp(PlatformNotRegistryError)}
}
