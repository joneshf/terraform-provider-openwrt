//go:build acceptance.test

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
)

func TestNetworkGlobalsResourceAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	providerBlock := acceptancetest.RunOpenWrtServerWithProviderBlock(
		ctx,
		*dockerPool,
		t,
	)

	createAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_globals" "this" {
	id = "globals"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_globals.this", "id", "globals"),
			resource.TestCheckNoResourceAttr("openwrt_network_globals.this", "network_steering"),
			resource.TestCheckNoResourceAttr("openwrt_network_globals.this", "ula_prefix"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_network_globals.this",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_globals" "this" {
	id = "globals"
	packet_steering = false
	ula_prefix = "fd12:3456:789a::/48"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_globals.this", "id", "globals"),
			resource.TestCheckResourceAttr("openwrt_network_globals.this", "packet_steering", "false"),
			resource.TestCheckResourceAttr("openwrt_network_globals.this", "ula_prefix", "fd12:3456:789a::/48"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
