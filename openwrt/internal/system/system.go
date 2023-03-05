package system

import (
	"context"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	ReadOnly AttributeExistence = iota
	NoValidation
	Optional
	Required

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

	systemTypeName      = "system_system"
	systemUCIConfig     = "system"
	systemUCISection    = "@system[0]"
	systemUCISystemType = "system"

	systemZonenameAttribute = "zonename"
	systemZonenameUCIOption = "zonename"
)

var (
	systemConLogLevelSchemaAttribute = int64SchemaAttribute{
		Description:       "The maximum log level for kernel messages to be logged to the console.",
		ResourceExistence: NoValidation,
	}

	systemCronLogLevelSchemaAttribute = int64SchemaAttribute{
		Description:       "The minimum level for cron messages to be logged to syslog.",
		ResourceExistence: NoValidation,
	}

	systemDescriptionSchemaAttribute = stringSchemaAttribute{
		Description:       "The hostname for the system.",
		ResourceExistence: NoValidation,
	}

	systemHostnameSchemaAttribute = stringSchemaAttribute{
		Description:       "A short single-line description for the system.",
		ResourceExistence: NoValidation,
	}

	systemIdSchemaAttribute = stringSchemaAttribute{
		DataSourceExistence: Required,
		Description:         "Placeholder identifier attribute.",
		ResourceExistence:   Required,
	}

	systemLogSizeSchemaAttribute = int64SchemaAttribute{
		Description:       "Size of the file based log buffer in KiB.",
		ResourceExistence: NoValidation,
	}

	systemNotesSchemaAttribute = stringSchemaAttribute{
		Description:       "Multi-line free-form text about the system.",
		ResourceExistence: NoValidation,
	}

	systemSchemaAttributes = map[string]schemaAttribute{
		systemConLogLevelAttribute:  systemConLogLevelSchemaAttribute,
		systemCronLogLevelAttribute: systemCronLogLevelSchemaAttribute,
		systemDescriptionAttribute:  systemDescriptionSchemaAttribute,
		systemHostnameAttribute:     systemHostnameSchemaAttribute,
		systemIdAttribute:           systemIdSchemaAttribute,
		systemLogSizeAttribute:      systemLogSizeSchemaAttribute,
		systemNotesAttribute:        systemNotesSchemaAttribute,
		systemTimezoneAttribute:     systemTimezoneSchemaAttribute,
		systemTTYLoginAttribute:     systemTtyLoginSchemaAttribute,
		systemZonenameAttribute:     systemZonenameSchemaAttribute,
	}

	systemTimezoneSchemaAttribute = stringSchemaAttribute{
		Description:       "The POSIX.1 time zone string. This has no corresponding value in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ResourceExistence: NoValidation,
	}

	systemTtyLoginSchemaAttribute = boolSchemaAttribute{
		Description:       "Require authentication for local users to log in the system.",
		ResourceExistence: NoValidation,
	}

	systemZonenameSchemaAttribute = stringSchemaAttribute{
		Description:       "The IANA/Olson time zone string. This corresponds to \"Timezone\" in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ResourceExistence: NoValidation,
	}
)

type AttributeExistence int

func (e AttributeExistence) ToComputed() bool {
	return e == NoValidation || e == ReadOnly
}

func (e AttributeExistence) ToOptional() bool {
	return e == NoValidation || e == Optional
}

func (e AttributeExistence) ToRequired() bool {
	return e == Required
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

func ReadModel(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	client lucirpc.Client,
	sectionName string,
) (context.Context, systemModel, diag.Diagnostics) {
	tflog.Info(ctx, "Reading system model")
	var (
		allDiagnostics diag.Diagnostics
		model          systemModel
	)

	section, diagnostics := lucirpcglue.GetSection(ctx, client, systemUCIConfig, sectionName)
	allDiagnostics.Append(diagnostics...)
	if allDiagnostics.HasError() {
		return ctx, model, allDiagnostics
	}

	ctx, model.ConLogLevel, diagnostics = lucirpcglue.GetOptionInt64(ctx, fullTypeName, terraformType, section, path.Root(systemConLogLevelAttribute), systemConLogLevelUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.CronLogLevel, diagnostics = lucirpcglue.GetOptionInt64(ctx, fullTypeName, terraformType, section, path.Root(systemCronLogLevelAttribute), systemCronLogLevelUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Description, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemDescriptionAttribute), systemDescriptionUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Hostname, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemHostnameAttribute), systemHostnameUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.LogSize, diagnostics = lucirpcglue.GetOptionInt64(ctx, fullTypeName, terraformType, section, path.Root(systemLogSizeAttribute), systemLogSizeUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Notes, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemNotesAttribute), systemNotesUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Timezone, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemTimezoneAttribute), systemTimezoneUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.TTYLogin, diagnostics = lucirpcglue.GetOptionBool(ctx, fullTypeName, terraformType, section, path.Root(systemTTYLoginAttribute), systemTTYLoginUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Zonename, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemZonenameAttribute), systemZonenameUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Id, diagnostics = lucirpcglue.GetMetadataString(ctx, fullTypeName, terraformType, section, systemIdUCISection)
	allDiagnostics.Append(diagnostics...)

	return ctx, model, diagnostics
}

type boolSchemaAttribute struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ResourceExistence   AttributeExistence
	Sensitive           bool
	Validators          []validator.Bool
}

func (a boolSchemaAttribute) ToDataSource() datasourceschema.Attribute {
	return datasourceschema.BoolAttribute{
		Computed:            a.DataSourceExistence.ToComputed(),
		DeprecationMessage:  a.DeprecationMessage,
		Description:         a.Description,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.DataSourceExistence.ToOptional(),
		Required:            a.DataSourceExistence.ToRequired(),
		Sensitive:           a.Sensitive,
		Validators:          a.Validators,
	}
}

func (a boolSchemaAttribute) ToResource() resourceschema.Attribute {
	return resourceschema.BoolAttribute{
		Computed:            a.ResourceExistence.ToComputed(),
		DeprecationMessage:  a.DeprecationMessage,
		Description:         a.Description,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.ResourceExistence.ToOptional(),
		Required:            a.ResourceExistence.ToRequired(),
		Sensitive:           a.Sensitive,
		Validators:          a.Validators,
	}
}

type int64SchemaAttribute struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ResourceExistence   AttributeExistence
	Sensitive           bool
	Validators          []validator.Int64
}

func (a int64SchemaAttribute) ToDataSource() datasourceschema.Attribute {
	return datasourceschema.Int64Attribute{
		Computed:            a.DataSourceExistence.ToComputed(),
		DeprecationMessage:  a.DeprecationMessage,
		Description:         a.Description,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.DataSourceExistence.ToOptional(),
		Required:            a.DataSourceExistence.ToRequired(),
		Sensitive:           a.Sensitive,
		Validators:          a.Validators,
	}
}

func (a int64SchemaAttribute) ToResource() resourceschema.Attribute {
	return resourceschema.Int64Attribute{
		Computed:            a.ResourceExistence.ToComputed(),
		DeprecationMessage:  a.DeprecationMessage,
		Description:         a.Description,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.ResourceExistence.ToOptional(),
		Required:            a.ResourceExistence.ToRequired(),
		Sensitive:           a.Sensitive,
		Validators:          a.Validators,
	}
}

type schemaAttribute interface {
	ToDataSource() datasourceschema.Attribute
	ToResource() resourceschema.Attribute
}

type stringSchemaAttribute struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ResourceExistence   AttributeExistence
	Sensitive           bool
	Validators          []validator.String
}

func (a stringSchemaAttribute) ToDataSource() datasourceschema.Attribute {
	return datasourceschema.StringAttribute{
		Computed:            a.DataSourceExistence.ToComputed(),
		DeprecationMessage:  a.DeprecationMessage,
		Description:         a.Description,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.DataSourceExistence.ToOptional(),
		Required:            a.DataSourceExistence.ToRequired(),
		Sensitive:           a.Sensitive,
		Validators:          a.Validators,
	}
}

func (a stringSchemaAttribute) ToResource() resourceschema.Attribute {
	return resourceschema.StringAttribute{
		Computed:            a.ResourceExistence.ToComputed(),
		DeprecationMessage:  a.DeprecationMessage,
		Description:         a.Description,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.ResourceExistence.ToOptional(),
		Required:            a.ResourceExistence.ToRequired(),
		Sensitive:           a.Sensitive,
		Validators:          a.Validators,
	}
}
