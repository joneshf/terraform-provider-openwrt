package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

var (
	_ resource.Resource                = &deviceResource{}
	_ resource.ResourceWithConfigure   = &deviceResource{}
	_ resource.ResourceWithImportState = &deviceResource{}
)

func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

type deviceResource struct {
	client        lucirpc.Client
	fullTypeName  string
	terraformType string
}

// Configure adds the provider configured client to the resource.
func (d *deviceResource) Configure(
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
func (d *deviceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	res *resource.CreateResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Creating %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Retrieving values from plan")
	var model deviceModel
	diagnostics := req.Plan.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx, options, diagnostics := model.generateAPIBody(
		ctx,
		d.fullTypeName,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	id := model.Id.ValueString()
	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", deviceUCIConfig, id))
	diagnostics = lucirpcglue.CreateSection(
		ctx,
		d.client,
		deviceUCIConfig,
		deviceUCIType,
		id,
		options,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading updated section")
	ctx, model, diagnostics = lucirpcglue.ReadModel(
		ctx,
		d.fullTypeName,
		d.terraformType,
		d.client,
		deviceSchemaAttributes,
		deviceUCIConfig,
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
func (d *deviceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	res *resource.DeleteResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Deleting %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Getting the current state")
	var model deviceModel
	diagnostics := req.State.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx = logger.SetFieldString(ctx, d.fullTypeName, d.terraformType, deviceIdAttribute, model.Id)
	id := model.Id.ValueString()
	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", deviceUCIConfig, id))
	tflog.Debug(ctx, "Deleting existing section")
	diagnostics = lucirpcglue.DeleteSection(
		ctx,
		d.client,
		deviceUCIConfig,
		id,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

// ImportState brings an existing resource into Terraform state.
func (d *deviceResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	res *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Retrieving import id and saving to id attribute")
	resource.ImportStatePassthroughID(ctx, path.Root(deviceIdAttribute), req, res)
}

// Metadata sets the resource type name.
func (d *deviceResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	res *resource.MetadataResponse,
) {
	fullTypeName := fmt.Sprintf("%s_%s", req.ProviderTypeName, deviceTypeName)
	d.fullTypeName = fullTypeName
	d.terraformType = lucirpcglue.ResourceTerraformType
	res.TypeName = fullTypeName
}

// Read refreshes the Terraform state with the latest data.
func (d *deviceResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	res *resource.ReadResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Reading %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Getting the current state")
	var model deviceModel
	diagnostics := req.State.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx, model, diagnostics = lucirpcglue.ReadModel(
		ctx,
		d.fullTypeName,
		d.terraformType,
		d.client,
		deviceSchemaAttributes,
		deviceUCIConfig,
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
func (d *deviceResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	res *resource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{}
	for k, v := range deviceSchemaAttributes {
		attributes[k] = v.ToResource()
	}

	res.Schema = schema.Schema{
		Attributes:  attributes,
		Description: deviceSchemaDescription,
	}
}

// Update modifies part of the resource and sets the Terraform state on success.
func (d *deviceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	res *resource.UpdateResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Updating %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Retrieving values from plan")
	var model deviceModel
	diagnostics := req.Plan.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx, options, diagnostics := model.generateAPIBody(
		ctx,
		d.fullTypeName,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	id := model.Id.ValueString()
	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", deviceUCIConfig, id))
	diagnostics = lucirpcglue.UpdateSection(
		ctx,
		d.client,
		deviceUCIConfig,
		id,
		options,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading updated section")
	ctx, model, diagnostics = lucirpcglue.ReadModel(
		ctx,
		d.fullTypeName,
		d.terraformType,
		d.client,
		deviceSchemaAttributes,
		deviceUCIConfig,
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
