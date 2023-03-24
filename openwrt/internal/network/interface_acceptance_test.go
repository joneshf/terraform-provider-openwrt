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

func TestNetworkInterfaceDataSourceAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	client, providerBlock := acceptancetest.AuthenticatedClientWithProviderBlock(
		ctx,
		*dockerPool,
		t,
	)
	options := lucirpc.Options{
		"device":  lucirpc.String("br-testing"),
		"ipaddr":  lucirpc.String("192.168.3.1"),
		"netmask": lucirpc.String("255.255.255.0"),
		"proto":   lucirpc.String("static"),
	}
	ok, err := client.CreateSection(ctx, "network", "interface", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_network_interface" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "ipaddr", "192.168.3.1"),
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "netmask", "255.255.255.0"),
			resource.TestCheckResourceAttr("data.openwrt_network_interface.testing", "proto", "static"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

func TestNetworkInterfaceResourceAcceptance(t *testing.T) {
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

resource "openwrt_network_interface" "testing" {
	device = "br-testing"
	id = "testing"
	ipaddr = "192.168.3.1"
	netmask = "255.255.255.0"
	proto = "static"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "ipaddr", "192.168.3.1"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "netmask", "255.255.255.0"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "proto", "static"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_network_interface.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_network_interface" "testing" {
	device = "br-testing"
	dns = [
		"9.9.9.9",
		"1.1.1.1",
	]
	id = "testing"
	ipaddr = "192.168.3.1"
	macaddr = "12:34:56:78:90:ab"
	mtu = 1505
	netmask = "255.255.255.0"
	proto = "static"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "device", "br-testing"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "dns.1", "1.1.1.1"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "dns.0", "9.9.9.9"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "ipaddr", "192.168.3.1"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "macaddr", "12:34:56:78:90:ab"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "mtu", "1505"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "netmask", "255.255.255.0"),
			resource.TestCheckResourceAttr("openwrt_network_interface.testing", "proto", "static"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
