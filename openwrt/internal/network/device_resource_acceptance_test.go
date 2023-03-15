//go:build acceptance.test

package network_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
)

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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"openwrt": providerserver.NewProtocol6WithError(openwrt.New("test", os.LookupEnv)),
		},
		Steps: []resource.TestStep{
			createAndReadResource,
			importValidation,
			updateAndReadResource,
		},
	})
}
