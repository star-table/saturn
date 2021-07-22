package saturn

import (
	"errors"
	"github.com/galaxy-book/saturn/model/context"
	"github.com/galaxy-book/saturn/model/req"
	"github.com/galaxy-book/saturn/model/resp"
	"github.com/galaxy-book/saturn/proxy"
	"time"
)

var (
	PlatformNotRegistryError = errors.New("platform not registry. ")
)

type SDK struct {
	platforms map[string]proxy.Proxy
}

type Tenant struct {
	s *SDK
	c *context.Context
	p proxy.Proxy
}

type App struct {
	s         *SDK
	p         proxy.Proxy
	appKey    string
	appSecret string
}

func New() *SDK {
	s := &SDK{platforms: map[string]proxy.Proxy{}}
	return s
}

func (s *SDK) RegistryPlatform(platform string, p proxy.Proxy) {
	s.platforms[platform] = p
}

func (s *SDK) GetContext(platform string, tenantKey string) (*context.Context, proxy.Proxy, error) {
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

func (s *SDK) GetTenant(platform string, tenantKey string) (*Tenant, error) {
	c, p, err := s.GetContext(platform, tenantKey)
	if err != nil {
		return nil, err
	}
	return &Tenant{
		s: s,
		c: c,
		p: p,
	}, nil
}

func (s *SDK) GetApp(platform string) (*App, error) {
	p, ok := s.platforms[platform]
	if !ok {
		return nil, PlatformNotRegistryError
	}
	return &App{
		s: s,
		p: p,
	}, nil
}

func (c *Tenant) context() (*context.Context, error) {
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

func (c *Tenant) GetUsers(req req.GetUsersReq) resp.GetUsersResp {
	ctx, err := c.context()
	if err != nil {
		return resp.GetUsersResp{Resp: resp.ErrResp(err)}
	}
	return c.p.GetUsers(ctx, req)
}

func (c *Tenant) GetDeptIds(req req.GetDeptIdsReq) resp.GetDeptIdsResp {
	ctx, err := c.context()
	if err != nil {
		return resp.GetDeptIdsResp{Resp: resp.ErrResp(err)}
	}
	return c.p.GetDeptIds(ctx, req)
}

func (c *Tenant) GetDepts(req req.GetDeptsReq) resp.GetDeptsResp {
	ctx, err := c.context()
	if err != nil {
		return resp.GetDeptsResp{Resp: resp.ErrResp(err)}
	}
	return c.p.GetDepts(ctx, req)
}

func (c *Tenant) SendMsg(req req.SendMsgReq) resp.SendMsgResp {
	ctx, err := c.context()
	if err != nil {
		return resp.SendMsgResp{Resp: resp.ErrResp(err)}
	}
	return c.p.SendMsg(ctx, req)
}

func (c *Tenant) GetRootDept() resp.GetRootDeptResp {
	ctx, err := c.context()
	if err != nil {
		return resp.GetRootDeptResp{Resp: resp.ErrResp(err)}
	}
	return c.p.GetRootDept(ctx)
}

func (c *Tenant) GetUser(id string) resp.GetUserResp {
	ctx, err := c.context()
	if err != nil {
		return resp.GetUserResp{Resp: resp.ErrResp(err)}
	}
	return c.p.GetUser(ctx, id)
}

func (c *App) CodeLogin(tenantKey, code string) resp.CodeLoginResp {
	return c.p.CodeLogin(tenantKey, code)
}
