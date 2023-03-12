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

func TestNetworkGlobalsResourceAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	openWrt, hostname, port := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	defer openWrt.Close()

	createAndReadTestStep := resource.TestStep{
		Config: fmt.Sprintf(`
provider "openwrt" {
	hostname = %q
	password = %q
	port = %d
	username = %q
}

resource "openwrt_network_globals" "this" {
	id = "globals"
}
`,
			hostname,
			acceptancetest.Password,
			port,
			acceptancetest.Username,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_globals.this", "id", "globals"),
			resource.TestCheckNoResourceAttr("openwrt_network_globals.this", "network_steering"),
			resource.TestCheckNoResourceAttr("openwrt_network_globals.this", "ula_prefix"),
		),
	}
	importTestStep := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_network_globals.this",
	}
	updateAndReadTestStep := resource.TestStep{
		Config: fmt.Sprintf(`
provider "openwrt" {
	hostname = %q
	password = %q
	port = %d
	username = %q
}

resource "openwrt_network_globals" "this" {
	id = "globals"
	packet_steering = false
	ula_prefix = "fd12:3456:789a::/48"
}
`,
			hostname,
			acceptancetest.Password,
			port,
			acceptancetest.Username,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_network_globals.this", "id", "globals"),
			resource.TestCheckResourceAttr("openwrt_network_globals.this", "packet_steering", "false"),
			resource.TestCheckResourceAttr("openwrt_network_globals.this", "ula_prefix", "fd12:3456:789a::/48"),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"openwrt": providerserver.NewProtocol6WithError(openwrt.New("test", os.LookupEnv)),
		},
		Steps: []resource.TestStep{
			createAndReadTestStep,
			importTestStep,
			updateAndReadTestStep,
		},
	})
}
