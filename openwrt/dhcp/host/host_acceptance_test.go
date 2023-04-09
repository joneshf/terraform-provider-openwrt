//go:build acceptance.test

package host_test

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
		"ip":   lucirpc.String("192.168.1.50"),
		"mac":  lucirpc.String("12:34:56:78:90:ab"),
		"name": lucirpc.String("testing"),
	}
	ok, err := client.CreateSection(ctx, "dhcp", "host", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_dhcp_host" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_dhcp_host.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_host.testing", "ip", "192.168.1.50"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_host.testing", "mac", "12:34:56:78:90:ab"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_host.testing", "name", "testing"),
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

resource "openwrt_dhcp_host" "testing" {
	id = "testing"
	ip = "192.168.1.50"
	mac = "12:34:56:78:90:ab"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_host.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_host.testing", "ip", "192.168.1.50"),
			resource.TestCheckResourceAttr("openwrt_dhcp_host.testing", "mac", "12:34:56:78:90:ab"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_dhcp_host.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_dhcp_host" "testing" {
	id = "testing"
	ip = "192.168.1.50"
	mac = "12:34:56:78:90:ab"
	name = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_host.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_host.testing", "ip", "192.168.1.50"),
			resource.TestCheckResourceAttr("openwrt_dhcp_host.testing", "mac", "12:34:56:78:90:ab"),
			resource.TestCheckResourceAttr("openwrt_dhcp_host.testing", "name", "testing"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
