package openwrt_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
	"gotest.tools/v3/assert"
)

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

func TestOpenWrtProviderSchemaHostnameAttribute(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		_, ok := res.Schema.Attributes["hostname"]
		assert.Check(t, ok)
	})

	t.Run("is optional", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		got := res.Schema.Attributes["hostname"]
		assert.Check(t, got.IsOptional())
	})

}

func TestOpenWrtProviderSchemaPasswordAttribute(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		_, ok := res.Schema.Attributes["password"]
		assert.Check(t, ok)
	})

	t.Run("is optional", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		got := res.Schema.Attributes["password"]
		assert.Check(t, got.IsOptional())
	})

	t.Run("is sensitive", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		got := res.Schema.Attributes["password"]
		assert.Check(t, got.IsSensitive())
	})
}

func TestOpenWrtProviderSchemaPortAttribute(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		_, ok := res.Schema.Attributes["port"]
		assert.Check(t, ok)
	})

	t.Run("is optional", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		got := res.Schema.Attributes["port"]
		assert.Check(t, got.IsOptional())
	})
}

func TestOpenWrtProviderSchemaSchemeAttribute(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		_, ok := res.Schema.Attributes["scheme"]
		assert.Check(t, ok)
	})

	t.Run("is optional", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		got := res.Schema.Attributes["scheme"]
		assert.Check(t, got.IsOptional())
	})
}

func TestOpenWrtProviderSchemaUsernameAttribute(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		_, ok := res.Schema.Attributes["username"]
		assert.Check(t, ok)
	})

	t.Run("is optional", func(t *testing.T) {
		// Given
		ctx := context.Background()
		openWrtProvider := openwrt.New()
		req := provider.SchemaRequest{}
		res := &provider.SchemaResponse{}

		// When
		openWrtProvider.Schema(ctx, req, res)

		// Then
		got := res.Schema.Attributes["username"]
		assert.Check(t, got.IsOptional())
	})
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
