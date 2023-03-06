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
		Description:       "The maximum log level for kernel messages to be logged to the console.",
		ReadResponse:      ReadResponseOptionInt64(systemModelSetConLogLevel, systemConLogLevelAttribute, systemConLogLevelUCIOption),
		ResourceExistence: NoValidation,
		UpsertRequest:     UpsertRequestOptionInt64(systemModelGetConLogLevel, systemConLogLevelAttribute, systemConLogLevelUCIOption),
	}

	systemCronLogLevelSchemaAttribute = int64SchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "The minimum level for cron messages to be logged to syslog.",
		ReadResponse:      ReadResponseOptionInt64(systemModelSetCronLogLevel, systemCronLogLevelAttribute, systemCronLogLevelUCIOption),
		ResourceExistence: NoValidation,
		UpsertRequest:     UpsertRequestOptionInt64(systemModelGetCronLogLevel, systemCronLogLevelAttribute, systemCronLogLevelUCIOption),
	}

	systemDescriptionSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "The hostname for the system.",
		ReadResponse:      ReadResponseOptionString(systemModelSetDescription, systemDescriptionAttribute, systemDescriptionUCIOption),
		ResourceExistence: NoValidation,
		UpsertRequest:     UpsertRequestOptionString(systemModelGetDescription, systemDescriptionAttribute, systemDescriptionUCIOption),
	}

	systemHostnameSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "A short single-line description for the system.",
		ReadResponse:      ReadResponseOptionString(systemModelSetHostname, systemHostnameAttribute, systemHostnameUCIOption),
		ResourceExistence: NoValidation,
		UpsertRequest:     UpsertRequestOptionString(systemModelGetHostname, systemHostnameAttribute, systemHostnameUCIOption),
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
		Description:       "Size of the file based log buffer in KiB.",
		ReadResponse:      ReadResponseOptionInt64(systemModelSetLogSize, systemLogSizeAttribute, systemLogSizeUCIOption),
		ResourceExistence: NoValidation,
		UpsertRequest:     UpsertRequestOptionInt64(systemModelGetLogSize, systemLogSizeAttribute, systemLogSizeUCIOption),
	}

	systemNotesSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "Multi-line free-form text about the system.",
		ReadResponse:      ReadResponseOptionString(systemModelSetNotes, systemNotesAttribute, systemNotesUCIOption),
		ResourceExistence: NoValidation,
		UpsertRequest:     UpsertRequestOptionString(systemModelGetNotes, systemNotesAttribute, systemNotesUCIOption),
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
		Description:       "The POSIX.1 time zone string. This has no corresponding value in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ReadResponse:      ReadResponseOptionString(systemModelSetTimezone, systemTimezoneAttribute, systemTimezoneUCIOption),
		ResourceExistence: NoValidation,
		UpsertRequest:     UpsertRequestOptionString(systemModelGetTimezone, systemTimezoneAttribute, systemTimezoneUCIOption),
	}

	systemTtyLoginSchemaAttribute = boolSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "Require authentication for local users to log in the system.",
		ReadResponse:      ReadResponseOptionBool(systemModelSetTTYLogin, systemTTYLoginAttribute, systemTTYLoginUCIOption),
		ResourceExistence: NoValidation,
		UpsertRequest:     UpsertRequestOptionBool(systemModelGetTTYLogin, systemTTYLoginAttribute, systemTTYLoginUCIOption),
	}

	systemZonenameSchemaAttribute = stringSchemaAttribute[systemModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "The IANA/Olson time zone string. This corresponds to \"Timezone\" in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ReadResponse:      ReadResponseOptionString(systemModelSetZonename, systemZonenameAttribute, systemZonenameUCIOption),
		ResourceExistence: NoValidation,
		UpsertRequest:     UpsertRequestOptionString(systemModelGetZonename, systemZonenameAttribute, systemZonenameUCIOption),
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

func ReadResponseOptionBool[Model any](
	set func(Model, types.Bool) Model,
	attribute string,
	option string,
) func(context.Context, string, string, map[string]json.RawMessage, Model) (context.Context, Model, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		section map[string]json.RawMessage,
		model Model,
	) (context.Context, Model, diag.Diagnostics) {
		ctx, value, diagnostics := lucirpcglue.GetOptionBool(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
		model = set(model, value)
		return ctx, model, diagnostics
	}
}

func ReadResponseOptionInt64[Model any](
	set func(Model, types.Int64) Model,
	attribute string,
	option string,
) func(context.Context, string, string, map[string]json.RawMessage, Model) (context.Context, Model, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		section map[string]json.RawMessage,
		model Model,
	) (context.Context, Model, diag.Diagnostics) {
		ctx, value, diagnostics := lucirpcglue.GetOptionInt64(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
		model = set(model, value)
		return ctx, model, diagnostics
	}
}

func ReadResponseOptionString[Model any](
	set func(Model, types.String) Model,
	attribute string,
	option string,
) func(context.Context, string, string, map[string]json.RawMessage, Model) (context.Context, Model, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		section map[string]json.RawMessage,
		model Model,
	) (context.Context, Model, diag.Diagnostics) {
		ctx, value, diagnostics := lucirpcglue.GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
		model = set(model, value)
		return ctx, model, diagnostics
	}
}

func UpsertRequestOptionBool[Model any](
	get func(Model) types.Bool,
	attribute string,
	option string,
) func(context.Context, string, string, map[string]json.RawMessage, Model) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		options map[string]json.RawMessage,
		model Model,
	) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
		str := get(model)
		if !hasValue(str) {
			return ctx, options, diag.Diagnostics{}
		}

		value, diagnostics := serializeBool(str, path.Root(attribute))
		if diagnostics.HasError() {
			return ctx, options, diagnostics
		}

		ctx = logger.SetFieldBool(ctx, fullTypeName, terraformType, attribute, str)
		options[option] = value
		return ctx, options, diag.Diagnostics{}
	}
}

func UpsertRequestOptionInt64[Model any](
	get func(Model) types.Int64,
	attribute string,
	option string,
) func(context.Context, string, string, map[string]json.RawMessage, Model) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		options map[string]json.RawMessage,
		model Model,
	) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
		str := get(model)
		if !hasValue(str) {
			return ctx, options, diag.Diagnostics{}
		}

		value, diagnostics := serializeInt64(str, path.Root(attribute))
		if diagnostics.HasError() {
			return ctx, options, diagnostics
		}

		ctx = logger.SetFieldInt64(ctx, fullTypeName, terraformType, attribute, str)
		options[option] = value
		return ctx, options, diag.Diagnostics{}
	}
}

func UpsertRequestOptionString[Model any](
	get func(Model) types.String,
	attribute string,
	option string,
) func(context.Context, string, string, map[string]json.RawMessage, Model) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		options map[string]json.RawMessage,
		model Model,
	) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
		str := get(model)
		if !hasValue(str) {
			return ctx, options, diag.Diagnostics{}
		}

		value, diagnostics := serializeString(str, path.Root(attribute))
		if diagnostics.HasError() {
			return ctx, options, diagnostics
		}

		ctx = logger.SetFieldString(ctx, fullTypeName, terraformType, attribute, str)
		options[option] = value
		return ctx, options, diag.Diagnostics{}
	}
}

type attributeHasValue interface {
	IsNull() bool
	IsUnknown() bool
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

func hasValue(
	attribute attributeHasValue,
) bool {
	return !attribute.IsNull() && !attribute.IsUnknown()
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

func serializeBool(
	attribute interface{ ValueBool() bool },
	attributePath path.Path,
) (json.RawMessage, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	value, err := json.Marshal(attribute.ValueBool())
	if err != nil {
		diagnostics.AddAttributeError(
			attributePath,
			"Could not serialize",
			err.Error(),
		)
		return nil, diagnostics
	}

	return json.RawMessage(value), diagnostics
}

func serializeInt64(
	attribute interface{ ValueInt64() int64 },
	attributePath path.Path,
) (json.RawMessage, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	value, err := json.Marshal(attribute.ValueInt64())
	if err != nil {
		diagnostics.AddAttributeError(
			attributePath,
			"Could not serialize",
			err.Error(),
		)
		return nil, diagnostics
	}

	return json.RawMessage(value), diagnostics
}

func serializeString(
	attribute interface{ ValueString() string },
	attributePath path.Path,
) (json.RawMessage, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	value, err := json.Marshal(attribute.ValueString())
	if err != nil {
		diagnostics.AddAttributeError(
			attributePath,
			"Could not serialize",
			err.Error(),
		)
		return nil, diagnostics
	}

	return json.RawMessage(value), diagnostics
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

func systemModelGetConLogLevel(model systemModel) types.Int64  { return model.ConLogLevel }
func systemModelGetCronLogLevel(model systemModel) types.Int64 { return model.CronLogLevel }
func systemModelGetDescription(model systemModel) types.String { return model.Description }
func systemModelGetHostname(model systemModel) types.String    { return model.Hostname }
func systemModelGetLogSize(model systemModel) types.Int64      { return model.LogSize }
func systemModelGetNotes(model systemModel) types.String       { return model.Notes }
func systemModelGetTimezone(model systemModel) types.String    { return model.Timezone }
func systemModelGetTTYLogin(model systemModel) types.Bool      { return model.TTYLogin }
func systemModelGetZonename(model systemModel) types.String    { return model.Zonename }

func systemModelSetConLogLevel(model systemModel, value types.Int64) systemModel {
	model.ConLogLevel = value
	return model
}

func systemModelSetCronLogLevel(model systemModel, value types.Int64) systemModel {
	model.CronLogLevel = value
	return model
}

func systemModelSetDescription(model systemModel, value types.String) systemModel {
	model.Description = value
	return model
}

func systemModelSetHostname(model systemModel, value types.String) systemModel {
	model.Hostname = value
	return model
}

func systemModelSetLogSize(model systemModel, value types.Int64) systemModel {
	model.LogSize = value
	return model
}

func systemModelSetNotes(model systemModel, value types.String) systemModel {
	model.Notes = value
	return model
}

func systemModelSetTimezone(model systemModel, value types.String) systemModel {
	model.Timezone = value
	return model
}

func systemModelSetTTYLogin(model systemModel, value types.Bool) systemModel {
	model.TTYLogin = value
	return model
}

func systemModelSetZonename(model systemModel, value types.String) systemModel {
	model.Zonename = value
	return model
}
