package system

import (
	"context"
	"encoding/json"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
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
	systemConLogLevelSchemaAttribute = int64SchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description: "The maximum log level for kernel messages to be logged to the console.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetOptionInt64(ctx, fullTypeName, terraformType, section, path.Root(systemConLogLevelAttribute), systemConLogLevelUCIOption)
			model.ConLogLevel = value
			return ctx, model, diagnostics
		},
		ResourceExistence: NoValidation,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			if !hasValue(model.ConLogLevel) {
				return ctx, options, diag.Diagnostics{}
			}

			value, diagnostics := serializeInt64(model.ConLogLevel, path.Root(systemConLogLevelAttribute))
			if diagnostics.HasError() {
				return ctx, options, diagnostics
			}

			ctx = logger.SetFieldInt64(ctx, fullTypeName, resourceTerraformType, systemConLogLevelAttribute, model.ConLogLevel)
			options[systemConLogLevelUCIOption] = value
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemCronLogLevelSchemaAttribute = int64SchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description: "The minimum level for cron messages to be logged to syslog.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetOptionInt64(ctx, fullTypeName, terraformType, section, path.Root(systemCronLogLevelAttribute), systemCronLogLevelUCIOption)
			model.CronLogLevel = value
			return ctx, model, diagnostics
		},
		ResourceExistence: NoValidation,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			if !hasValue(model.CronLogLevel) {
				return ctx, options, diag.Diagnostics{}
			}

			value, diagnostics := serializeInt64(model.CronLogLevel, path.Root(systemCronLogLevelAttribute))
			if diagnostics.HasError() {
				return ctx, options, diagnostics
			}

			ctx = logger.SetFieldInt64(ctx, fullTypeName, resourceTerraformType, systemCronLogLevelAttribute, model.CronLogLevel)
			options[systemCronLogLevelUCIOption] = value
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemDescriptionSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description: "The hostname for the system.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemDescriptionAttribute), systemDescriptionUCIOption)
			model.Description = value
			return ctx, model, diagnostics
		},
		ResourceExistence: NoValidation,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			if !hasValue(model.Description) {
				return ctx, options, diag.Diagnostics{}
			}

			value, diagnostics := serializeString(model.Description, path.Root(systemDescriptionAttribute))
			if diagnostics.HasError() {
				return ctx, options, diagnostics
			}

			ctx = logger.SetFieldString(ctx, fullTypeName, resourceTerraformType, systemDescriptionAttribute, model.Description)
			options[systemDescriptionUCIOption] = value
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemHostnameSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description: "A short single-line description for the system.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemHostnameAttribute), systemHostnameUCIOption)
			model.Hostname = value
			return ctx, model, diagnostics
		},
		ResourceExistence: NoValidation,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			if !hasValue(model.Hostname) {
				return ctx, options, diag.Diagnostics{}
			}

			value, diagnostics := serializeString(model.Hostname, path.Root(systemHostnameAttribute))
			if diagnostics.HasError() {
				return ctx, options, diagnostics
			}

			ctx = logger.SetFieldString(ctx, fullTypeName, resourceTerraformType, systemHostnameAttribute, model.Hostname)
			options[systemHostnameUCIOption] = value
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemIdSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		DataSourceExistence: Required,
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
		ResourceExistence: Required,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			ctx = logger.SetFieldString(ctx, fullTypeName, resourceTerraformType, systemIdAttribute, model.Id)
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemLogSizeSchemaAttribute = int64SchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description: "Size of the file based log buffer in KiB.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetOptionInt64(ctx, fullTypeName, terraformType, section, path.Root(systemLogSizeAttribute), systemLogSizeUCIOption)
			model.LogSize = value
			return ctx, model, diagnostics
		},
		ResourceExistence: NoValidation,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			if !hasValue(model.LogSize) {
				return ctx, options, diag.Diagnostics{}
			}

			value, diagnostics := serializeInt64(model.LogSize, path.Root(systemLogSizeAttribute))
			if diagnostics.HasError() {
				return ctx, options, diagnostics
			}

			ctx = logger.SetFieldInt64(ctx, fullTypeName, resourceTerraformType, systemLogSizeAttribute, model.LogSize)
			options[systemLogSizeUCIOption] = value
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemNotesSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description: "Multi-line free-form text about the system.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemNotesAttribute), systemNotesUCIOption)
			model.Notes = value
			return ctx, model, diagnostics
		},
		ResourceExistence: NoValidation,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			if !hasValue(model.Notes) {
				return ctx, options, diag.Diagnostics{}
			}

			value, diagnostics := serializeString(model.Notes, path.Root(systemNotesAttribute))
			if diagnostics.HasError() {
				return ctx, options, diagnostics
			}

			ctx = logger.SetFieldString(ctx, fullTypeName, resourceTerraformType, systemNotesAttribute, model.Notes)
			options[systemNotesUCIOption] = value
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemSchemaAttributes = map[string]schemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
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

	systemTimezoneSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description: "The POSIX.1 time zone string. This has no corresponding value in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemTimezoneAttribute), systemTimezoneUCIOption)
			model.Timezone = value
			return ctx, model, diagnostics
		},
		ResourceExistence: NoValidation,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			if !hasValue(model.Timezone) {
				return ctx, options, diag.Diagnostics{}
			}

			value, diagnostics := serializeString(model.Timezone, path.Root(systemTimezoneAttribute))
			if diagnostics.HasError() {
				return ctx, options, diagnostics
			}

			ctx = logger.SetFieldString(ctx, fullTypeName, resourceTerraformType, systemTimezoneAttribute, model.Timezone)
			options[systemTimezoneUCIOption] = value
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemTtyLoginSchemaAttribute = boolSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description: "Require authentication for local users to log in the system.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetOptionBool(ctx, fullTypeName, terraformType, section, path.Root(systemTTYLoginAttribute), systemTTYLoginUCIOption)
			model.TTYLogin = value
			return ctx, model, diagnostics
		},
		ResourceExistence: NoValidation,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			if !hasValue(model.TTYLogin) {
				return ctx, options, diag.Diagnostics{}
			}

			value, diagnostics := serializeBool(model.TTYLogin, path.Root(systemTTYLoginAttribute))
			if diagnostics.HasError() {
				return ctx, options, diagnostics
			}

			ctx = logger.SetFieldBool(ctx, fullTypeName, resourceTerraformType, systemTTYLoginAttribute, model.TTYLogin)
			options[systemTTYLoginUCIOption] = value
			return ctx, options, diag.Diagnostics{}
		},
	}

	systemZonenameSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description: "The IANA/Olson time zone string. This corresponds to \"Timezone\" in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model systemModel,
		) (context.Context, systemModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(systemZonenameAttribute), systemZonenameUCIOption)
			model.Zonename = value
			return ctx, model, diagnostics
		},
		ResourceExistence: NoValidation,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model systemModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			if !hasValue(model.Zonename) {
				return ctx, options, diag.Diagnostics{}
			}

			value, diagnostics := serializeString(model.Zonename, path.Root(systemZonenameAttribute))
			if diagnostics.HasError() {
				return ctx, options, diagnostics
			}

			ctx = logger.SetFieldString(ctx, fullTypeName, resourceTerraformType, systemZonenameAttribute, model.Zonename)
			options[systemZonenameUCIOption] = value
			return ctx, options, diag.Diagnostics{}
		},
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

	for _, attribute := range systemSchemaAttributes {
		ctx, model, diagnostics = attribute.Read(ctx, fullTypeName, terraformType, section, model)
		allDiagnostics.Append(diagnostics...)
	}

	return ctx, model, diagnostics
}

type boolSchemaAttribute[Model any, Request any, Response any] struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ReadResponse        func(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ResourceExistence   AttributeExistence
	Sensitive           bool
	UpsertRequest       func(context.Context, string, string, Request, Model) (context.Context, Request, diag.Diagnostics)
	Validators          []validator.Bool
}

func (a boolSchemaAttribute[Model, Request, Response]) Read(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	response Response,
	model Model,
) (context.Context, Model, diag.Diagnostics) {
	if a.ReadResponse == nil {
		return ctx, model, diag.Diagnostics{}
	}

	return a.ReadResponse(ctx, fullTypeName, terraformType, response, model)
}

func (a boolSchemaAttribute[Model, Request, Response]) ToDataSource() datasourceschema.Attribute {
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

func (a boolSchemaAttribute[Model, Request, Response]) ToResource() resourceschema.Attribute {
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

func (a boolSchemaAttribute[Model, Request, Response]) Upsert(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	request Request,
	model Model,
) (context.Context, Request, diag.Diagnostics) {
	if a.UpsertRequest == nil {
		return ctx, request, diag.Diagnostics{}
	}

	return a.UpsertRequest(ctx, fullTypeName, terraformType, request, model)
}

type int64SchemaAttribute[Model any, Request any, Response any] struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ReadResponse        func(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ResourceExistence   AttributeExistence
	Sensitive           bool
	UpsertRequest       func(context.Context, string, string, Request, Model) (context.Context, Request, diag.Diagnostics)
	Validators          []validator.Int64
}

func (a int64SchemaAttribute[Model, Request, Response]) Read(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	response Response,
	model Model,
) (context.Context, Model, diag.Diagnostics) {
	if a.ReadResponse == nil {
		return ctx, model, diag.Diagnostics{}
	}

	return a.ReadResponse(ctx, fullTypeName, terraformType, response, model)
}

func (a int64SchemaAttribute[Model, Request, Response]) ToDataSource() datasourceschema.Attribute {
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

func (a int64SchemaAttribute[Model, Request, Response]) ToResource() resourceschema.Attribute {
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

func (a int64SchemaAttribute[Model, Request, Response]) Upsert(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	request Request,
	model Model,
) (context.Context, Request, diag.Diagnostics) {
	if a.UpsertRequest == nil {
		return ctx, request, diag.Diagnostics{}
	}

	return a.UpsertRequest(ctx, fullTypeName, terraformType, request, model)
}

type schemaAttribute[Model any, Request any, Response any] interface {
	Read(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ToDataSource() datasourceschema.Attribute
	ToResource() resourceschema.Attribute
	Upsert(context.Context, string, string, Request, Model) (context.Context, Request, diag.Diagnostics)
}

type stringSchemaAttribute[Model any, Request any, Response any] struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ReadResponse        func(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ResourceExistence   AttributeExistence
	Sensitive           bool
	UpsertRequest       func(context.Context, string, string, Request, Model) (context.Context, Request, diag.Diagnostics)
	Validators          []validator.String
}

func (a stringSchemaAttribute[Model, Request, Response]) Read(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	response Response,
	model Model,
) (context.Context, Model, diag.Diagnostics) {
	if a.ReadResponse == nil {
		return ctx, model, diag.Diagnostics{}
	}

	return a.ReadResponse(ctx, fullTypeName, terraformType, response, model)
}

func (a stringSchemaAttribute[Model, Request, Response]) ToDataSource() datasourceschema.Attribute {
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

func (a stringSchemaAttribute[Model, Request, Response]) ToResource() resourceschema.Attribute {
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

func (a stringSchemaAttribute[Model, Request, Response]) Upsert(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	request Request,
	model Model,
) (context.Context, Request, diag.Diagnostics) {
	if a.UpsertRequest == nil {
		return ctx, request, diag.Diagnostics{}
	}

	return a.UpsertRequest(ctx, fullTypeName, terraformType, request, model)
}
