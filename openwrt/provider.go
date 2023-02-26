package openwrt

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

const (
	hostnameAttribute           = "hostname"
	hostnameDefaultValue        = "192.168.1.1"
	hostnameEnvironmentVariable = "OPENWRT_HOSTNAME"
	hostnameHumanReadableName   = "hostname"

	passwordAttribute           = "password"
	passwordDefaultValue        = ""
	passwordEnvironmentVariable = "OPENWRT_PASSWORD"
	passwordHumanReadableName   = "password"

	portAttribute           = "port"
	portDefaultValue        = 80
	portEnvironmentVariable = "OPENWRT_PORT"
	portHumanReadableName   = "port"

	schemeAttribute           = "scheme"
	schemeDefaultValue        = "http"
	schemeEnvironmentVariable = "OPENWRT_SCHEME"
	schemeHumanReadableName   = "URI scheme"

	usernameAttribute           = "username"
	usernameDefaultValue        = "root"
	usernameEnvironmentVariable = "OPENWRT_USERNAME"
	usernameHumanReadableName   = "username"

	providerTypeName = "openwrt"
)

var (
	_ provider.Provider = &openWrtProvider{}
)

func New() provider.Provider {
	return &openWrtProvider{}
}

type openWrtProvider struct {
}

// Configure prepares an OpenWrt API client for data sources and resources.
func (p *openWrtProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	res *provider.ConfigureResponse,
) {
	tflog.Info(ctx, "Configuring OpenWrt API client")

	model := newProviderModel(ctx, req, res)
	if res.Diagnostics.HasError() {
		return
	}

	validateConfigurationKnown(ctx, model, res)
	if res.Diagnostics.HasError() {
		return
	}

	hostname := defaultStringAttributeValue(
		model.Hostname,
		hostnameEnvironmentVariable,
		hostnameDefaultValue,
	)
	password := defaultStringAttributeValue(
		model.Password,
		passwordEnvironmentVariable,
		passwordDefaultValue,
	)
	port := defaultInt64AttributeValue(
		model.Port,
		portEnvironmentVariable,
		portDefaultValue,
	)
	scheme := defaultStringAttributeValue(
		model.Scheme,
		schemeEnvironmentVariable,
		schemeDefaultValue,
	)
	username := defaultStringAttributeValue(
		model.Username,
		usernameEnvironmentVariable,
		usernameDefaultValue,
	)

	ctx = setField(ctx, hostnameAttribute, hostname)
	ctx = setField(ctx, passwordAttribute, password)
	ctx = setField(ctx, portAttribute, port)
	ctx = setField(ctx, schemeAttribute, scheme)
	ctx = setField(ctx, usernameAttribute, username)

	client := newOpenWrtClient(
		ctx,
		scheme,
		hostname,
		port,
		username,
		password,
		res,
	)
	if res.Diagnostics.HasError() {
		return
	}

	provideClient(ctx, client, res)
	if res.Diagnostics.HasError() {
		return
	}

	ctx = setField(ctx, "configure_success", true)
	tflog.Info(ctx, "Configured OpenWrt API client")
}

// DataSources defines the data sources implemented in the provider.
func (p *openWrtProvider) DataSources(
	ctx context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Metadata returns the provider type name.
func (p *openWrtProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	res *provider.MetadataResponse,
) {
	res.TypeName = providerTypeName
}

// Resources defines the resources implemented in the provider.
func (p *openWrtProvider) Resources(
	ctx context.Context,
) []func() resource.Resource {
	return []func() resource.Resource{}
}

// Schema defines the provider-level schema for configuration data.
func (p *openWrtProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	res *provider.SchemaResponse,
) {
	hostname := schema.StringAttribute{
		Description: fmt.Sprintf(
			"The %s to use. Defaults to %q.",
			hostnameHumanReadableName,
			hostnameDefaultValue,
		),
		Optional: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}

	password := schema.StringAttribute{
		Description: fmt.Sprintf(
			"The %s to use. Defaults to %q.",
			passwordHumanReadableName,
			passwordDefaultValue,
		),
		Optional:  true,
		Sensitive: true,
	}

	port := schema.Int64Attribute{
		Description: fmt.Sprintf(
			"The %s to use. Defaults to %d.",
			portHumanReadableName,
			portDefaultValue,
		),
		Optional: true,
		Validators: []validator.Int64{
			int64validator.Between(1, 65535),
		},
	}

	scheme := schema.StringAttribute{
		Description: fmt.Sprintf(
			"The %s to use. Defaults to %q.",
			schemeHumanReadableName,
			schemeDefaultValue,
		),
		Optional: true,
		Validators: []validator.String{
			stringvalidator.OneOf(
				"http",
				"https",
			),
		},
	}

	username := schema.StringAttribute{
		Description: fmt.Sprintf(
			"The %s to use. Defaults to %q.",
			usernameHumanReadableName,
			usernameDefaultValue,
		),
		Optional: true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	}

	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			hostnameAttribute: hostname,
			passwordAttribute: password,
			portAttribute:     port,
			schemeAttribute:   scheme,
			usernameAttribute: username,
		},
		Description: "Interfaces with an OpenWrt device through LuCI RPC. See https://github.com/openwrt/luci/wiki/JsonRpcHowTo#basics for setup instructions.",
	}
}

// openWrtProviderModel maps provider schema data to a Go type.
type openWrtProviderModel struct {
	Hostname types.String `tfsdk:"hostname"`
	Password types.String `tfsdk:"password"`
	Port     types.Int64  `tfsdk:"port"`
	Scheme   types.String `tfsdk:"scheme"`
	Username types.String `tfsdk:"username"`
}

type attributeInt64Default interface {
	IsNull() bool
	ValueInt64() int64
}

type attributeStringDefault interface {
	IsNull() bool
	ValueString() string
}

type attributeKnown interface {
	IsUnknown() bool
}

func defaultInt64AttributeValue(
	attribute attributeInt64Default,
	environmentVariable string,
	defaultValue int64,
) int64 {
	value := defaultValue
	variable, ok := os.LookupEnv(environmentVariable)
	if ok {
		parsed, err := strconv.Atoi(variable)
		if err != nil {
			value = int64(parsed)
		}
	}

	if !attribute.IsNull() {
		value = attribute.ValueInt64()
	}

	return value
}

func defaultStringAttributeValue(
	attribute attributeStringDefault,
	environmentVariable string,
	defaultValue string,
) string {
	value := defaultValue
	variable, ok := os.LookupEnv(environmentVariable)
	if ok {
		value = variable
	}

	if !attribute.IsNull() {
		value = attribute.ValueString()
	}

	return value
}

func newOpenWrtClient(
	ctx context.Context,
	scheme string,
	hostname string,
	port int64,
	username string,
	password string,
	res *provider.ConfigureResponse,
) *lucirpc.Client {
	tflog.Debug(ctx, "Creating OpenWrt API Client")

	client, err := lucirpc.NewClient(
		ctx,
		scheme,
		hostname,
		uint16(port),
		username,
		password,
	)
	if err != nil {
		res.Diagnostics.AddError(
			"problem creating LuCI RPC client",
			err.Error(),
		)
	}

	return client
}

func newProviderModel(
	ctx context.Context,
	req provider.ConfigureRequest,
	res *provider.ConfigureResponse,
) openWrtProviderModel {
	tflog.Debug(ctx, "Retrieving provider data from configuration")

	var config openWrtProviderModel
	diagnostics := req.Config.Get(ctx, &config)
	res.Diagnostics.Append(diagnostics...)
	return config
}

func provideClient(
	ctx context.Context,
	client *lucirpc.Client,
	res *provider.ConfigureResponse,
) {
	tflog.Debug(ctx, "Making OpenWrt client available during DataSource, and Resource type Configure methods")

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
	model openWrtProviderModel,
	res *provider.ConfigureResponse,
) {
	tflog.Debug(ctx, "Validating configuration values are known")
	validateKnown(
		model.Hostname,
		path.Root(hostnameAttribute),
		hostnameEnvironmentVariable,
		hostnameHumanReadableName,
		res,
	)
	validateKnown(
		model.Password,
		path.Root(passwordAttribute),
		passwordEnvironmentVariable,
		passwordHumanReadableName,
		res,
	)
	validateKnown(
		model.Port,
		path.Root(portAttribute),
		portEnvironmentVariable,
		portHumanReadableName,
		res,
	)
	validateKnown(
		model.Scheme,
		path.Root(schemeAttribute),
		schemeEnvironmentVariable,
		schemeHumanReadableName,
		res,
	)
	validateKnown(
		model.Username,
		path.Root(usernameAttribute),
		usernameEnvironmentVariable,
		usernameHumanReadableName,
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
			fmt.Sprintf("Unknown OpenWrt %s", humanReadableName),
			validateKnownMessage(attributePath, environmentVariable),
		)
	}
}

func validateKnownMessage(
	attributePath path.Path,
	environmentVariable string,
) string {
	pathPart := fmt.Sprintf(
		"The provider cannot create the OpenWrt API client as there is an unknown configuration value for the OpenWrt %s.",
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
