package system

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

const (
	providerTypeName = "openwrt"

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

	systemFullTypeName = fmt.Sprintf("%s_%s", providerTypeName, systemTypeName)
)

func NewSystemDataSource() datasource.DataSource {
	return &systemDataSource{}
}

type systemDataSource struct {
	client lucirpc.Client
}

// Configure prepares the data source.
func (d *systemDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	res *datasource.ConfigureResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Configuring %s Data Source", systemFullTypeName))
	if req.ProviderData == nil {
		tflog.Debug(ctx, "No provider data")
		return
	}

	client := newUCIClient(req, res)
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
	res.TypeName = systemFullTypeName
}

// Read prepares the data source.
func (d *systemDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	res *datasource.ReadResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Reading %s data source", systemFullTypeName))

	section := getSection(ctx, d.client, res)
	if res.Diagnostics.HasError() {
		return
	}

	var model systemModel
	ctx, model.ConLogLevel = getOptionInt64(ctx, section, path.Root(systemConLogLevelAttribute), systemConLogLevelUCIOption, res)
	ctx, model.CronLogLevel = getOptionInt64(ctx, section, path.Root(systemCronLogLevelAttribute), systemCronLogLevelUCIOption, res)
	ctx, model.Description = getOptionString(ctx, section, path.Root(systemDescriptionAttribute), systemDescriptionUCIOption, res)
	ctx, model.Hostname = getOptionString(ctx, section, path.Root(systemHostnameAttribute), systemHostnameUCIOption, res)
	ctx, model.LogSize = getOptionInt64(ctx, section, path.Root(systemLogSizeAttribute), systemLogSizeUCIOption, res)
	ctx, model.Notes = getOptionString(ctx, section, path.Root(systemNotesAttribute), systemNotesUCIOption, res)
	ctx, model.Timezone = getOptionString(ctx, section, path.Root(systemTimezoneAttribute), systemTimezoneUCIOption, res)
	ctx, model.TTYLogin = getOptionBool(ctx, section, path.Root(systemTTYLoginAttribute), systemTTYLoginUCIOption, res)
	ctx, model.Zonename = getOptionString(ctx, section, path.Root(systemZonenameAttribute), systemZonenameUCIOption, res)
	ctx, model.Id = getMetadataString(ctx, section, systemIdUCISection, res)
	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Setting the %s data source state", systemFullTypeName))
	diagnostics := res.State.Set(ctx, model)
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
	section map[string]json.RawMessage,
	key string,
	res *datasource.ReadResponse,
) (context.Context, types.String) {
	result := types.StringNull()
	raw, ok := section[key]
	if !ok {
		return ctx, result
	}

	var value string
	err := json.Unmarshal(raw, &value)
	if err != nil {
		res.Diagnostics.AddError(
			fmt.Sprintf("unable to parse metadata: %q", key),
			err.Error(),
		)
		return ctx, result
	}

	result = types.StringValue(value)
	ctx = logSetFieldString(ctx, key, result)
	return ctx, result
}

func getOptionBool(
	ctx context.Context,
	section map[string]json.RawMessage,
	attribute path.Path,
	option string,
	res *datasource.ReadResponse,
) (context.Context, types.Bool) {
	result := types.BoolNull()
	raw, ok := section[option]
	if !ok {
		return ctx, result
	}

	// Booleans in UCI can be any number of things:
	// - True: "1", "yes", "on", "true", "enabled"
	// - False: "0", "no", "off", "false", "disabled"
	// We try to parse on of these out of the string.
	var boolish string
	err := json.Unmarshal(raw, &boolish)
	if err != nil {
		res.Diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result
	}

	switch boolish {
	case "1", "yes", "on", "true", "enabled":
		result = types.BoolValue(true)

	case "0", "no", "off", "false", "disabled":
		result = types.BoolValue(false)

	default:
		res.Diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("Unexpected value for option: %q", option),
			fmt.Sprintf(`expected one of "1", "yes", "on", "true", "enabled", "0", "no", "off", "false", or "disabled"; got: %q`, boolish),
		)
		return ctx, result
	}

	ctx = logSetFieldBool(ctx, option, result)
	return ctx, result
}

func getOptionInt64(
	ctx context.Context,
	section map[string]json.RawMessage,
	attribute path.Path,
	option string,
	res *datasource.ReadResponse,
) (context.Context, types.Int64) {
	result := types.Int64Null()
	raw, ok := section[option]
	if !ok {
		return ctx, result
	}

	// Integers in UCI are stored as strtings.
	// We have to unmarshall first, then parse the string.
	var intish string
	err := json.Unmarshal(raw, &intish)
	if err != nil {
		res.Diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result
	}

	value, err := strconv.Atoi(intish)
	if err != nil {
		res.Diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to convert option: %q to a string", option),
			err.Error(),
		)
		return ctx, result
	}

	result = types.Int64Value(int64(value))
	ctx = logSetFieldInt64(ctx, option, result)
	return ctx, result
}

func getOptionString(
	ctx context.Context,
	section map[string]json.RawMessage,
	attribute path.Path,
	option string,
	res *datasource.ReadResponse,
) (context.Context, types.String) {
	result := types.StringNull()
	raw, ok := section[option]
	if !ok {
		return ctx, result
	}

	var value string
	err := json.Unmarshal(raw, &value)
	if err != nil {
		res.Diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result
	}

	result = types.StringValue(value)
	ctx = logSetFieldString(ctx, option, result)
	return ctx, result
}

func getSection(
	ctx context.Context,
	client lucirpc.Client,
	res *datasource.ReadResponse,
) map[string]json.RawMessage {
	section, err := client.GetSection(ctx, systemUCIConfig, systemUCISection)
	if err != nil {
		res.Diagnostics.AddError(
			fmt.Sprintf("problem getting %s.%s section", systemUCIConfig, systemUCISection),
			err.Error(),
		)
		return map[string]json.RawMessage{}
	}

	return section
}

func logSetFieldBool(
	ctx context.Context,
	key string,
	value logValueBool,
) context.Context {
	ctx = tflog.SetField(ctx, fmt.Sprintf("%s_data_source_%s", systemFullTypeName, key), value.ValueBool())
	return ctx
}

func logSetFieldInt64(
	ctx context.Context,
	key string,
	value logValueInt64,
) context.Context {
	ctx = tflog.SetField(ctx, fmt.Sprintf("%s_data_source_%s", systemFullTypeName, key), value.ValueInt64())
	return ctx
}

func logSetFieldString(
	ctx context.Context,
	key string,
	value logValueString,
) context.Context {
	ctx = tflog.SetField(ctx, fmt.Sprintf("%s_data_source_%s", systemFullTypeName, key), value.ValueString())
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
	res *datasource.ConfigureResponse,
) *lucirpc.Client {
	client, ok := req.ProviderData.(*lucirpc.Client)
	if !ok {
		res.Diagnostics.AddError(
			"OpenWrt provider not configured correctly",
			"Expected UCI tree, but one was not provided. This is a problem with the provider implementation. Please report this to https://github.com/joneshf/terraform-provider-openwrt",
		)
		return nil
	}

	return client
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
