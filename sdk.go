package saturn

import (
	"errors"
	"gitea.bjx.cloud/allstar/saturn/model/context"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	"gitea.bjx.cloud/allstar/saturn/model/resp"
	"gitea.bjx.cloud/allstar/saturn/proxy"
	"time"
)

var (
	PlatformNotRegistryError = errors.New("platform not registry. ")
)

type sdk struct {
	platforms map[string]proxy.Proxy
}

type caller struct {
	s *sdk
	c *context.Context
	p proxy.Proxy
}

func New() *sdk {
	s := &sdk{platforms: map[string]proxy.Proxy{}}
	return s
}

func (s *sdk) RegistryPlatform(platform string, p proxy.Proxy) {
	s.platforms[platform] = p
}

func (s *sdk) GetContext(platform string, tenantKey string) (*context.Context, proxy.Proxy, error) {
	p, ok := s.platforms[platform]
	if !ok {
		return nil, nil, PlatformNotRegistryError
	}
	tenantAccessTokenResp := p.GetTenantAccessToken(tenantKey)
	if !tenantAccessTokenResp.Suc {
		return nil, nil, tenantAccessTokenResp.Error()
	}
	return &context.Context{
		Platform:                    platform,
		TenantKey:                   tenantKey,
		TenantAccessToken:           tenantAccessTokenResp.Data.Token,
		TenantAccessTokenExpire:     tenantAccessTokenResp.Data.Expire,
		TenantAccessTokenExpireTime: time.Now().Add(time.Duration(tenantAccessTokenResp.Data.Expire) * time.Second),
	}, p, nil
}

func (s *sdk) GetCaller(platform string, tenantKey string) (*caller, error) {
	c, p, err := s.GetContext(platform, tenantKey)
	if err != nil {
		return nil, err
	}
	return &caller{
		s: s,
		c: c,
		p: p,
	}, nil
}

func (c *caller) context() (*context.Context, error) {
	if c.c.Valid() {
		return c.c, nil
	}
	ctx, _, err := c.s.GetContext(c.c.Platform, c.c.TenantKey)
	if err != nil {
		return nil, err
	}
	c.c = ctx
	return c.c, nil
}

func (c *caller) GetUsers(req req.GetUsersReq) resp.GetUsersResp {
	ctx, err := c.context()
	if err != nil {
		return resp.GetUsersResp{Resp: resp.ErrResp(err)}
	}
	return c.p.GetUsers(ctx, req)
}

func (c *caller) GetDeptIds(req req.GetDeptIdsReq) resp.GetDeptIdsResp {
	ctx, err := c.context()
	if err != nil {
		return resp.GetDeptIdsResp{Resp: resp.ErrResp(err)}
	}
	return c.p.GetDeptIds(ctx, req)
}

func (c *caller) GetDepts(req req.GetDeptsReq) resp.GetDeptsResp {
	ctx, err := c.context()
	if err != nil {
		return resp.GetDeptsResp{Resp: resp.ErrResp(err)}
	}
	return c.p.GetDepts(ctx, req)
}
