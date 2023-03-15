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

func TestSystemSystemDataSourceAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	providerBlock := acceptancetest.RunOpenWrtServerWithProviderBlock(
		ctx,
		*dockerPool,
		t,
	)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_system_system" "this" {
	id = "cfg01e48a"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "id", "cfg01e48a"),
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "hostname", "OpenWrt"),
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "log_size", "64"),
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "timezone", "UTC"),
			resource.TestCheckResourceAttr("data.openwrt_system_system.this", "ttylogin", "false"),
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"openwrt": providerserver.NewProtocol6WithError(openwrt.New("test", os.LookupEnv)),
		},
		Steps: []resource.TestStep{
			readDataSource,
		},
	})
}
