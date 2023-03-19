//go:build acceptance.test

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestNetworkGlobalsDataSourceAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client, providerBlock := acceptancetest.AuthenticatedClientWithProviderBlock(
		ctx,
		*dockerPool,
		t,
	)
	options := lucirpc.Options{
		"packet_steering": lucirpc.Boolean(false),
		"ula_prefix":      lucirpc.String("fd12:3456:789a::/48"),
	}
	ok, err := client.CreateSection(ctx, "network", "globals", "globals", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_network_globals" "this" {
	id = "globals"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_network_globals.this", "id", "globals"),
			resource.TestCheckResourceAttr("data.openwrt_network_globals.this", "packet_steering", "false"),
			resource.TestCheckResourceAttr("data.openwrt_network_globals.this", "ula_prefix", "fd12:3456:789a::/48"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

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
