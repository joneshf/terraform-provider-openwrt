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

func TestOpenWrtProviderConfigureDoesNotErrorWithNoConfiguration(t *testing.T) {
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
					"configuration_directory": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"configuration_directory": tftypes.NewValue(tftypes.String, ""),
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

func TestOpenWrtProviderMetadataDoesNotSetVersion(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtProvider := openwrt.New()
	req := provider.MetadataRequest{}
	res := &provider.MetadataResponse{}

	// When
	openWrtProvider.Metadata(ctx, req, res)

	// Then
	assert.DeepEqual(t, res.Version, "")
}

func TestOpenWrtProviderMetadataSetsTypeName(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtProvider := openwrt.New()
	req := provider.MetadataRequest{}
	res := &provider.MetadataResponse{}

	// When
	openWrtProvider.Metadata(ctx, req, res)

	// Then
	assert.DeepEqual(t, res.TypeName, "openwrt")
}

func TestOpenWrtProviderSchemaHasOptionalConfigurationDirectory(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtProvider := openwrt.New()
	req := provider.SchemaRequest{}
	res := &provider.SchemaResponse{}

	// When
	openWrtProvider.Schema(ctx, req, res)

	// Then
	attributes := res.Schema.Attributes
	configurationDirectory, ok := attributes["configuration_directory"]
	assert.Check(t, ok)
	assert.Check(t, configurationDirectory.IsOptional())
}

func TestOpenWrtProviderSchemaDoesNotUseInvalidAttributes(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtProvider := openwrt.New()
	req := provider.SchemaRequest{}
	res := &provider.SchemaResponse{}
	openWrtProvider.Schema(ctx, req, res)

	// When
	diagnostics := res.Schema.Validate()

	// Then
	assert.DeepEqual(t, diagnostics, diag.Diagnostics{})
}
