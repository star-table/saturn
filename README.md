## 介绍
多平台聚合SDK，目前已支持钉钉、飞书、企业微信、WeLink等平台。

## 已支持
 - [x] 企业认证
 - [x] 获取成员列表
 - [x] 获取成员信息  
 - [x] 获取部门id列表
 - [x] 获取部门列表
 - [x] 获取根部门  
 - [x] 发送通知(部分)
 - [x] 免密登录

目前可以通过saturn完成组织架构的同步以及免密登录，已满足己用。

## 使用
初始化saturn，配置不同平台应用的key和secret
```go
s := New()
s.RegistryPlatform("ding", ding.NewDingProxy(75917, "suitegq7xfvnj3unkkbig", "CBEWjSdJ2aQV5w9crGM7TD5icSIc5tyU2VOX2UUpYq75Dh22VUOfVNYs3r3HX2oI", "12345645615313", "1234567890123456789012345678901234567890123", func() (string, error) {
    return "YMYpD97ELlElBCYnd3orCBWeINvMnWIPWXo2xCbmhepKj8wYmgYTlAq7d9lqUt9uGcwhmn8bXatODrgPFeCqCA", nil
}))
s.RegistryPlatform("lark", lark.NewLarkProxy("cli_9d5e49aae9ae9101", "HDzPYfWmf8rmhsF2hHSvmhTffojOYCdI", func() (string, error) {
    return "fa5140497af97fab6b768ea212f0a2ec4e0eff62", nil
}))
s.RegistryPlatform("welink", welink.NewWelinkProxy("20210716161159595718742", "241e87e6-4825-4bff-8274-3c763a2fef20", func() (string, error) {
    return "", nil
}))
s.RegistryPlatform("wechat", wechat.NewWechatProxy("wwf36b5e6ef0b569ac", "", "ww9b85ae8ff033ee89", "BCLomiIeq8je52OqsXusskBMSMO8LSLnuIxpxMnfhrc", func() (string, error) {
    return "zHIaXmHYu-UWu_hOXICtNz0omNxiAxzCUfziZi-72hHxYJwOfzGfDbBfbv2EJfZr", nil
}))
```
应用级别接口调用
```go
app, err := sdk.GetApp(platform)
if err != nil {
    return err
}
app.CodeLogin
```
企业级别接口调用
```go
tenant, err := sdk.GetTenant(platform, tenantKey)
if err != nil {
    return err
}
tenant.GetUser
tenant.GetUsers
tenant.GetDeptIds
tenant.GetDepts
tenant.GetRootDept
```

## 相关资源
- 企业微信：https://github.com/xen0n/go-workwx
- 飞书：https://github.com/galaxy-book/feishu-sdk-golang
- 钉钉：https://github.com/galaxy-book/dingtalk-sdk-golang
- WeLink: https://open.welink.huaweicloud.com/docs/#/qdmtm8/uug541/sx2v22
