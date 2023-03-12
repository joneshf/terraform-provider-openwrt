package openwrt_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
	"gotest.tools/v3/assert"
)

func TestOpenWrtProviderMetadataDoesNotSetVersion(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtProvider := openwrt.New("test", os.LookupEnv)
	req := provider.MetadataRequest{}
	res := &provider.MetadataResponse{}

	// When
	openWrtProvider.Metadata(ctx, req, res)

	// Then
	assert.DeepEqual(t, res.Version, "test")
}

func TestOpenWrtProviderMetadataSetsTypeName(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtProvider := openwrt.New("test", os.LookupEnv)
	req := provider.MetadataRequest{}
	res := &provider.MetadataResponse{}

	// When
	openWrtProvider.Metadata(ctx, req, res)

	// Then
	assert.DeepEqual(t, res.TypeName, "openwrt")
}

func TestOpenWrtProviderSchemaHostnameAttribute(t *testing.T) {
	attribute := "hostname"
	t.Run("exists", schemaAttributeExists(attribute))
	t.Run("is optional", schemaAttributeIsOptional(attribute))
}

func TestOpenWrtProviderSchemaPasswordAttribute(t *testing.T) {
	attribute := "password"
	t.Run("exists", schemaAttributeExists(attribute))
	t.Run("is optional", schemaAttributeIsOptional(attribute))
	t.Run("is sensitive", schemaAttributeIsSensitive(attribute))
}

func TestOpenWrtProviderSchemaPortAttribute(t *testing.T) {
	attribute := "port"
	t.Run("exists", schemaAttributeExists(attribute))
	t.Run("is optional", schemaAttributeIsOptional(attribute))
}

func TestOpenWrtProviderSchemaSchemeAttribute(t *testing.T) {
	attribute := "scheme"
	t.Run("exists", schemaAttributeExists(attribute))
	t.Run("is optional", schemaAttributeIsOptional(attribute))
}

func TestOpenWrtProviderSchemaUsernameAttribute(t *testing.T) {
	attribute := "username"
	t.Run("exists", schemaAttributeExists(attribute))
	t.Run("is optional", schemaAttributeIsOptional(attribute))
}

func TestOpenWrtProviderSchemaDoesNotUseInvalidAttributes(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtProvider := openwrt.New("test", os.LookupEnv)
	req := provider.SchemaRequest{}
	res := &provider.SchemaResponse{}
	openWrtProvider.Schema(ctx, req, res)

	// When
	diagnostics := res.Schema.Validate()

	// Then
	assert.DeepEqual(t, diagnostics, diag.Diagnostics{})
}

func schemaAttributeExists(
	attribute string,
) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New("test", os.LookupEnv)
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		_, ok := res.Schema.Attributes[attribute]
		assert.Check(t, ok)
	}
}

func schemaAttributeIsOptional(
	attribute string,
) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New("test", os.LookupEnv)
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		got := res.Schema.Attributes[attribute]
		assert.Check(t, got.IsOptional())
	}
}

func schemaAttributeIsSensitive(
	attribute string,
) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New("test", os.LookupEnv)
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		got := res.Schema.Attributes[attribute]
		assert.Check(t, got.IsSensitive())
	}
}
