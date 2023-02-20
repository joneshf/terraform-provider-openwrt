package openwrt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider = &openWRTProvider{}
)

func New() provider.Provider {
	return &openWRTProvider{}
}

type openWRTProvider struct {
}

// Configure prepares an OpenWRT API client for data sources and resources.
func (p *openWRTProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	res *provider.ConfigureResponse,
) {
}

// DataSources defines the data sources implemented in the provider.
func (p *openWRTProvider) DataSources(
	ctx context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Metadata returns the provider type name.
func (p *openWRTProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	res *provider.MetadataResponse,
) {
}

// Resources defines the resources implemented in the provider.
func (p *openWRTProvider) Resources(
	ctx context.Context,
) []func() resource.Resource {
	return []func() resource.Resource{}
}

// Schema defines the provider-level schema for configuration data.
func (p *openWRTProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	res *provider.SchemaResponse,
) {
	// res.Schema = schema.Schema{}
}
