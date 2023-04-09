//go:build acceptance.test

package domain_test

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
		"name": lucirpc.String("testing"),
	}
	ok, err := client.CreateSection(ctx, "dhcp", "domain", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_dhcp_domain" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_dhcp_domain.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_domain.testing", "ip", "192.168.1.50"),
			resource.TestCheckResourceAttr("data.openwrt_dhcp_domain.testing", "name", "testing"),
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

resource "openwrt_dhcp_domain" "testing" {
	id = "testing"
	ip = "192.168.1.50"
	name = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_domain.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_domain.testing", "ip", "192.168.1.50"),
			resource.TestCheckResourceAttr("openwrt_dhcp_domain.testing", "name", "testing"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_dhcp_domain.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_dhcp_domain" "testing" {
	id = "testing"
	ip = "192.168.1.51"
	name = "testing-1"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_dhcp_domain.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_dhcp_domain.testing", "ip", "192.168.1.51"),
			resource.TestCheckResourceAttr("openwrt_dhcp_domain.testing", "name", "testing-1"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
