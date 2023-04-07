//go:build acceptance.test

package wifiiface_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestDataSourceAcceptance(t *testing.T) {
	ctx := context.Background()
	client, providerBlock := runOpenWrtServerWithWireless(
		ctx,
		*dockerPool,
		t,
	)
	options := lucirpc.Options{
		"device":  lucirpc.String("device-testing"),
		"mode":    lucirpc.String("ap"),
		"network": lucirpc.String("network-testing"),
		"ssid":    lucirpc.String("ssid-testing"),
	}
	ok, err := client.CreateSection(ctx, "wireless", "wifi-iface", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_wireless_wifi_iface" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_iface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_iface.testing", "device", "device-testing"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_iface.testing", "mode", "ap"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_iface.testing", "network", "network-testing"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_iface.testing", "ssid", "ssid-testing"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

func TestResourceAcceptance(t *testing.T) {
	ctx := context.Background()
	_, providerBlock := runOpenWrtServerWithWireless(
		ctx,
		*dockerPool,
		t,
	)

	createAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_wireless_wifi_iface" "testing" {
	device = "device-testing"
	id = "testing"
	mode = "ap"
	network = "network-testing"
	ssid = "ssid-testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "device", "device-testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "mode", "ap"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "network", "network-testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "ssid", "ssid-testing"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_wireless_wifi_iface.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_wireless_wifi_iface" "testing" {
	device = "device-testing"
	encryption = "sae"
	id = "testing"
	key = "password"
	mode = "ap"
	network = "network-testing"
	ssid = "ssid-testing"
	wpa_disable_eapol_key_retries = true
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "device", "device-testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "encryption", "sae"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "key", "password"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "mode", "ap"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "network", "network-testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "ssid", "ssid-testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_iface.testing", "wpa_disable_eapol_key_retries", "true"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
