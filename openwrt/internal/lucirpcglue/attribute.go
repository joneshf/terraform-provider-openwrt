package lucirpcglue

import (
	"context"
	"encoding/json"

	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
)

const (
	ReadOnly AttributeExistence = iota
	NoValidation
	Optional
	Required
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

type BoolSchemaAttribute[Model any, Request any, Response any] struct {
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

func (a BoolSchemaAttribute[Model, Request, Response]) Read(
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

func (a BoolSchemaAttribute[Model, Request, Response]) ToDataSource() datasourceschema.Attribute {
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

func (a BoolSchemaAttribute[Model, Request, Response]) ToResource() resourceschema.Attribute {
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

func (a BoolSchemaAttribute[Model, Request, Response]) Upsert(
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

type Int64SchemaAttribute[Model any, Request any, Response any] struct {
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

func (a Int64SchemaAttribute[Model, Request, Response]) Read(
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

func (a Int64SchemaAttribute[Model, Request, Response]) ToDataSource() datasourceschema.Attribute {
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

func (a Int64SchemaAttribute[Model, Request, Response]) ToResource() resourceschema.Attribute {
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

func (a Int64SchemaAttribute[Model, Request, Response]) Upsert(
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

func ReadResponseOptionBool[Model any](
	set func(*Model, types.Bool),
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
		ctx, value, diagnostics := GetOptionBool(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
		set(&model, value)
		return ctx, model, diagnostics
	}
}

func ReadResponseOptionInt64[Model any](
	set func(*Model, types.Int64),
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
		ctx, value, diagnostics := GetOptionInt64(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
		set(&model, value)
		return ctx, model, diagnostics
	}
}

func ReadResponseOptionString[Model any](
	set func(*Model, types.String),
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
		ctx, value, diagnostics := GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
		set(&model, value)
		return ctx, model, diagnostics
	}
}

type SchemaAttribute[Model any, Request any, Response any] interface {
	Read(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ToDataSource() datasourceschema.Attribute
	ToResource() resourceschema.Attribute
	Upsert(context.Context, string, string, Request, Model) (context.Context, Request, diag.Diagnostics)
}

type StringSchemaAttribute[Model any, Request any, Response any] struct {
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

func (a StringSchemaAttribute[Model, Request, Response]) Read(
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

func (a StringSchemaAttribute[Model, Request, Response]) ToDataSource() datasourceschema.Attribute {
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

func (a StringSchemaAttribute[Model, Request, Response]) ToResource() resourceschema.Attribute {
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

func (a StringSchemaAttribute[Model, Request, Response]) Upsert(
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

func hasValue(
	attribute attributeHasValue,
) bool {
	return !attribute.IsNull() && !attribute.IsUnknown()
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
