package system

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	resourceTerraformType = "resource"
)

var (
	_ resource.Resource                = &systemResource{}
	_ resource.ResourceWithConfigure   = &systemResource{}
	_ resource.ResourceWithImportState = &systemResource{}
)

func NewSystemResource() resource.Resource {
	return &systemResource{}
}

type systemResource struct {
	client       lucirpc.Client
	fullTypeName string
}

// Configure adds the provider configured client to the resource.
func (d *systemResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	res *resource.ConfigureResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Configuring %s Resource", d.fullTypeName))
	if req.ProviderData == nil {
		tflog.Debug(ctx, "No provider data")
		return
	}

	client, diagnostics := lucirpcglue.NewClient(lucirpcglue.ConfigureRequest(req))
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	d.client = *client
}

// Create constructs a new resource and sets the initial Terraform state.
func (d *systemResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	res *resource.CreateResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Creating %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Retrieving values from plan")
	var model systemModel
	diagnostics := req.Plan.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Generating API request body")
	options := map[string]json.RawMessage{}

	tflog.Debug(ctx, "Handling required attributes")
	ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemIdAttribute, model.Id)
	id := model.Id.ValueString()

	tflog.Debug(ctx, "Handling optional attributes")
	if hasValue(model.ConLogLevel) {
		value, diagnostics := serializeInt64(model.ConLogLevel, path.Root(systemConLogLevelAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldInt64(ctx, d.fullTypeName, resourceTerraformType, systemConLogLevelAttribute, model.ConLogLevel)
			options[systemConLogLevelUCIOption] = value
		}
	}

	if hasValue(model.CronLogLevel) {
		value, diagnostics := serializeInt64(model.CronLogLevel, path.Root(systemCronLogLevelAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldInt64(ctx, d.fullTypeName, resourceTerraformType, systemCronLogLevelAttribute, model.CronLogLevel)
			options[systemCronLogLevelUCIOption] = value
		}
	}

	if hasValue(model.Description) {
		value, diagnostics := serializeString(model.Description, path.Root(systemDescriptionAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemDescriptionAttribute, model.Description)
			options[systemDescriptionUCIOption] = value
		}
	}

	if hasValue(model.Hostname) {
		value, diagnostics := serializeString(model.Hostname, path.Root(systemHostnameAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemHostnameAttribute, model.Hostname)
			options[systemHostnameUCIOption] = value
		}
	}

	if hasValue(model.LogSize) {
		value, diagnostics := serializeInt64(model.LogSize, path.Root(systemLogSizeAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldInt64(ctx, d.fullTypeName, resourceTerraformType, systemLogSizeAttribute, model.LogSize)
			options[systemLogSizeUCIOption] = value
		}
	}

	if hasValue(model.Notes) {
		value, diagnostics := serializeString(model.Notes, path.Root(systemNotesAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemNotesAttribute, model.Notes)
			options[systemNotesUCIOption] = value
		}
	}

	if hasValue(model.Timezone) {
		value, diagnostics := serializeString(model.Timezone, path.Root(systemTimezoneAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemTimezoneAttribute, model.Timezone)
			options[systemTimezoneUCIOption] = value
		}
	}

	if hasValue(model.TTYLogin) {
		value, diagnostics := serializeBool(model.TTYLogin, path.Root(systemTTYLoginAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldBool(ctx, d.fullTypeName, resourceTerraformType, systemTTYLoginAttribute, model.TTYLogin)
			options[systemTTYLoginUCIOption] = value
		}
	}

	if hasValue(model.Zonename) {
		value, diagnostics := serializeString(model.Zonename, path.Root(systemZonenameAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemZonenameAttribute, model.Zonename)
			options[systemZonenameUCIOption] = value
		}
	}

	if res.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", systemUCIConfig, systemUCISection))
	ok, diagnostics := lucirpcglue.CreateSection(
		ctx,
		d.client,
		systemUCIConfig,
		systemUCISystemType,
		id,
		options,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	if !ok {
		res.Diagnostics.AddError(
			fmt.Sprintf("Could not create %s.%s section", systemUCIConfig, systemUCISection),
			"It is not currently known why this happens. It is unclear if this is a problem with the provider. Please double check the values provided are acceptable.",
		)
		return
	}

	tflog.Debug(ctx, "Reading updated section")
	ctx, model, diagnostics = ReadModel(
		ctx,
		d.fullTypeName,
		resourceTerraformType,
		d.client,
		id,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating state with values")
	diagnostics = res.State.Set(ctx, model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

// Delete removes the actual resource and remove the Terraform state on success.
func (d *systemResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	res *resource.DeleteResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Deleting %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Getting the current state")
	var model systemModel
	diagnostics := req.State.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemIdAttribute, model.Id)
	id := model.Id.ValueString()
	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", systemUCIConfig, id))
	tflog.Debug(ctx, "Deleting existing section")
	ok, diagnostics := lucirpcglue.DeleteSection(
		ctx,
		d.client,
		systemUCIConfig,
		id,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	if !ok {
		res.Diagnostics.AddError(
			fmt.Sprintf("Could not delete %s.%s section", systemUCIConfig, systemUCISection),
			"It is not currently known why this happens. It is unclear if this is a problem with the provider. Please double check the values provided are acceptable.",
		)
		return
	}
}

// ImportState brings an existing resource into Terraform state.
func (d *systemResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	res *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Retrieving import id and saving to id attribute")
	resource.ImportStatePassthroughID(ctx, path.Root(systemIdAttribute), req, res)
}

// Metadata sets the resource type name.
func (d *systemResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	res *resource.MetadataResponse,
) {
	fullTypeName := fmt.Sprintf("%s_%s", req.ProviderTypeName, systemTypeName)
	d.fullTypeName = fullTypeName
	res.TypeName = fullTypeName
}

// Read refreshes the Terraform state with the latest data.
func (d *systemResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	res *resource.ReadResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Reading %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Getting the current state")
	var model systemModel
	diagnostics := req.State.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx, model, diagnostics = ReadModel(
		ctx,
		d.fullTypeName,
		resourceTerraformType,
		d.client,
		model.Id.ValueString(),
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Setting the %s resource state", d.fullTypeName))
	diagnostics = res.State.Set(ctx, model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

// Schema defines the schema for the resource.
func (d *systemResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	res *resource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{}
	for k, v := range systemSchemaAttributes {
		attributes[k] = v.ToResource()
	}

	res.Schema = schema.Schema{
		Attributes:  attributes,
		Description: "Provides system data about an OpenWrt device",
	}
}

// Update modifies part of the resource and sets the Terraform state on success.
func (d *systemResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	res *resource.UpdateResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Updating %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Retrieving values from plan")
	var model systemModel
	diagnostics := req.Plan.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Generating API request body")
	options := map[string]json.RawMessage{}

	tflog.Debug(ctx, "Handling required attributes")
	ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemIdAttribute, model.Id)
	id := model.Id.ValueString()

	tflog.Debug(ctx, "Handling optional attributes")
	if hasValue(model.ConLogLevel) {
		value, diagnostics := serializeInt64(model.ConLogLevel, path.Root(systemConLogLevelAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldInt64(ctx, d.fullTypeName, resourceTerraformType, systemConLogLevelAttribute, model.ConLogLevel)
			options[systemConLogLevelUCIOption] = value
		}
	}

	if hasValue(model.CronLogLevel) {
		value, diagnostics := serializeInt64(model.CronLogLevel, path.Root(systemCronLogLevelAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldInt64(ctx, d.fullTypeName, resourceTerraformType, systemCronLogLevelAttribute, model.CronLogLevel)
			options[systemCronLogLevelUCIOption] = value
		}
	}

	if hasValue(model.Description) {
		value, diagnostics := serializeString(model.Description, path.Root(systemDescriptionAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemDescriptionAttribute, model.Description)
			options[systemDescriptionUCIOption] = value
		}
	}

	if hasValue(model.Hostname) {
		value, diagnostics := serializeString(model.Hostname, path.Root(systemHostnameAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemHostnameAttribute, model.Hostname)
			options[systemHostnameUCIOption] = value
		}
	}

	if hasValue(model.LogSize) {
		value, diagnostics := serializeInt64(model.LogSize, path.Root(systemLogSizeAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldInt64(ctx, d.fullTypeName, resourceTerraformType, systemLogSizeAttribute, model.LogSize)
			options[systemLogSizeUCIOption] = value
		}
	}

	if hasValue(model.Notes) {
		value, diagnostics := serializeString(model.Notes, path.Root(systemNotesAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemNotesAttribute, model.Notes)
			options[systemNotesUCIOption] = value
		}
	}

	if hasValue(model.Timezone) {
		value, diagnostics := serializeString(model.Timezone, path.Root(systemTimezoneAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemTimezoneAttribute, model.Timezone)
			options[systemTimezoneUCIOption] = value
		}
	}

	if hasValue(model.TTYLogin) {
		value, diagnostics := serializeBool(model.TTYLogin, path.Root(systemTTYLoginAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldBool(ctx, d.fullTypeName, resourceTerraformType, systemTTYLoginAttribute, model.TTYLogin)
			options[systemTTYLoginUCIOption] = value
		}
	}

	if hasValue(model.Zonename) {
		value, diagnostics := serializeString(model.Zonename, path.Root(systemZonenameAttribute))
		res.Diagnostics.Append(diagnostics...)
		if !res.Diagnostics.HasError() {
			ctx = logger.SetFieldString(ctx, d.fullTypeName, resourceTerraformType, systemZonenameAttribute, model.Zonename)
			options[systemZonenameUCIOption] = value
		}
	}

	if res.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", systemUCIConfig, systemUCISection))
	ok, diagnostics := lucirpcglue.UpdateSection(
		ctx,
		d.client,
		systemUCIConfig,
		id,
		options,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	if !ok {
		res.Diagnostics.AddError(
			fmt.Sprintf("Could not create %s.%s section", systemUCIConfig, systemUCISection),
			"It is not currently known why this happens. It is unclear if this is a problem with the provider. Please double check the values provided are acceptable.",
		)
		return
	}

	tflog.Debug(ctx, "Reading updated section")
	ctx, model, diagnostics = ReadModel(
		ctx,
		d.fullTypeName,
		resourceTerraformType,
		d.client,
		id,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating state with values")
	diagnostics = res.State.Set(ctx, model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
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
