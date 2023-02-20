package openwrt

import (
	"context"
	"fmt"
	"os"

	"github.com/digineo/go-uci"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	configurationDirectoryAttribute           = "configuration_directory"
	configurationDirectoryDefaultValue        = "/etc/config"
	configurationDirectoryEnvironmentVariable = "OPENWRT_CONFIGURATION_DIRECTORY"
	configurationDirectoryHumanReadableName   = "configuration directory"

	providerTypeName = "openwrt"
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
	tflog.Info(ctx, "Configuring OpenWRT API client")

	model := newProviderModel(ctx, req, res)
	if res.Diagnostics.HasError() {
		return
	}

	validateConfigurationKnown(ctx, model, res)
	if res.Diagnostics.HasError() {
		return
	}

	configurationDirectory := defaultAttributeValue(
		model.ConfigurationDirectory,
		configurationDirectoryEnvironmentVariable,
		configurationDirectoryDefaultValue,
	)

	ctx = setField(ctx, configurationDirectoryAttribute, configurationDirectory)
	client := newOpenWRTClient(ctx, configurationDirectory, res)
	if res.Diagnostics.HasError() {
		return
	}

	provideClient(ctx, client, res)
	if res.Diagnostics.HasError() {
		return
	}

	ctx = setField(ctx, "configure_success", true)
	tflog.Info(ctx, "Configured HashiCups client")
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
	res.TypeName = providerTypeName
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
	configurationDirectory := schema.StringAttribute{
		Description: fmt.Sprintf(
			"The configuration directory to use. Defaults to %q.",
			configurationDirectoryDefaultValue,
		),
		Optional: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}

	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			configurationDirectoryAttribute: configurationDirectory,
		},
		Description: "Interfaces with an OpenWRT device through UCI.",
	}
}

// openWRTProviderModel maps provider schema data to a Go type.
type openWRTProviderModel struct {
	ConfigurationDirectory types.String `tfsdk:"configuration_directory"`
}

type attributeDefault interface {
	IsNull() bool
	ValueString() string
}

type attributeKnown interface {
	IsUnknown() bool
}

func defaultAttributeValue(
	attribute attributeDefault,
	environmentVariable string,
	defaultValue string,
) string {
	configurationDirectory := os.Getenv(environmentVariable)
	if !attribute.IsNull() {
		configurationDirectory = attribute.ValueString()
	}

	return configurationDirectory
}

func newOpenWRTClient(
	ctx context.Context,
	configurationDirectory string,
	res *provider.ConfigureResponse,
) uci.Tree {
	tflog.Debug(ctx, "Creating OpenWRT API Client")

	client := uci.NewTree(configurationDirectory)

	return client
}

func newProviderModel(
	ctx context.Context,
	req provider.ConfigureRequest,
	res *provider.ConfigureResponse,
) openWRTProviderModel {
	tflog.Debug(ctx, "Retrieving provider data from configuration")

	var config openWRTProviderModel
	diagnostics := req.Config.Get(ctx, &config)
	res.Diagnostics.Append(diagnostics...)
	return config
}

func provideClient(
	ctx context.Context,
	client uci.Tree,
	res *provider.ConfigureResponse,
) {
	tflog.Debug(ctx, "Making OpenWRT client available during DataSource, and Resource type Configure methods")

	res.DataSourceData = client
	res.ResourceData = client
}

func setField(
	ctx context.Context,
	key string,
	value any,
) context.Context {
	field := fmt.Sprintf("%s_%s", providerTypeName, key)
	ctx = tflog.SetField(ctx, field, value)
	return ctx
}

func validateConfigurationKnown(
	ctx context.Context,
	model openWRTProviderModel,
	res *provider.ConfigureResponse,
) {
	tflog.Debug(ctx, "Validating configuration values are known")
	validateKnown(
		model.ConfigurationDirectory,
		path.Root(configurationDirectoryAttribute),
		configurationDirectoryEnvironmentVariable,
		configurationDirectoryHumanReadableName,
		res,
	)
}

func validateKnown(
	attribute attributeKnown,
	attributePath path.Path,
	environmentVariable string,
	humanReadableName string,
	res *provider.ConfigureResponse,
) {
	if attribute.IsUnknown() {
		res.Diagnostics.AddAttributeError(
			attributePath,
			fmt.Sprintf("Unknown OpenWRT %s", humanReadableName),
			validateKnownMessage(attributePath, environmentVariable),
		)
	}
}

func validateKnownMessage(
	attributePath path.Path,
	environmentVariable string,
) string {
	pathPart := fmt.Sprintf(
		"The provider cannot create the OpenWRT API client as there is an unknown configuration value for the OpenWRT %s.",
		attributePath.String(),
	)
	environmentVariablePart := fmt.Sprintf(
		"Either target apply the source of the value first, set the value statically in the configuration, or use the %s environment variable.",
		environmentVariable,
	)
	return fmt.Sprintf(
		"%s %s",
		pathPart,
		environmentVariablePart,
	)
}
