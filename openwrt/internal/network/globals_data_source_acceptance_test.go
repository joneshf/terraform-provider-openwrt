//go:build acceptance.test

package network_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
	"gotest.tools/v3/assert"
)

func TestNetworkGlobalsDataSourceAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	openWrt, hostname, port := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	defer openWrt.Close()
	client, err := lucirpc.NewClient(ctx, acceptancetest.Scheme, hostname, port, acceptancetest.Username, acceptancetest.Password)
	assert.NilError(t, err)
	options := map[string]json.RawMessage{
		"packet_steering": json.RawMessage("false"),
		"ula_prefix":      json.RawMessage(`"fd12:3456:789a::/48"`),
	}
	ok, err := client.CreateSection(ctx, "network", "globals", "globals", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"openwrt": providerserver.NewProtocol6WithError(openwrt.New("test")),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "openwrt" {
	hostname = %q
	password = %q
	port = %d
	username = %q
}

data "openwrt_network_globals" "this" {
	id = "globals"
}
`,
					hostname,
					acceptancetest.Password,
					port,
					acceptancetest.Username,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.openwrt_network_globals.this", "id", "globals"),
					resource.TestCheckResourceAttr("data.openwrt_network_globals.this", "packet_steering", "false"),
					resource.TestCheckResourceAttr("data.openwrt_network_globals.this", "ula_prefix", "fd12:3456:789a::/48"),
				),
			},
		},
	})
}
