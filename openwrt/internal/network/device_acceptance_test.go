//go:build acceptance.test

package network_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestNetworkDeviceDataSourceAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client, providerBlock := acceptancetest.AuthenticatedClientWithProviderBlock(
		ctx,
		*dockerPool,
		t,
	)
	options := lucirpc.Options{
		"name":  json.RawMessage(`"br-testing"`),
		"ports": json.RawMessage(`["eth0", "eth1"]`),
		"type":  json.RawMessage(`"bridge"`),
	}
	ok, err := client.CreateSection(ctx, "network", "device", "br_testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_network_device" "this" {
	id = "br_testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "id", "br_testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "name", "br-testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "ports.0", "eth0"),
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "ports.1", "eth1"),
			resource.TestCheckResourceAttr("data.openwrt_network_device.this", "type", "bridge"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

func TestNetworkDeviceResourceAcceptance(t *testing.T) {
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

resource "openwrt_network_device" "br_testing" {
	id = "br_testing"
	name = "br-testing"
	ports = [
		"eth0",
		"eth1",
	]
	type = "bridge"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "id", "br_testing"),
			resource.TestCheckNoResourceAttr("openwrt_network_device.br_testing", "macaddr"),
			resource.TestCheckNoResourceAttr("openwrt_network_device.br_testing", "mtu"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "name", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.0", "eth0"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.1", "eth1"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "type", "bridge"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_network_device.br_testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_device" "br_testing" {
	id = "br_testing"
	macaddr = "12:34:56:78:90:ab"
	mtu = 1505
	name = "br-testing"
	ports = [
		"eth0",
		"eth1",
		"eth2.10",
		"eth2.20",
	]
	type = "bridge"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "id", "br_testing"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "macaddr", "12:34:56:78:90:ab"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "mtu", "1505"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "name", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.0", "eth0"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.1", "eth1"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.2", "eth2.10"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "ports.3", "eth2.20"),
			resource.TestCheckResourceAttr("openwrt_network_device.br_testing", "type", "bridge"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
