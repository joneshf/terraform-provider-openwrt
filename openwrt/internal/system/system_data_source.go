package system

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
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

	client, diagnostics := lucirpcglue.NewClient(req)
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

	section, diagnostics := lucirpcglue.GetSection(ctx, d.client, systemUCIConfig, systemUCISection)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	var model systemModel
	ctx, model.ConLogLevel, diagnostics = lucirpcglue.GetOptionInt64(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemConLogLevelAttribute), systemConLogLevelUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.CronLogLevel, diagnostics = lucirpcglue.GetOptionInt64(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemCronLogLevelAttribute), systemCronLogLevelUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Description, diagnostics = lucirpcglue.GetOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemDescriptionAttribute), systemDescriptionUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Hostname, diagnostics = lucirpcglue.GetOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemHostnameAttribute), systemHostnameUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.LogSize, diagnostics = lucirpcglue.GetOptionInt64(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemLogSizeAttribute), systemLogSizeUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Notes, diagnostics = lucirpcglue.GetOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemNotesAttribute), systemNotesUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Timezone, diagnostics = lucirpcglue.GetOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemTimezoneAttribute), systemTimezoneUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.TTYLogin, diagnostics = lucirpcglue.GetOptionBool(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemTTYLoginAttribute), systemTTYLoginUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Zonename, diagnostics = lucirpcglue.GetOptionString(ctx, d.fullTypeName, dataSourceTerraformType, section, path.Root(systemZonenameAttribute), systemZonenameUCIOption)
	res.Diagnostics.Append(diagnostics...)
	ctx, model.Id, diagnostics = lucirpcglue.GetMetadataString(ctx, d.fullTypeName, dataSourceTerraformType, section, systemIdUCISection)
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
