//go:build acceptance.test

package networkswitch_test

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
		"name":        lucirpc.String("switch0"),
		"enable_vlan": lucirpc.Boolean(true),
	}
	ok, err := client.CreateSection(ctx, "network", "switch", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_network_switch" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_network_switch.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_switch.testing", "enable_vlan", "true"),
			resource.TestCheckResourceAttr("data.openwrt_network_switch.testing", "name", "switch0"),
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

resource "openwrt_network_switch" "testing" {
	id = "testing"
	name = "switch0"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "name", "switch0"),
			resource.TestCheckNoResourceAttr("openwrt_network_switch.testing", "enable_vlan"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_network_switch.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_switch" "testing" {
	enable_vlan = true
	id = "testing"
	name = "switch0"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "enable_vlan", "true"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "name", "switch0"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}

func TestResourceMirrorMonitorPortWithEnableMirrorReceivedAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	step := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_switch" "testing" {
	enable_mirror_rx = true
	id = "testing"
	mirror_monitor_port = 3
	name = "switch0"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "enable_mirror_rx", "true"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "mirror_monitor_port", "3"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "name", "switch0"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		step,
	)
}

func TestResourceMirrorMonitorPortWithEnableMirrorTransmittedAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	step := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_switch" "testing" {
	enable_mirror_tx = true
	id = "testing"
	mirror_monitor_port = 3
	name = "switch0"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "enable_mirror_tx", "true"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "mirror_monitor_port", "3"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "name", "switch0"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		step,
	)
}

func TestResourceMirrorSourcePortWithEnableMirrorReceivedAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	step := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_switch" "testing" {
	enable_mirror_rx = true
	id = "testing"
	mirror_source_port = 3
	name = "switch0"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "enable_mirror_rx", "true"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "mirror_source_port", "3"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "name", "switch0"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		step,
	)
}

func TestResourceMirrorSourcePortWithEnableMirrorTransmittedAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	step := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_switch" "testing" {
	enable_mirror_tx = true
	id = "testing"
	mirror_source_port = 3
	name = "switch0"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "enable_mirror_tx", "true"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "mirror_source_port", "3"),
			resource.TestCheckResourceAttr("openwrt_network_switch.testing", "name", "switch0"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		step,
	)
}
