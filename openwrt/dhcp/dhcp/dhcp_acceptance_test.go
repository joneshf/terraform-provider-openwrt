//go:build acceptance.test

package dhcp_test

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
		"interface": lucirpc.String("testing"),
		"leasetime": lucirpc.String("12h"),
		"limit":     lucirpc.Integer(150),
		"start":     lucirpc.Integer(100),
	}
	ok, err := client.CreateSection(ctx, "dhcp", "dhcp", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_dhcp_dhcp" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_dhcp_dhcp.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_dhcp.testing", "interface", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_dhcp.testing", "leasetime", "12h"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_dhcp.testing", "limit", "150"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_dhcp.testing", "start", "100"),
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

resource "openwrt_dhcp_dhcp" "testing" {
	id = "testing"
	ignore = true
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "ignore", "true"),
			resource.TestCheckNoResourceAttr("openwrt_dhcp_dhcp.testing", "dhcpv4"),
			resource.TestCheckNoResourceAttr("openwrt_dhcp_dhcp.testing", "dhcpv6"),
			resource.TestCheckNoResourceAttr("openwrt_dhcp_dhcp.testing", "interface"),
			resource.TestCheckNoResourceAttr("openwrt_dhcp_dhcp.testing", "leasetime"),
			resource.TestCheckNoResourceAttr("openwrt_dhcp_dhcp.testing", "limit"),
			resource.TestCheckNoResourceAttr("openwrt_dhcp_dhcp.testing", "ra_flags"),
			resource.TestCheckNoResourceAttr("openwrt_dhcp_dhcp.testing", "start"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_dhcp_dhcp.testing",
	}
	ignoreWithOtherAttributes := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_dhcp_dhcp" "testing" {
	dhcpv4 = "server"
	dhcpv6 = "disabled"
	id = "testing"
	ignore = true
	interface = "testing"
	leasetime = "12h"
	limit = 150
	ra_flags = [
		"managed-config",
		"other-config",
	]
	start = 100
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "dhcpv4", "server"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "dhcpv6", "disabled"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "ignore", "true"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "interface", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "leasetime", "12h"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "limit", "150"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "ra_flags.0", "managed-config"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "ra_flags.1", "other-config"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "start", "100"),
		),
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_dhcp_dhcp" "testing" {
	dhcpv4 = "server"
	dhcpv6 = "disabled"
	id = "testing"
	interface = "testing"
	leasetime = "12h"
	limit = 150
	ra_flags = [
		"managed-config",
		"other-config",
	]
	start = 100
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "dhcpv4", "server"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "dhcpv6", "disabled"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "interface", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "leasetime", "12h"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "limit", "150"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "ra_flags.0", "managed-config"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "ra_flags.1", "other-config"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dhcp.testing", "start", "100"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		ignoreWithOtherAttributes,
		updateAndReadResource,
	)
}
