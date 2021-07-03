package saturn

import (
	"fmt"
	"gitea.bjx.cloud/allstar/saturn/proxy"
	"testing"
)

const (
	ding = "ding"
	lark = "lark"
)

type TestCase struct {
	platform  string
	tenantKey string
}

var testcases = []TestCase{
	{platform: ding, tenantKey: "ding696a8496a96bcf58a1320dcb25e91351"},
	{platform: lark, tenantKey: "2e99b3ab0b0f1654"},
}

func assertEqual(t *testing.T, val interface{}, want interface{}) {
	if val != want {
		t.Fatal(fmt.Sprintf("%v != %v err", val, want))
	}
}

func NewSDK() *sdk {
	s := New()
	s.RegistryPlatform(ding, proxy.NewDingProxy(36342, "suiteocpiljyoalvbhrbi", "d1XKtyVpocDrOVJrDqPfqysmLGX7pinWS7iA8l5T7OWPd8aWZWNRfXEJrHoyb5Ng", "", "12345645615313", "1234567890123456789012345678901234567890123"))
	s.RegistryPlatform(lark, proxy.NewLarkProxy("cli_9d5e49aae9ae9101", "HDzPYfWmf8rmhsF2hHSvmhTffojOYCdI", "a63005f5f7e2ba8e492485969c3880c8e4aa8f4d"))
	return s
}

func TestNew(t *testing.T) {
	s := NewSDK()
	for _, testcase := range testcases {
		resp := s.GetTenantAccessToken(testcase.platform, testcase.tenantKey)
		fmt.Println(resp)
		assertEqual(t, resp.Suc, true)
	}
}
