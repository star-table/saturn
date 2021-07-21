package saturn

import (
	"fmt"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	ding2 "gitea.bjx.cloud/allstar/saturn/proxy/ding"
	lark2 "gitea.bjx.cloud/allstar/saturn/proxy/lark"
	wechat2 "gitea.bjx.cloud/allstar/saturn/proxy/wechat"
	welink2 "gitea.bjx.cloud/allstar/saturn/proxy/welink"
	"gitea.bjx.cloud/allstar/saturn/util/json"
	"testing"
)

const (
	ding   = "ding"
	lark   = "lark"
	wechat = "wechat"
	welink = "welink"
)

type TestCase struct {
	platform  string
	tenantKey string
}

var testcases = []TestCase{
	//{platform: ding, tenantKey: "ding95ff008aad4bebd0acaaa37764f94726"},
	//{platform: lark, tenantKey: "2ed263bf32cf1651"},
	//{platform: welink, tenantKey: "6A2F303224A44EA7999F417E54DE0B1F"},
	{platform: wechat, tenantKey: "wwf36b5e6ef0b569ac:JGTO2F8SLhbqKezTw_zjB0lYBxbi0BToHfcQmFXG9TM"},
}

func assertEqual(t *testing.T, val interface{}, want interface{}) {
	if val != want {
		t.Fatal(fmt.Sprintf("%v != %v err", val, want))
	}
}

func NewTestTenant() *SDK {
	s := New()
	s.RegistryPlatform(ding, ding2.NewDingProxy(75917, "suitegq7xfvnj3unkkbig", "CBEWjSdJ2aQV5w9crGM7TD5icSIc5tyU2VOX2UUpYq75Dh22VUOfVNYs3r3HX2oI", "12345645615313", "1234567890123456789012345678901234567890123", func() (string, error) {
		return "YMYpD97ELlElBCYnd3orCBWeINvMnWIPWXo2xCbmhepKj8wYmgYTlAq7d9lqUt9uGcwhmn8bXatODrgPFeCqCA", nil
	}))
	s.RegistryPlatform(lark, lark2.NewLarkProxy("cli_9d5e49aae9ae9101", "HDzPYfWmf8rmhsF2hHSvmhTffojOYCdI", func() (string, error) {
		return "fa5140497af97fab6b768ea212f0a2ec4e0eff62", nil
	}))
	s.RegistryPlatform(welink, welink2.NewWelinkProxy("20210716161159595718742", "241e87e6-4825-4bff-8274-3c763a2fef20", func() (string, error) {
		return "", nil
	}))
	s.RegistryPlatform(wechat, wechat2.NewWechatProxy("wwf36b5e6ef0b569ac", "", "ww9b85ae8ff033ee89", "BCLomiIeq8je52OqsXusskBMSMO8LSLnuIxpxMnfhrc", func() (string, error) {
		return "zHIaXmHYu-UWu_hOXICtN1AZ8hzh0Qv5WrIOAHTphSDjO_VpIri-pOZc9UwGRECc", nil
	}))
	return s
}

func TestTenant_GetTenant(t *testing.T) {
	s := NewTestTenant()
	for _, testcase := range testcases {
		_, err := s.GetTenant(testcase.platform, testcase.tenantKey)
		assertEqual(t, err, nil)
	}
}

func TestCaller_GetUsers(t *testing.T) {
	s := NewTestTenant()
	for _, testcase := range testcases {
		t.Log("testcase:", json.ToJsonIgnoreError(testcase))
		cer, err := s.GetTenant(testcase.platform, testcase.tenantKey)
		assertEqual(t, err, nil)
		r := cer.GetUsers(req.GetUsersReq{
			FetchChild: true,
		})
		t.Log(json.ToJsonIgnoreError(r))
		assertEqual(t, r.Suc, true)
	}
}

func TestCaller_GetUser(t *testing.T) {
	s := NewTestTenant()
	for _, testcase := range testcases {
		t.Log("testcase:", json.ToJsonIgnoreError(testcase))
		cer, err := s.GetTenant(testcase.platform, testcase.tenantKey)
		assertEqual(t, err, nil)
		r := cer.GetUser("wokYCoBwAACY3pYd5XZIuiBPB1FsItxA")
		t.Log(json.ToJsonIgnoreError(r))
		assertEqual(t, r.Suc, true)
	}
}

func TestCaller_GetDeptIds(t *testing.T) {
	s := NewTestTenant()
	for _, testcase := range testcases {
		t.Log("testcase:", json.ToJsonIgnoreError(testcase))
		cer, err := s.GetTenant(testcase.platform, testcase.tenantKey)
		assertEqual(t, err, nil)
		r := cer.GetDeptIds(req.GetDeptIdsReq{})
		t.Log(json.ToJsonIgnoreError(r))
		assertEqual(t, r.Suc, true)
	}
}

func TestCaller_GetDepts(t *testing.T) {
	s := NewTestTenant()
	for _, testcase := range testcases {
		t.Log("testcase:", json.ToJsonIgnoreError(testcase))
		cer, err := s.GetTenant(testcase.platform, testcase.tenantKey)
		assertEqual(t, err, nil)
		r := cer.GetDepts(req.GetDeptsReq{
			FetchChild: true,
		})
		t.Log(json.ToJsonIgnoreError(r))
		assertEqual(t, r.Suc, true)
	}
}
