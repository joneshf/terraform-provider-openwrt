//go:build acceptance.test

package odhcpd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestDataSourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	client := openWrtServer.LuCIRPCClient(
		ctx,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()
	options := lucirpc.Options{
		"leasefile":    lucirpc.String("/tmp/leasefile"),
		"leasetrigger": lucirpc.String("/tmp/leasetrigger"),
	}
	ok, err := client.CreateSection(ctx, "dhcp", "odhcpd", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_dhcp_odhcpd" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_dhcp_odhcpd.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_odhcpd.testing", "leasefile", "/tmp/leasefile"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_odhcpd.testing", "leasetrigger", "/tmp/leasetrigger"),
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

	createAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_dhcp_odhcpd" "testing" {
	id = "testing"
	leasefile = "/tmp/leasefile"
	leasetrigger = "/tmp/leasetrigger"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_odhcpd.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_odhcpd.testing", "leasefile", "/tmp/leasefile"),
			resource.TestCheckResourceAttr("openwrt_dhcp_odhcpd.testing", "leasetrigger", "/tmp/leasetrigger"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_dhcp_odhcpd.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_dhcp_odhcpd" "testing" {
	id = "testing"
	leasefile = "/tmp/leasefile"
	leasetrigger = "/tmp/leasetrigger"
	legacy = true
	loglevel = 6
	maindhcp = true
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_odhcpd.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_odhcpd.testing", "leasefile", "/tmp/leasefile"),
			resource.TestCheckResourceAttr("openwrt_dhcp_odhcpd.testing", "leasetrigger", "/tmp/leasetrigger"),
			resource.TestCheckResourceAttr("openwrt_dhcp_odhcpd.testing", "legacy", "true"),
			resource.TestCheckResourceAttr("openwrt_dhcp_odhcpd.testing", "loglevel", "6"),
			resource.TestCheckResourceAttr("openwrt_dhcp_odhcpd.testing", "maindhcp", "true"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
