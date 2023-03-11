package network

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

var (
	_ datasource.DataSource              = &globalsDataSource{}
	_ datasource.DataSourceWithConfigure = &globalsDataSource{}
)

func NewGlobalsDataSource() datasource.DataSource {
	return &globalsDataSource{}
}

type globalsDataSource struct {
	client        lucirpc.Client
	fullTypeName  string
	terraformType string
}

// Configure prepares the data source.
func (d *globalsDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	res *datasource.ConfigureResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Configuring %s Data Source", d.fullTypeName))
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

// Metadata sets the data source name.
func (d *globalsDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	res *datasource.MetadataResponse,
) {
	fullTypeName := fmt.Sprintf("%s_%s", req.ProviderTypeName, globalsTypeName)
	d.fullTypeName = fullTypeName
	d.terraformType = lucirpcglue.DataSourceTerraformType
	res.TypeName = fullTypeName
}

// Read prepares the data source.
func (d *globalsDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	res *datasource.ReadResponse,
) {
	tflog.Info(ctx, fmt.Sprintf("Reading %s data source", d.fullTypeName))

	tflog.Debug(ctx, "Retrieving values from config")
	var model globalsModel
	diagnostics := req.Config.Get(ctx, &model)
	res.Diagnostics.Append(diagnostics...)
	if res.Diagnostics.HasError() {
		return
	}

	ctx, model, diagnostics = lucirpcglue.ReadModel(
		ctx,
		d.fullTypeName,
		d.terraformType,
		d.client,
		globalsSchemaAttributes,
		globalsUCIConfig,
		model.Id.ValueString(),
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

// Schema prepares the data source.
func (d *globalsDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	res *datasource.SchemaResponse,
) {
	attributes := map[string]schema.Attribute{}
	for k, v := range globalsSchemaAttributes {
		attributes[k] = v.ToDataSource()
	}

	res.Schema = schema.Schema{
		Attributes:  attributes,
		Description: globalsSchemaDescription,
	}
}
