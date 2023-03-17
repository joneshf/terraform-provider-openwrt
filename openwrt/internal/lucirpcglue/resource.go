package lucirpcglue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
)

var (
	_ resource.Resource                = &deviceResource[any]{}
	_ resource.ResourceWithConfigure   = &deviceResource[any]{}
	_ resource.ResourceWithImportState = &deviceResource[any]{}
)

func NewResource[Model any](
	getId func(Model) types.String,
	schemaAttributes map[string]SchemaAttribute[Model, map[string]json.RawMessage, map[string]json.RawMessage],
	schemaDescription string,
	uciConfig string,
	uciType string,
) resource.Resource {
	return &deviceResource[Model]{
		getId:             getId,
		schemaAttributes:  schemaAttributes,
		schemaDescription: schemaDescription,
		terraformType:     ResourceTerraformType,
		uciConfig:         uciConfig,
		uciType:           uciType,
	}
}

type deviceResource[Model any] struct {
	client            lucirpc.Client
	fullTypeName      string
	getId             func(Model) types.String
	schemaAttributes  map[string]SchemaAttribute[Model, map[string]json.RawMessage, map[string]json.RawMessage]
	schemaDescription string
	terraformType     string
	uciConfig         string
	uciType           string
}

// Configure adds the provider configured client to the resource.
func (d *deviceResource[Model]) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	res *resource.ConfigureResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Configuring %s.%s resource", d.uciConfig, d.uciType))
	if req.ProviderData == nil {
		tflog.Debug(ctx, "No provider data")
		return
	}

	providerData, diagnostics := ParseProviderData(ConfigureRequest(req))
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	d.client = providerData.Client
	d.fullTypeName = d.getFullTypeName(providerData.TypeName)
}

// Create constructs a new resource and sets the initial Terraform state.
func (d *deviceResource[Model]) Create(
	ctx context.Context,
	req resource.CreateRequest,
	res *resource.CreateResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Creating %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Retrieving values from plan")
	var model Model
	diagnostics := req.Plan.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx, options, diagnostics := GenerateUpsertBody(
		ctx,
		d.fullTypeName,
		model,
		d.schemaAttributes,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	id := d.getId(model).ValueString()
	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", d.uciConfig, id))
	diagnostics = CreateSection(
		ctx,
		d.client,
		d.uciConfig,
		d.uciType,
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
		d.terraformType,
		d.client,
		d.schemaAttributes,
		d.uciConfig,
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
func (d *deviceResource[Model]) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	res *resource.DeleteResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Deleting %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Getting the current state")
	var model Model
	diagnostics := req.State.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx = logger.SetFieldString(ctx, d.fullTypeName, d.terraformType, IdAttribute, d.getId(model))
	id := d.getId(model).ValueString()
	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", d.uciConfig, id))
	tflog.Debug(ctx, "Deleting existing section")
	diagnostics = DeleteSection(
		ctx,
		d.client,
		d.uciConfig,
		id,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

// ImportState brings an existing resource into Terraform state.
func (d *deviceResource[Model]) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	res *resource.ImportStateResponse,
) {
	tflog.Debug(ctx, "Retrieving import id and saving to id attribute")
	resource.ImportStatePassthroughID(ctx, path.Root(IdAttribute), req, res)
}

// Metadata sets the resource type name.
func (d *deviceResource[Model]) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	res *resource.MetadataResponse,
) {
	res.TypeName = d.getFullTypeName(req.ProviderTypeName)
}

// Read refreshes the Terraform state with the latest data.
func (d *deviceResource[Model]) Read(
	ctx context.Context,
	req resource.ReadRequest,
	res *resource.ReadResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Reading %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Getting the current state")
	var model Model
	diagnostics := req.State.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx, model, diagnostics = ReadModel(
		ctx,
		d.fullTypeName,
		d.terraformType,
		d.client,
		d.schemaAttributes,
		d.uciConfig,
		d.getId(model).ValueString(),
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
func (d *deviceResource[Model]) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	res *resource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{}
	for k, v := range d.schemaAttributes {
		attributes[k] = v.ToResource()
	}

	res.Schema = schema.Schema{
		Attributes:  attributes,
		Description: d.schemaDescription,
	}
}

// Update modifies part of the resource and sets the Terraform state on success.
func (d *deviceResource[Model]) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	res *resource.UpdateResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Updating %s resource", d.fullTypeName))

	tflog.Debug(ctx, "Retrieving values from plan")
	var model Model
	diagnostics := req.Plan.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx, options, diagnostics := GenerateUpsertBody(
		ctx,
		d.fullTypeName,
		model,
		d.schemaAttributes,
	)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	id := d.getId(model).ValueString()
	ctx = tflog.SetField(ctx, "section", fmt.Sprintf("%s.%s", d.uciConfig, id))
	diagnostics = UpdateSection(
		ctx,
		d.client,
		d.uciConfig,
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
		d.terraformType,
		d.client,
		d.schemaAttributes,
		d.uciConfig,
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

func (d deviceResource[Model]) getFullTypeName(
	providerTypeName string,
) string {
	return fmt.Sprintf("%s_%s_%s", providerTypeName, d.uciConfig, d.uciType)
}
