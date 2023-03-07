package system

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
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
	systemConLogLevelSchemaAttribute = lucirpcglue.Int64SchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "The maximum log level for kernel messages to be logged to the console.",
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(systemModelSetConLogLevel, systemConLogLevelAttribute, systemConLogLevelUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(systemModelGetConLogLevel, systemConLogLevelAttribute, systemConLogLevelUCIOption),
	}

	systemCronLogLevelSchemaAttribute = lucirpcglue.Int64SchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "The minimum level for cron messages to be logged to syslog.",
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(systemModelSetCronLogLevel, systemCronLogLevelAttribute, systemCronLogLevelUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(systemModelGetCronLogLevel, systemCronLogLevelAttribute, systemCronLogLevelUCIOption),
	}

	systemDescriptionSchemaAttribute = lucirpcglue.StringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "The hostname for the system.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(systemModelSetDescription, systemDescriptionAttribute, systemDescriptionUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(systemModelGetDescription, systemDescriptionAttribute, systemDescriptionUCIOption),
	}

	systemHostnameSchemaAttribute = lucirpcglue.StringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "A short single-line description for the system.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(systemModelSetHostname, systemHostnameAttribute, systemHostnameUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(systemModelGetHostname, systemHostnameAttribute, systemHostnameUCIOption),
	}

	systemIdSchemaAttribute = lucirpcglue.StringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		DataSourceExistence: lucirpcglue.Required,
		Description:         "Placeholder identifier attribute.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetMetadataString(ctx, fullTypeName, terraformType, section, systemIdUCISection)
			model.Id = value
			return ctx, model, diagnostics
		},
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			ctx = logger.SetFieldString(ctx, fullTypeName, lucirpcglue.ResourceTerraformType, systemIdAttribute, model.Id)
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemLogSizeSchemaAttribute = lucirpcglue.Int64SchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "Size of the file based log buffer in KiB.",
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(systemModelSetLogSize, systemLogSizeAttribute, systemLogSizeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(systemModelGetLogSize, systemLogSizeAttribute, systemLogSizeUCIOption),
	}

	systemNotesSchemaAttribute = lucirpcglue.StringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "Multi-line free-form text about the system.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(systemModelSetNotes, systemNotesAttribute, systemNotesUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(systemModelGetNotes, systemNotesAttribute, systemNotesUCIOption),
	}

	systemSchemaAttributes = map[string]lucirpcglue.SchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
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

	systemTimezoneSchemaAttribute = lucirpcglue.StringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "The POSIX.1 time zone string. This has no corresponding value in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(systemModelSetTimezone, systemTimezoneAttribute, systemTimezoneUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(systemModelGetTimezone, systemTimezoneAttribute, systemTimezoneUCIOption),
	}

	systemTtyLoginSchemaAttribute = lucirpcglue.BoolSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "Require authentication for local users to log in the system.",
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(systemModelSetTTYLogin, systemTTYLoginAttribute, systemTTYLoginUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(systemModelGetTTYLogin, systemTTYLoginAttribute, systemTTYLoginUCIOption),
	}

	systemZonenameSchemaAttribute = lucirpcglue.StringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "The IANA/Olson time zone string. This corresponds to \"Timezone\" in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(systemModelSetZonename, systemZonenameAttribute, systemZonenameUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(systemModelGetZonename, systemZonenameAttribute, systemZonenameUCIOption),
	}
)

func readSystemModel(
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

	for _, attribute := range systemSchemaAttributes {
		ctx, model, diagnostics = attribute.Read(ctx, fullTypeName, terraformType, section, model)
		allDiagnostics.Append(diagnostics...)
	}

	return ctx, model, diagnostics
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

func systemModelGetConLogLevel(model systemModel) types.Int64  { return model.ConLogLevel }
func systemModelGetCronLogLevel(model systemModel) types.Int64 { return model.CronLogLevel }
func systemModelGetDescription(model systemModel) types.String { return model.Description }
func systemModelGetHostname(model systemModel) types.String    { return model.Hostname }
func systemModelGetLogSize(model systemModel) types.Int64      { return model.LogSize }
func systemModelGetNotes(model systemModel) types.String       { return model.Notes }
func systemModelGetTimezone(model systemModel) types.String    { return model.Timezone }
func systemModelGetTTYLogin(model systemModel) types.Bool      { return model.TTYLogin }
func systemModelGetZonename(model systemModel) types.String    { return model.Zonename }

func systemModelSetConLogLevel(model *systemModel, value types.Int64)  { model.ConLogLevel = value }
func systemModelSetCronLogLevel(model *systemModel, value types.Int64) { model.CronLogLevel = value }
func systemModelSetDescription(model *systemModel, value types.String) { model.Description = value }
func systemModelSetHostname(model *systemModel, value types.String)    { model.Hostname = value }
func systemModelSetLogSize(model *systemModel, value types.Int64)      { model.LogSize = value }
func systemModelSetNotes(model *systemModel, value types.String)       { model.Notes = value }
func systemModelSetTimezone(model *systemModel, value types.String)    { model.Timezone = value }
func systemModelSetTTYLogin(model *systemModel, value types.Bool)      { model.TTYLogin = value }
func systemModelSetZonename(model *systemModel, value types.String)    { model.Zonename = value }
