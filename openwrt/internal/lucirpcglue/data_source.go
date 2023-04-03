package lucirpcglue

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

var (
	_ datasource.DataSource              = &dataSource[any]{}
	_ datasource.DataSourceWithConfigure = &dataSource[any]{}
)

func NewDataSource[Model any](
	getId func(Model) types.String,
	schemaAttributes map[string]SchemaAttribute[Model, lucirpc.Options, lucirpc.Options],
	schemaDescription string,
	uciConfig string,
	uciType string,
) datasource.DataSource {
	return &dataSource[Model]{
		getId:             getId,
		schemaAttributes:  schemaAttributes,
		schemaDescription: schemaDescription,
		terraformType:     DataSourceTerraformType,
		uciConfig:         uciConfig,
		uciType:           uciType,
	}
}

type dataSource[Model any] struct {
	client            lucirpc.Client
	fullTypeName      string
	getId             func(Model) types.String
	schemaAttributes  map[string]SchemaAttribute[Model, lucirpc.Options, lucirpc.Options]
	schemaDescription string
	terraformType     string
	uciConfig         string
	uciType           string
}

// Configure adds the provider configured client to the data source.
func (d *dataSource[Model]) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	res *datasource.ConfigureResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Configuring %s.%s data source", d.uciConfig, d.uciType))
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

// Metadata sets the data source name.
func (d *dataSource[Model]) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	res *datasource.MetadataResponse,
) {
	res.TypeName = d.getFullTypeName(req.ProviderTypeName)
}

// Read refreshes the Terraform state with the latest data.
func (d *dataSource[Model]) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	res *datasource.ReadResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Reading %s data source", d.fullTypeName))

	tflog.Debug(ctx, "Retrieving values from config")
	var model Model
	diagnostics := req.Config.Get(ctx, &model)
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

	tflog.Debug(ctx, fmt.Sprintf("Setting the %s data source state", d.fullTypeName))
	diagnostics = res.State.Set(ctx, model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}
}

// Schema defines the schema for the data source.
func (d *dataSource[Model]) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	res *datasource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{}
	for k, v := range d.schemaAttributes {
		attributes[k] = v.ToDataSource()
	}

	res.Schema = schema.Schema{
		Attributes:  attributes,
		Description: d.schemaDescription,
	}
}

func (d dataSource[Model]) getFullTypeName(
	providerTypeName string,
) string {
	uciConfig := strings.ReplaceAll(d.uciConfig, "-", "_")
	uciType := strings.ReplaceAll(d.uciType, "-", "_")
	return fmt.Sprintf("%s_%s_%s", providerTypeName, uciConfig, uciType)
}
