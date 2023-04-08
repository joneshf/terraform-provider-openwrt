//go:build acceptance.test

package dnsmasq_test

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
		"domain":            lucirpc.String("testing"),
		"rebind_protection": lucirpc.Boolean(true),
	}
	ok, err := client.CreateSection(ctx, "dhcp", "dnsmasq", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_dhcp_dnsmasq" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_dhcp_dnsmasq.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_dnsmasq.testing", "domain", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_dnsmasq.testing", "rebind_protection", "true"),
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

resource "openwrt_dhcp_dnsmasq" "testing" {
	domain = "testing"
	id = "testing"
	rebind_protection = true
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_dnsmasq.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dnsmasq.testing", "domain", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dnsmasq.testing", "rebind_protection", "true"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_dhcp_dnsmasq.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_dhcp_dnsmasq" "testing" {
	domain = "testing"
	expandhosts = true
	id = "testing"
	local = "/testing/"
	rebind_localhost = true
	rebind_protection = true
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_dnsmasq.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dnsmasq.testing", "domain", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dnsmasq.testing", "expandhosts", "true"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dnsmasq.testing", "local", "/testing/"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dnsmasq.testing", "rebind_localhost", "true"),
			resource.TestCheckResourceAttr("openwrt_dhcp_dnsmasq.testing", "rebind_protection", "true"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
