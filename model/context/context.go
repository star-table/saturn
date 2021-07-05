package context

import "time"

type Context struct {
	Platform                    string    `json:"platform"`
	TenantKey                   string    `json:"tenantKey"`
	TenantAccessToken           string    `json:"tenantAccessToken"`
	TenantAccessTokenExpire     int64     `json:"tenantAccessTokenExpire"`
	TenantAccessTokenExpireTime time.Time `json:"tenantAccessTokenExpireTime"`
}

func (c *Context) Valid() bool {
	return time.Now().Add(5 * time.Minute).Before(c.TenantAccessTokenExpireTime)
}
