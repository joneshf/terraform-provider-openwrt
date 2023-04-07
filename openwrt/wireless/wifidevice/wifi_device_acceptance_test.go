//go:build acceptance.test

package wifidevice_test

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
		"channel": lucirpc.String("auto"),
		"type":    lucirpc.String("mac80211"),
	}
	ok, err := client.CreateSection(ctx, "wireless", "wifi-device", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_wireless_wifi_device" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_device.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_device.testing", "channel", "auto"),
			resource.TestCheckResourceAttr("data.openwrt_wireless_wifi_device.testing", "type", "mac80211"),
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

resource "openwrt_wireless_wifi_device" "testing" {
	channel = "auto"
	id = "testing"
	type = "mac80211"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "channel", "auto"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "type", "mac80211"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_wireless_wifi_device.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_wireless_wifi_device" "testing" {
	band = "6g"
	channel = "auto"
	id = "testing"
	type = "mac80211"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "band", "6g"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "channel", "auto"),
			resource.TestCheckResourceAttr("openwrt_wireless_wifi_device.testing", "type", "mac80211"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
