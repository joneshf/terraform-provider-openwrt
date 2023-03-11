package lucirpcglue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
)

const (
	ReadOnly AttributeExistence = iota
	NoValidation
	Optional
	Required
)

var (
	_ validator.Bool   = requiresAttribute[any]{}
	_ validator.Int64  = requiresAttribute[any]{}
	_ validator.Set    = requiresAttribute[any]{}
	_ validator.String = requiresAttribute[any]{}
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
	UpsertRequest       func(context.Context, string, Request, Model) (context.Context, Request, diag.Diagnostics)
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
	request Request,
	model Model,
) (context.Context, Request, diag.Diagnostics) {
	if a.UpsertRequest == nil {
		return ctx, request, diag.Diagnostics{}
	}

	return a.UpsertRequest(ctx, fullTypeName, request, model)
}

type Int64SchemaAttribute[Model any, Request any, Response any] struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ReadResponse        func(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ResourceExistence   AttributeExistence
	Sensitive           bool
	UpsertRequest       func(context.Context, string, Request, Model) (context.Context, Request, diag.Diagnostics)
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
	request Request,
	model Model,
) (context.Context, Request, diag.Diagnostics) {
	if a.UpsertRequest == nil {
		return ctx, request, diag.Diagnostics{}
	}

	return a.UpsertRequest(ctx, fullTypeName, request, model)
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

func ReadResponseOptionSetString[Model any](
	set func(*Model, types.Set),
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
		ctx, value, diagnostics := GetOptionSetString(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
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
func RequiresAttributeEqualBool(
	expression path.Expression,
	expected bool,
) requiresAttribute[bool] {
	return requiresAttributeEqual(
		types.BoolType,
		expression,
		expected,
	)
}

func RequiresAttributeEqualString(
	expression path.Expression,
	expected string,
) requiresAttribute[string] {
	return requiresAttributeEqual(
		types.StringType,
		expression,
		expected,
	)
}

type SchemaAttribute[Model any, Request any, Response any] interface {
	Read(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ToDataSource() datasourceschema.Attribute
	ToResource() resourceschema.Attribute
	Upsert(context.Context, string, Request, Model) (context.Context, Request, diag.Diagnostics)
}

type SetStringSchemaAttribute[Model any, Request any, Response any] struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ReadResponse        func(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ResourceExistence   AttributeExistence
	Sensitive           bool
	UpsertRequest       func(context.Context, string, Request, Model) (context.Context, Request, diag.Diagnostics)
	Validators          []validator.Set
}

func (a SetStringSchemaAttribute[Model, Request, Response]) Read(
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

func (a SetStringSchemaAttribute[Model, Request, Response]) ToDataSource() datasourceschema.Attribute {
	return datasourceschema.SetAttribute{
		Computed:            a.DataSourceExistence.ToComputed(),
		DeprecationMessage:  a.DeprecationMessage,
		Description:         a.Description,
		ElementType:         types.StringType,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.DataSourceExistence.ToOptional(),
		Required:            a.DataSourceExistence.ToRequired(),
		Sensitive:           a.Sensitive,
		Validators:          a.Validators,
	}
}

func (a SetStringSchemaAttribute[Model, Request, Response]) ToResource() resourceschema.Attribute {
	return resourceschema.SetAttribute{
		Computed:            a.ResourceExistence.ToComputed(),
		DeprecationMessage:  a.DeprecationMessage,
		Description:         a.Description,
		ElementType:         types.StringType,
		MarkdownDescription: a.MarkdownDescription,
		Optional:            a.ResourceExistence.ToOptional(),
		Required:            a.ResourceExistence.ToRequired(),
		Sensitive:           a.Sensitive,
		Validators:          a.Validators,
	}
}

func (a SetStringSchemaAttribute[Model, Request, Response]) Upsert(
	ctx context.Context,
	fullTypeName string,
	request Request,
	model Model,
) (context.Context, Request, diag.Diagnostics) {
	if a.UpsertRequest == nil {
		return ctx, request, diag.Diagnostics{}
	}

	return a.UpsertRequest(ctx, fullTypeName, request, model)
}

type StringSchemaAttribute[Model any, Request any, Response any] struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ReadResponse        func(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ResourceExistence   AttributeExistence
	Sensitive           bool
	UpsertRequest       func(context.Context, string, Request, Model) (context.Context, Request, diag.Diagnostics)
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
	request Request,
	model Model,
) (context.Context, Request, diag.Diagnostics) {
	if a.UpsertRequest == nil {
		return ctx, request, diag.Diagnostics{}
	}

	return a.UpsertRequest(ctx, fullTypeName, request, model)
}

func UpsertRequestOptionBool[Model any](
	get func(Model) types.Bool,
	attribute string,
	option string,
) func(context.Context, string, map[string]json.RawMessage, Model) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
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

		ctx = logger.SetFieldBool(ctx, fullTypeName, ResourceTerraformType, attribute, str)
		options[option] = value
		return ctx, options, diag.Diagnostics{}
	}
}

func UpsertRequestOptionInt64[Model any](
	get func(Model) types.Int64,
	attribute string,
	option string,
) func(context.Context, string, map[string]json.RawMessage, Model) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
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

		ctx = logger.SetFieldInt64(ctx, fullTypeName, ResourceTerraformType, attribute, str)
		options[option] = value
		return ctx, options, diag.Diagnostics{}
	}
}

func UpsertRequestOptionSetString[Model any](
	get func(Model) types.Set,
	attribute string,
	option string,
) func(context.Context, string, map[string]json.RawMessage, Model) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		options map[string]json.RawMessage,
		model Model,
	) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
		str := get(model)
		if !hasValue(str) {
			return ctx, options, diag.Diagnostics{}
		}

		value, diagnostics := serializeSetString(ctx, str, path.Root(attribute))
		if diagnostics.HasError() {
			return ctx, options, diagnostics
		}

		ctx = logger.SetFieldSetString(ctx, fullTypeName, ResourceTerraformType, attribute, str)
		options[option] = value
		return ctx, options, diag.Diagnostics{}
	}
}

func UpsertRequestOptionString[Model any](
	get func(Model) types.String,
	attribute string,
	option string,
) func(context.Context, string, map[string]json.RawMessage, Model) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
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

		ctx = logger.SetFieldString(ctx, fullTypeName, ResourceTerraformType, attribute, str)
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

type requiresAttribute[Value any] struct {
	attrType   attr.Type
	expected   Value
	expression path.Expression
}

func (a requiresAttribute[Value]) Description(ctx context.Context) string {
	return a.MarkdownDescription(ctx)
}

func (a requiresAttribute[Value]) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures that if an attribute is set, %q is also set to %v", a.expression, a.expected)
}

func (a requiresAttribute[Value]) ValidateBool(
	ctx context.Context,
	req validator.BoolRequest,
	res *validator.BoolResponse,
) {
	diagnostics := a.validate(ctx, req.Config, req.Path, req.ConfigValue)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

func (a requiresAttribute[Value]) ValidateInt64(
	ctx context.Context,
	req validator.Int64Request,
	res *validator.Int64Response,
) {
	diagnostics := a.validate(ctx, req.Config, req.Path, req.ConfigValue)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

func (a requiresAttribute[Value]) ValidateSet(
	ctx context.Context,
	req validator.SetRequest,
	res *validator.SetResponse,
) {
	diagnostics := a.validate(ctx, req.Config, req.Path, req.ConfigValue)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

func (a requiresAttribute[Value]) ValidateString(
	ctx context.Context,
	req validator.StringRequest,
	res *validator.StringResponse,
) {
	diagnostics := a.validate(ctx, req.Config, req.Path, req.ConfigValue)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

func (a requiresAttribute[Value]) validate(
	ctx context.Context,
	config tfsdk.Config,
	requestPath path.Path,
	configValue interface{ IsNull() bool },
) (allDiagnostics diag.Diagnostics) {
	if configValue.IsNull() {
		return
	}

	matchedPaths, diagnostics := config.PathMatches(ctx, a.expression)
	allDiagnostics.Append(diagnostics...)
	if allDiagnostics.HasError() {
		return
	}

	for _, matchedPath := range matchedPaths {
		if matchedPath.Equal(requestPath) {
			allDiagnostics.Append(
				validatordiag.BugInProviderDiagnostic(
					fmt.Sprintf("Attribute %q cannot require itself to have a specific value", requestPath),
				),
			)
			continue
		}

		var actual attr.Value
		diagnostics = config.GetAttribute(ctx, matchedPath, &actual)
		allDiagnostics.Append(diagnostics...)
		if allDiagnostics.HasError() {
			continue
		}

		if actual.IsUnknown() {
			// Ignore this value until it is known.
			continue
		}

		if actual.IsNull() {
			allDiagnostics.Append(
				validatordiag.InvalidAttributeCombinationDiagnostic(
					requestPath,
					fmt.Sprintf("Attribute %q must be %v when %q is specified, but %q was not specified", matchedPath, a.expected, requestPath, matchedPath),
				),
			)
			continue
		}

		var expected attr.Value
		diagnostics = tfsdk.ValueFrom(ctx, a.expected, a.attrType, &expected)
		allDiagnostics.Append(diagnostics...)
		if allDiagnostics.HasError() {
			continue
		}

		if !actual.Equal(expected) {
			allDiagnostics.Append(
				validatordiag.InvalidAttributeCombinationDiagnostic(
					requestPath,
					fmt.Sprintf("Attribute %q must be %v when %q is specified, but %q was %v", matchedPath, expected, requestPath, matchedPath, actual),
				),
			)
			continue
		}
	}

	if allDiagnostics.HasError() {
		return
	}

	return
}

func requiresAttributeEqual[Value any](
	attrType attr.Type,
	expression path.Expression,
	expected Value,
) requiresAttribute[Value] {
	return requiresAttribute[Value]{
		attrType:   attrType,
		expected:   expected,
		expression: expression,
	}
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

func serializeSetString(
	ctx context.Context,
	attribute interface{ Elements() []attr.Value },
	attributePath path.Path,
) (json.RawMessage, diag.Diagnostics) {
	allDiagnostics := diag.Diagnostics{}
	elements := attribute.Elements()
	values := []string{}
	for _, element := range elements {
		var value string
		diagnostics := tfsdk.ValueAs(ctx, element, &value)
		allDiagnostics.Append(diagnostics...)
		if allDiagnostics.HasError() {
			// We don't want to exit early.
			// We want to continue to accumulate diagnostics.
			continue
		}

		values = append(values, value)
	}

	if allDiagnostics.HasError() {
		return nil, allDiagnostics
	}

	value, err := json.Marshal(values)
	if err != nil {
		allDiagnostics.AddAttributeError(
			attributePath,
			"Could not serialize",
			err.Error(),
		)
		return nil, allDiagnostics
	}

	return json.RawMessage(value), allDiagnostics
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
