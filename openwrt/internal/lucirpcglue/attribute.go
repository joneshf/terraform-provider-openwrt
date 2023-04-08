package lucirpcglue

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
)

const (
	ReadOnly AttributeExistence = iota
	NoValidation
	Optional
	Required

	idAttributeDescription = "Name of the section. This name is only used when interacting with UCI directly."
	idUCISection           = ".name"

	IdAttribute = "id"
)

var (
	_ validator.Bool = anyValidatorBool{}

	_ validator.Bool   = requiredIfAttributeNot[any]{}
	_ validator.Int64  = requiredIfAttributeNot[any]{}
	_ validator.List   = requiredIfAttributeNot[any]{}
	_ validator.Set    = requiredIfAttributeNot[any]{}
	_ validator.String = requiredIfAttributeNot[any]{}

	_ validator.Bool   = requiresAttribute[any]{}
	_ validator.Int64  = requiresAttribute[any]{}
	_ validator.List   = requiresAttribute[any]{}
	_ validator.Set    = requiresAttribute[any]{}
	_ validator.String = requiresAttribute[any]{}
)

// AnyBool returns a validator which ensures that any configured attribute value passes at least one of the given validators.
func AnyBool(validators ...validator.Bool) validator.Bool {
	return anyValidatorBool{
		validators: validators,
	}
}

type anyValidatorBool struct {
	validators []validator.Bool
}

func (v anyValidatorBool) Description(ctx context.Context) string {
	var descriptions []string

	for _, subValidator := range v.validators {
		descriptions = append(descriptions, subValidator.Description(ctx))
	}

	return fmt.Sprintf("Value must satisfy at least one of the validations: %s", strings.Join(descriptions, " + "))
}

func (v anyValidatorBool) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v anyValidatorBool) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	for _, subValidator := range v.validators {
		validateResp := &validator.BoolResponse{}

		subValidator.ValidateBool(ctx, req, validateResp)

		if !validateResp.Diagnostics.HasError() {
			resp.Diagnostics = validateResp.Diagnostics

			return
		}

		resp.Diagnostics.Append(validateResp.Diagnostics...)
	}
}

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

func IdSchemaAttribute[Model any](
	get func(Model) types.String,
	set func(*Model, types.String),
) SchemaAttribute[Model, lucirpc.Options, lucirpc.Options] {
	return StringSchemaAttribute[Model, lucirpc.Options, lucirpc.Options]{
		DataSourceExistence: Required,
		Description:         idAttributeDescription,
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section lucirpc.Options,
			model Model,
		) (context.Context, Model, diag.Diagnostics) {
			ctx, value, diagnostics := GetMetadataString(ctx, fullTypeName, terraformType, section, idUCISection)
			set(&model, value)
			return ctx, model, diagnostics
		},
		ResourceExistence: Required,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			options lucirpc.Options,
			model Model,
		) (context.Context, lucirpc.Options, diag.Diagnostics) {
			id := get(model)
			ctx = logger.SetFieldString(ctx, fullTypeName, ResourceTerraformType, IdAttribute, id)
			return ctx, options, diag.Diagnostics{}
		},
	}
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

type ListStringSchemaAttribute[Model any, Request any, Response any] struct {
	DataSourceExistence AttributeExistence
	DeprecationMessage  string
	Description         string
	MarkdownDescription string
	ReadResponse        func(context.Context, string, string, Response, Model) (context.Context, Model, diag.Diagnostics)
	ResourceExistence   AttributeExistence
	Sensitive           bool
	UpsertRequest       func(context.Context, string, Request, Model) (context.Context, Request, diag.Diagnostics)
	Validators          []validator.List
}

func (a ListStringSchemaAttribute[Model, Request, Response]) Read(
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

func (a ListStringSchemaAttribute[Model, Request, Response]) ToDataSource() datasourceschema.Attribute {
	return datasourceschema.ListAttribute{
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

func (a ListStringSchemaAttribute[Model, Request, Response]) ToResource() resourceschema.Attribute {
	return resourceschema.ListAttribute{
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

func (a ListStringSchemaAttribute[Model, Request, Response]) Upsert(
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
) func(context.Context, string, string, lucirpc.Options, Model) (context.Context, Model, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		section lucirpc.Options,
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
) func(context.Context, string, string, lucirpc.Options, Model) (context.Context, Model, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		section lucirpc.Options,
		model Model,
	) (context.Context, Model, diag.Diagnostics) {
		ctx, value, diagnostics := GetOptionInt64(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
		set(&model, value)
		return ctx, model, diagnostics
	}
}

func ReadResponseOptionListString[Model any](
	set func(*Model, types.List),
	attribute string,
	option string,
) func(context.Context, string, string, lucirpc.Options, Model) (context.Context, Model, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		section lucirpc.Options,
		model Model,
	) (context.Context, Model, diag.Diagnostics) {
		ctx, value, diagnostics := GetOptionListString(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
		set(&model, value)
		return ctx, model, diagnostics
	}
}

func ReadResponseOptionSetString[Model any](
	set func(*Model, types.Set),
	attribute string,
	option string,
) func(context.Context, string, string, lucirpc.Options, Model) (context.Context, Model, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		section lucirpc.Options,
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
) func(context.Context, string, string, lucirpc.Options, Model) (context.Context, Model, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		terraformType string,
		section lucirpc.Options,
		model Model,
	) (context.Context, Model, diag.Diagnostics) {
		ctx, value, diagnostics := GetOptionString(ctx, fullTypeName, terraformType, section, path.Root(attribute), option)
		set(&model, value)
		return ctx, model, diagnostics
	}
}

func RequiredIfAttributeNotEqualBool(
	expression path.Expression,
	expected bool,
) requiredIfAttributeNot[bool] {
	return requiredIfAttributeNotEqual(
		types.BoolType,
		expression,
		expected,
	)
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
) func(context.Context, string, lucirpc.Options, Model) (context.Context, lucirpc.Options, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		options lucirpc.Options,
		model Model,
	) (context.Context, lucirpc.Options, diag.Diagnostics) {
		str := get(model)
		if !hasValue(str) {
			return ctx, options, diag.Diagnostics{}
		}

		ctx = logger.SetFieldBool(ctx, fullTypeName, ResourceTerraformType, attribute, str)
		options[option] = lucirpc.Boolean(str.ValueBool())
		return ctx, options, diag.Diagnostics{}
	}
}

func UpsertRequestOptionInt64[Model any](
	get func(Model) types.Int64,
	attribute string,
	option string,
) func(context.Context, string, lucirpc.Options, Model) (context.Context, lucirpc.Options, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		options lucirpc.Options,
		model Model,
	) (context.Context, lucirpc.Options, diag.Diagnostics) {
		str := get(model)
		if !hasValue(str) {
			return ctx, options, diag.Diagnostics{}
		}

		ctx = logger.SetFieldInt64(ctx, fullTypeName, ResourceTerraformType, attribute, str)
		options[option] = lucirpc.Integer(int(str.ValueInt64()))
		return ctx, options, diag.Diagnostics{}
	}
}

func UpsertRequestOptionListString[Model any](
	get func(Model) types.List,
	attribute string,
	option string,
) func(context.Context, string, lucirpc.Options, Model) (context.Context, lucirpc.Options, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		options lucirpc.Options,
		model Model,
	) (context.Context, lucirpc.Options, diag.Diagnostics) {
		str := get(model)
		if !hasValue(str) {
			return ctx, options, diag.Diagnostics{}
		}

		value, diagnostics := serializeListString(ctx, str, path.Root(attribute))
		if diagnostics.HasError() {
			return ctx, options, diagnostics
		}

		ctx = logger.SetFieldListString(ctx, fullTypeName, ResourceTerraformType, attribute, str)
		options[option] = value
		return ctx, options, diag.Diagnostics{}
	}
}

func UpsertRequestOptionSetString[Model any](
	get func(Model) types.Set,
	attribute string,
	option string,
) func(context.Context, string, lucirpc.Options, Model) (context.Context, lucirpc.Options, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		options lucirpc.Options,
		model Model,
	) (context.Context, lucirpc.Options, diag.Diagnostics) {
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
) func(context.Context, string, lucirpc.Options, Model) (context.Context, lucirpc.Options, diag.Diagnostics) {
	return func(
		ctx context.Context,
		fullTypeName string,
		options lucirpc.Options,
		model Model,
	) (context.Context, lucirpc.Options, diag.Diagnostics) {
		str := get(model)
		if !hasValue(str) {
			return ctx, options, diag.Diagnostics{}
		}

		ctx = logger.SetFieldString(ctx, fullTypeName, ResourceTerraformType, attribute, str)
		options[option] = lucirpc.String(str.ValueString())
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

type requiredIfAttributeNot[Value any] struct {
	attrType   attr.Type
	expected   Value
	expression path.Expression
}

func (a requiredIfAttributeNot[Value]) Description(ctx context.Context) string {
	return a.MarkdownDescription(ctx)
}

func (a requiredIfAttributeNot[Value]) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Ensures that an attribute is set, if %q is also set to %v", a.expression, a.expected)
}

func (a requiredIfAttributeNot[Value]) ValidateBool(
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

func (a requiredIfAttributeNot[Value]) ValidateInt64(
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

func (a requiredIfAttributeNot[Value]) ValidateList(
	ctx context.Context,
	req validator.ListRequest,
	res *validator.ListResponse,
) {
	diagnostics := a.validate(ctx, req.Config, req.Path, req.ConfigValue)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

func (a requiredIfAttributeNot[Value]) ValidateSet(
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

func (a requiredIfAttributeNot[Value]) ValidateString(
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

func (a requiredIfAttributeNot[Value]) validate(
	ctx context.Context,
	config tfsdk.Config,
	requestPath path.Path,
	configValue interface{ IsNull() bool },
) (allDiagnostics diag.Diagnostics) {
	if !configValue.IsNull() {
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
			// If the value is null,
			// it cannot be what we expect.
			// We ignore the value.
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
				diag.NewAttributeErrorDiagnostic(
					requestPath,
					"Missing required argument",
					fmt.Sprintf("Attribute %q is required when %q is not %v", requestPath, matchedPath, expected),
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

func requiredIfAttributeNotEqual[Value any](
	attrType attr.Type,
	expression path.Expression,
	expected Value,
) requiredIfAttributeNot[Value] {
	return requiredIfAttributeNot[Value]{
		attrType:   attrType,
		expected:   expected,
		expression: expression,
	}
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

func (a requiresAttribute[Value]) ValidateList(
	ctx context.Context,
	req validator.ListRequest,
	res *validator.ListResponse,
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

func serializeListString(
	ctx context.Context,
	attribute interface{ Elements() []attr.Value },
	attributePath path.Path,
) (lucirpc.Option, diag.Diagnostics) {
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

	return lucirpc.ListString(values), allDiagnostics
}

func serializeSetString(
	ctx context.Context,
	attribute interface{ Elements() []attr.Value },
	attributePath path.Path,
) (lucirpc.Option, diag.Diagnostics) {
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

	return lucirpc.ListString(values), allDiagnostics
}
