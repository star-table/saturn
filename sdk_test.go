package saturn

import (
	"fmt"
	"gitea.bjx.cloud/allstar/saturn/model/req"
	ding2 "gitea.bjx.cloud/allstar/saturn/proxy/ding"
	lark2 "gitea.bjx.cloud/allstar/saturn/proxy/lark"
	"gitea.bjx.cloud/allstar/saturn/util/json"
	"testing"
)

const (
	ding   = "ding"
	lark   = "lark"
	wechat = "wechat"
)

type TestCase struct {
	platform  string
	tenantKey string
}

var testcases = []TestCase{
	{platform: ding, tenantKey: "ding696a8496a96bcf58a1320dcb25e91351"},
	{platform: lark, tenantKey: "2ed263bf32cf1651"},
}

func assertEqual(t *testing.T, val interface{}, want interface{}) {
	if val != want {
		t.Fatal(fmt.Sprintf("%v != %v err", val, want))
	}
}

func NewTestSDK() *sdk {
	s := New()
	s.RegistryPlatform(ding, ding2.NewDingProxy(36342, "suiteocpiljyoalvbhrbi", "d1XKtyVpocDrOVJrDqPfqysmLGX7pinWS7iA8l5T7OWPd8aWZWNRfXEJrHoyb5Ng", "", "12345645615313", "1234567890123456789012345678901234567890123"))
	s.RegistryPlatform(lark, lark2.NewLarkProxy("cli_9d5e49aae9ae9101", "HDzPYfWmf8rmhsF2hHSvmhTffojOYCdI", "fa5140497af97fab6b768ea212f0a2ec4e0eff62"))
	return s
}

func TestSdk_GetCaller(t *testing.T) {
	s := NewTestSDK()
	for _, testcase := range testcases {
		_, err := s.GetCaller(testcase.platform, testcase.tenantKey)
		assertEqual(t, err, nil)
	}
}

func TestCaller_GetUsers(t *testing.T) {
	s := NewTestSDK()
	for _, testcase := range testcases {
		t.Log("testcase:", json.ToJsonIgnoreError(testcase))
		cer, err := s.GetCaller(testcase.platform, testcase.tenantKey)
		assertEqual(t, err, nil)
		r := cer.GetUsers(req.GetUsersReq{})
		t.Log(json.ToJsonIgnoreError(r))
		assertEqual(t, r.Suc, true)
	}
}

func TestCaller_GetDeptIds(t *testing.T) {
	s := NewTestSDK()
	for _, testcase := range testcases {
		t.Log("testcase:", json.ToJsonIgnoreError(testcase))
		cer, err := s.GetCaller(testcase.platform, testcase.tenantKey)
		assertEqual(t, err, nil)
		r := cer.GetDeptIds(req.GetDeptIdsReq{})
		t.Log(json.ToJsonIgnoreError(r))
		assertEqual(t, r.Suc, true)
	}
}

func TestCaller_GetDepts(t *testing.T) {
	s := NewTestSDK()
	for _, testcase := range testcases {
		t.Log("testcase:", json.ToJsonIgnoreError(testcase))
		cer, err := s.GetCaller(testcase.platform, testcase.tenantKey)
		assertEqual(t, err, nil)
		r := cer.GetDepts(req.GetDeptsReq{})
		t.Log(json.ToJsonIgnoreError(r))
		assertEqual(t, r.Suc, true)
	}
}
