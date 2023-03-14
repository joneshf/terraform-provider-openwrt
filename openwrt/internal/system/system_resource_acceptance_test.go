//go:build acceptance.test

package system_test

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

func TestSystemSystemResourceAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	hostname, port := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)

	createAndReadTestStep := resource.TestStep{
		Config: fmt.Sprintf(`
provider "openwrt" {
	hostname = %q
	password = %q
	port = %d
	username = %q
}

resource "openwrt_system_system" "this" {
	id = "cfg01e48a"
}
`,
			hostname,
			acceptancetest.Password,
			port,
			acceptancetest.Username,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_system_system.this", "id", "cfg01e48a"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "hostname", "OpenWrt"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "log_size", "64"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "timezone", "UTC"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "ttylogin", "false"),
		),
	}
	importTestStep := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_system_system.this",
	}
	updateAndReadTestStep := resource.TestStep{
		Config: fmt.Sprintf(`
provider "openwrt" {
	hostname = %q
	password = %q
	port = %d
	username = %q
}

resource "openwrt_system_system" "this" {
	hostname = "OpenWRT"
	id = "cfg01e48a"
	log_size = 64
	timezone = "UTC"
	ttylogin = false
}
`,
			hostname,
			acceptancetest.Password,
			port,
			acceptancetest.Username,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_system_system.this", "id", "cfg01e48a"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "hostname", "OpenWRT"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "log_size", "64"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "timezone", "UTC"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "ttylogin", "false"),
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
