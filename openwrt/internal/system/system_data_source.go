package system

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

const (
	dataSourceTerraformType = "data_source"

	systemConLogLevelAttribute = "conloglevel"
	systemConLogLevelUCIOption = "conloglevel"

	systemCronLogLevelAttribute = "cronloglevel"
	systemCronLogLevelUCIOption = "cronloglevel"

	systemDescriptionAttribute = "description"
	systemDescriptionUCIOption = "description"

	systemHostnameAttribute = "hostname"
	systemHostnameUCIOption = "hostname"

	systemIdAttribute  = "id"
	systemIdUCISection = ".name"

	systemLogSizeAttribute = "log_size"
	systemLogSizeUCIOption = "log_size"

	systemNotesAttribute = "notes"
	systemNotesUCIOption = "notes"

	systemTimezoneAttribute = "timezone"
	systemTimezoneUCIOption = "timezone"

	systemTTYLoginAttribute = "ttylogin"
	systemTTYLoginUCIOption = "ttylogin"

	systemTypeName   = "system_system"
	systemUCIConfig  = "system"
	systemUCISection = "@system[0]"

	systemZonenameAttribute = "zonename"
	systemZonenameUCIOption = "zonename"
)

var (
	_ datasource.DataSource              = &systemDataSource{}
	_ datasource.DataSourceWithConfigure = &systemDataSource{}
)

func NewSystemDataSource() datasource.DataSource {
	return &systemDataSource{}
}

type systemDataSource struct {
	client       lucirpc.Client
	fullTypeName string
}

// Configure prepares the data source.
func (d *systemDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	res *datasource.ConfigureResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Configuring %s Data Source", d.fullTypeName))
	if req.ProviderData == nil {
		tflog.Debug(ctx, "No provider data")
		return
	}

	client, diagnostics := newUCIClient(req)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	d.client = *client
}

// Metadata sets the data source name.
func (d *systemDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	res *datasource.MetadataResponse,
) {
	fullTypeName := fmt.Sprintf("%s_%s", req.ProviderTypeName, systemTypeName)
	d.fullTypeName = fullTypeName
	res.TypeName = fullTypeName
}

// Read prepares the data source.
func (d *systemDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	res *datasource.ReadResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Reading %s data source", d.fullTypeName))

	section, diagnostics := getSection(ctx, d.client)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	var model systemModel
	ctx, model.ConLogLevel, diagnostics = getOptionInt64(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemConLogLevelAttribute), systemConLogLevelUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.CronLogLevel, diagnostics = getOptionInt64(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemCronLogLevelAttribute), systemCronLogLevelUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Description, diagnostics = getOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemDescriptionAttribute), systemDescriptionUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Hostname, diagnostics = getOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemHostnameAttribute), systemHostnameUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.LogSize, diagnostics = getOptionInt64(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemLogSizeAttribute), systemLogSizeUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Notes, diagnostics = getOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemNotesAttribute), systemNotesUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Timezone, diagnostics = getOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemTimezoneAttribute), systemTimezoneUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.TTYLogin, diagnostics = getOptionBool(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemTTYLoginAttribute), systemTTYLoginUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Zonename, diagnostics = getOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemZonenameAttribute), systemZonenameUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Id, diagnostics = getMetadataString(ctx, d.fullTypeName, dataSourceTerraformType, section, systemIdUCISection)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Setting the %s data source state", d.fullTypeName))
	diagnostics = res.State.Set(ctx, model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

// Schema prepares the data source.
func (d *systemDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	res *datasource.SchemaResponse,
) {
	conLogLevel := schema.Int64Attribute{
		Description: "The maximum log level for kernel messages to be logged to the console.",
		Optional:    true,
	}
	cronLogLevel := schema.Int64Attribute{
		Description: "The minimum level for cron messages to be logged to syslog.",
		Optional:    true,
	}
	description := schema.StringAttribute{
		Description: "The hostname for the system.",
		Optional:    true,
	}
	hostname := schema.StringAttribute{
		Description: "A short single-line description for the system.",
		Optional:    true,
	}
	id := schema.StringAttribute{
		Computed:    true,
		Description: "Placeholder identifier attribute.",
	}
	logSize := schema.Int64Attribute{
		Description: "Size of the file based log buffer in KiB.",
		Optional:    true,
	}
	notes := schema.StringAttribute{
		Description: "Multi-line free-form text about the system.",
		Optional:    true,
	}
	timezone := schema.StringAttribute{
		Description: "The POSIX.1 time zone string. This has no corresponding value in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		Optional:    true,
	}
	ttyLogin := schema.BoolAttribute{
		Description: "Require authentication for local users to log in the system.",
		Optional:    true,
	}
	zonename := schema.StringAttribute{
		Description: "The IANA/Olson time zone string. This corresponds to \"Timezone\" in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		Optional:    true,
	}

	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			systemConLogLevelAttribute:  conLogLevel,
			systemCronLogLevelAttribute: cronLogLevel,
			systemDescriptionAttribute:  description,
			systemHostnameAttribute:     hostname,
			systemIdAttribute:           id,
			systemLogSizeAttribute:      logSize,
			systemNotesAttribute:        notes,
			systemTimezoneAttribute:     timezone,
			systemTTYLoginAttribute:     ttyLogin,
			systemZonenameAttribute:     zonename,
		},
		Description: "Provides system data about an OpenWrt device",
	}
}

func getMetadataString(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section map[string]json.RawMessage,
	key string,
) (context.Context, types.String, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.StringNull()
	raw, ok := section[key]
	if !ok {
		return ctx, result, diagnostics
	}

	var value string
	err := json.Unmarshal(raw, &value)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("unable to parse metadata: %q", key),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	result = types.StringValue(value)
	ctx = logSetFieldString(ctx, fullTypeName, terraformType, key, result)
	return ctx, result, diagnostics
}

func getOptionBool(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section map[string]json.RawMessage,
	attribute path.Path,
	option string,
) (context.Context, types.Bool, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.BoolNull()
	raw, ok := section[option]
	if !ok {
		return ctx, result, diagnostics
	}

	// Booleans in UCI can be any number of things:
	// - True: "1", "yes", "on", "true", "enabled"
	// - False: "0", "no", "off", "false", "disabled"
	// We try to parse on of these out of the string.
	var boolish string
	err := json.Unmarshal(raw, &boolish)
	if err != nil {
		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	switch boolish {
	case "1", "yes", "on", "true", "enabled":
		result = types.BoolValue(true)

	case "0", "no", "off", "false", "disabled":
		result = types.BoolValue(false)

	default:
		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("Unexpected value for option: %q", option),
			fmt.Sprintf(`expected one of "1", "yes", "on", "true", "enabled", "0", "no", "off", "false", or "disabled"; got: %q`, boolish),
		)
		return ctx, result, diagnostics
	}

	ctx = logSetFieldBool(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, diagnostics
}

func getOptionInt64(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section map[string]json.RawMessage,
	attribute path.Path,
	option string,
) (context.Context, types.Int64, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.Int64Null()
	raw, ok := section[option]
	if !ok {
		return ctx, result, diagnostics
	}

	// Integers in UCI are stored as strtings.
	// We have to unmarshall first, then parse the string.
	var intish string
	err := json.Unmarshal(raw, &intish)
	if err != nil {
		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	value, err := strconv.Atoi(intish)
	if err != nil {
		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to convert option: %q to a string", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	result = types.Int64Value(int64(value))
	ctx = logSetFieldInt64(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, diagnostics
}

func getOptionString(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section map[string]json.RawMessage,
	attribute path.Path,
	option string,
) (context.Context, types.String, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.StringNull()
	raw, ok := section[option]
	if !ok {
		return ctx, result, diagnostics
	}

	var value string
	err := json.Unmarshal(raw, &value)
	if err != nil {
		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	result = types.StringValue(value)
	ctx = logSetFieldString(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, diagnostics
}

func getSection(
	ctx context.Context,
	client lucirpc.Client,
) (map[string]json.RawMessage, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	section, err := client.GetSection(ctx, systemUCIConfig, systemUCISection)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("problem getting %s.%s section", systemUCIConfig, systemUCISection),
			err.Error(),
		)
		return map[string]json.RawMessage{}, diagnostics
	}

	return section, diagnostics
}

func logSetFieldBool(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	key string,
	value logValueBool,
) context.Context {
	ctx = tflog.SetField(ctx, fmt.Sprintf("%s_%s_%s", fullTypeName, terraformType, key), value.ValueBool())
	return ctx
}

func logSetFieldInt64(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	key string,
	value logValueInt64,
) context.Context {
	ctx = tflog.SetField(ctx, fmt.Sprintf("%s_%s_%s", fullTypeName, terraformType, key), value.ValueInt64())
	return ctx
}

func logSetFieldString(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	key string,
	value logValueString,
) context.Context {
	ctx = tflog.SetField(ctx, fmt.Sprintf("%s_%s_%s", fullTypeName, terraformType, key), value.ValueString())
	return ctx
}

type logValueBool interface {
	ValueBool() bool
}

type logValueInt64 interface {
	ValueInt64() int64
}

type logValueString interface {
	ValueString() string
}

func newUCIClient(
	req datasource.ConfigureRequest,
) (*lucirpc.Client, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	client, ok := req.ProviderData.(*lucirpc.Client)
	if !ok {
		diagnostics.AddError(
			"OpenWrt provider not configured correctly",
			"Expected UCI tree, but one was not provided. This is a problem with the provider implementation. Please report this to https://github.com/joneshf/terraform-provider-openwrt",
		)
		return nil, diagnostics
	}

	return client, diagnostics
}

type systemModel struct {
	ConLogLevel  types.Int64  `tfsdk:"conloglevel"`
	CronLogLevel types.Int64  `tfsdk:"cronloglevel"`
	Description  types.String `tfsdk:"description"`
	Hostname     types.String `tfsdk:"hostname"`
	Id           types.String `tfsdk:"id"`
	LogSize      types.Int64  `tfsdk:"log_size"`
	Notes        types.String `tfsdk:"notes"`
	Timezone     types.String `tfsdk:"timezone"`
	TTYLogin     types.Bool   `tfsdk:"ttylogin"`
	Zonename     types.String `tfsdk:"zonename"`
}
