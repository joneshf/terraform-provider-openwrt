//go:build acceptance.test

package system_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
)

func TestDataSourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_system_system" "this" {
	id = "cfg01e48a"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "id", "cfg01e48a"),
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "hostname", "OpenWrt"),
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "log_size", "64"),
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "timezone", "UTC"),
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "ttylogin", "false"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

func TestResourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	importValidation := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_system_system" "this" {
	id = "cfg01e48a"
}
`,
			providerBlock,
		),
		ImportState:        true,
		ImportStateId:      "cfg01e48a",
		ImportStatePersist: true,
		ResourceName:       "openwrt_system_system.this",
	}

	readResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_system_system" "this" {
	id = "cfg01e48a"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_system_system.this", "id", "cfg01e48a"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "hostname", "OpenWrt"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "log_size", "64"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "timezone", "UTC"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "ttylogin", "false"),
		),
	}

	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_system_system" "this" {
	hostname = "OpenWRT"
	id = "cfg01e48a"
	log_size = 64
	timezone = "UTC"
	ttylogin = false
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_system_system.this", "id", "cfg01e48a"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "hostname", "OpenWRT"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "log_size", "64"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "timezone", "UTC"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "ttylogin", "false"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		importValidation,
		readResource,
		updateAndReadResource,
	)
}
