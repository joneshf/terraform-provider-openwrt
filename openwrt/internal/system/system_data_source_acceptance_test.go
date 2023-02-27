//go:build acceptance.test

package system_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
)

func TestSystemSystemDataSourceAcceptance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"openwrt": providerserver.NewProtocol6WithError(openwrt.New("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "openwrt" {
	hostname = "localhost"
	port = 8080
}

data "openwrt_system_system" "this" {
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "id", "cfg01e48a"),
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "hostname", "OpenWrt"),
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "log_size", "64"),
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "timezone", "UTC"),
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "ttylogin", "false"),
				),
			},
		},
	})
}
