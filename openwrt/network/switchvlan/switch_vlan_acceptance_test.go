//go:build acceptance.test

package switchvlan_test

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
		"device": lucirpc.String("switch0"),
		"ports":  lucirpc.String("0t 1t"),
		"vid":    lucirpc.Integer(10),
		"vlan":   lucirpc.Integer(2),
	}
	ok, err := client.CreateSection(ctx, "network", "switch_vlan", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_network_switch_vlan" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_network_switch_vlan.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_switch_vlan.testing", "device", "switch0"),
			resource.TestCheckResourceAttr("data.openwrt_network_switch_vlan.testing", "ports", "0t 1t"),
			resource.TestCheckResourceAttr("data.openwrt_network_switch_vlan.testing", "vid", "10"),
			resource.TestCheckResourceAttr("data.openwrt_network_switch_vlan.testing", "vlan", "2"),
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

resource "openwrt_network_switch_vlan" "testing" {
	device = "switch0"
	id = "testing"
	ports = "0t"
	vlan = 2
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_switch_vlan.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_switch_vlan.testing", "device", "switch0"),
			resource.TestCheckResourceAttr("openwrt_network_switch_vlan.testing", "ports", "0t"),
			resource.TestCheckResourceAttr("openwrt_network_switch_vlan.testing", "vlan", "2"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_network_switch_vlan.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_switch_vlan" "testing" {
	device = "switch0"
	id = "testing"
	ports = "0t 1t"
	vid = 10
	vlan = 2
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_switch_vlan.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_switch_vlan.testing", "device", "switch0"),
			resource.TestCheckResourceAttr("openwrt_network_switch_vlan.testing", "ports", "0t 1t"),
			resource.TestCheckResourceAttr("openwrt_network_switch_vlan.testing", "vid", "10"),
			resource.TestCheckResourceAttr("openwrt_network_switch_vlan.testing", "vlan", "2"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
