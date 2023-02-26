//go:build acceptance.test

package openwrt_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
	"gotest.tools/v3/assert"
)

func TestOpenWrtProviderConfigureConnectsWithoutError(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtProvider := openwrt.New()
	schemaReq := provider.SchemaRequest{}
	schemaRes := &provider.SchemaResponse{}
	openWrtProvider.Schema(ctx, schemaReq, schemaRes)
	config := tfsdk.Config{
		Schema: schemaRes.Schema,
		Raw: tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"hostname": tftypes.String,
					"password": tftypes.String,
					"port":     tftypes.Number,
					"scheme":   tftypes.String,
					"username": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"hostname": tftypes.NewValue(tftypes.String, "localhost"),
				"password": tftypes.NewValue(tftypes.String, ""),
				"port":     tftypes.NewValue(tftypes.Number, 8080),
				"scheme":   tftypes.NewValue(tftypes.String, "http"),
				"username": tftypes.NewValue(tftypes.String, "root"),
			},
		),
	}
	req := provider.ConfigureRequest{
		Config: config,
	}
	res := &provider.ConfigureResponse{}

	// When
	openWrtProvider.Configure(ctx, req, res)

	// Then
	assert.DeepEqual(t, res.Diagnostics, diag.Diagnostics{})
}
