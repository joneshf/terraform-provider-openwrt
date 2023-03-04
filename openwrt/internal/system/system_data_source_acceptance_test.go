//go:build acceptance.test

package system_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
	"github.com/ory/dockertest/v3"
)

var (
	dockerPool *dockertest.Pool
)

func TestMain(m *testing.M) {
	var (
		code     int
		err      error
		tearDown func()
	)
	ctx := context.Background()
	tearDown, dockerPool, err = acceptancetest.Setup(ctx, m)
	defer func() {
		tearDown()
		os.Exit(code)
	}()
	if err != nil {
		fmt.Printf("Problem setting up tests: %s", err)
		code = 1
		return
	}

	log.Printf("Running tests")
	code = m.Run()
}

func TestSystemSystemDataSourceAcceptance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	openWrt, hostname, port := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	defer openWrt.Close()

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

data "openwrt_system_system" "this" {
}
`,
					hostname,
					acceptancetest.Password,
					port,
					acceptancetest.Username,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "id", "cfg01e48a"),
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "hostname", "OpenWrt"),
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "log_size", "64"),
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "timezone", "UTC"),
					resource.TestCheckResourceAttr("data.openwrt_system_system.this", "ttylogin", "false"),
				),
			},
		},
	})
}
