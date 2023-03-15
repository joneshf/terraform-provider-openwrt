//go:build acceptance.test

package network_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
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
	options := map[string]json.RawMessage{
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
