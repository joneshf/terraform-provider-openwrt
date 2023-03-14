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
	providerBlock := acceptancetest.RunOpenWrtServerWithProviderBlock(
		ctx,
		*dockerPool,
		t,
	)

	importTestStep := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_system_system" "this" {
	id = "cfg01e48a"
}
`,
			providerBlock,
		),
		ImportState:        true,
		ImportStateId:      "cfg01e48a",
		ImportStatePersist: true,
		ResourceName:       "openwrt_system_system.this",
	}

	readTestStep := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_system_system" "this" {
	id = "cfg01e48a"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_system_system.this", "id", "cfg01e48a"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "hostname", "OpenWrt"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "log_size", "64"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "timezone", "UTC"),
			resource.TestCheckResourceAttr("openwrt_system_system.this", "ttylogin", "false"),
		),
	}

	updateAndReadTestStep := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_system_system" "this" {
	hostname = "OpenWRT"
	id = "cfg01e48a"
	log_size = 64
	timezone = "UTC"
	ttylogin = false
}
`,
			providerBlock,
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
			importTestStep,
			readTestStep,
			updateAndReadTestStep,
		},
	})
}
