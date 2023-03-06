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

	ctx, id, options, diagnostics := generateAPIBody(
		ctx,
		model,
		d.fullTypeName,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", systemUCIConfig, systemUCISection))
	diagnostics = lucirpcglue.CreateSection(
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
	diagnostics = lucirpcglue.DeleteSection(
		ctx,
		d.client,
		systemUCIConfig,
		id,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
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

	ctx, id, options, diagnostics := generateAPIBody(
		ctx,
		model,
		d.fullTypeName,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", systemUCIConfig, systemUCISection))
	diagnostics = lucirpcglue.UpdateSection(
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

func generateAPIBody(
	ctx context.Context,
	model systemModel,
	fullTypeName string,
) (context.Context, string, map[string]json.RawMessage, diag.Diagnostics) {
	tflog.Info(ctx, "Generating API request body")
	var diagnostics diag.Diagnostics
	allDiagnostics := diag.Diagnostics{}
	options := map[string]json.RawMessage{}

	tflog.Debug(ctx, "Handling attributes")
	id := model.Id.ValueString()
	for _, attribute := range systemSchemaAttributes {
		ctx, options, diagnostics = attribute.Upsert(ctx, fullTypeName, resourceTerraformType, options, model)
		allDiagnostics.Append(diagnostics...)
	}

	return ctx, id, options, allDiagnostics
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
