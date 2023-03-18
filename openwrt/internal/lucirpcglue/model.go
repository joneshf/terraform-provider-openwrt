package lucirpcglue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

func GenerateUpsertBody[Model any](
	ctx context.Context,
	fullTypeName string,
	model Model,
	attributes map[string]SchemaAttribute[Model, lucirpc.Options, lucirpc.Options],
) (context.Context, lucirpc.Options, diag.Diagnostics) {
	tflog.Info(ctx, "Generating API request body")
	var diagnostics diag.Diagnostics
	allDiagnostics := diag.Diagnostics{}
	options := lucirpc.Options{}

	tflog.Debug(ctx, "Handling attributes")
	for _, attribute := range attributes {
		ctx, options, diagnostics = attribute.Upsert(ctx, fullTypeName, options, model)
		allDiagnostics.Append(diagnostics...)
	}

	return ctx, options, allDiagnostics
}

func ReadModel[Model any](
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	client lucirpc.Client,
	attributes map[string]SchemaAttribute[Model, lucirpc.Options, lucirpc.Options],
	uciConfig string,
	uciSection string,
) (context.Context, Model, diag.Diagnostics) {
	tflog.Info(ctx, fmt.Sprintf("Reading %s model", fullTypeName))
	var (
		allDiagnostics diag.Diagnostics
		model          Model
	)

	section, diagnostics := GetSection(ctx, client, uciConfig, uciSection)
	allDiagnostics.Append(diagnostics...)
	if allDiagnostics.HasError() {
		return ctx, model, allDiagnostics
	}

	for _, attribute := range attributes {
		ctx, model, diagnostics = attribute.Read(ctx, fullTypeName, terraformType, section, model)
		allDiagnostics.Append(diagnostics...)
	}

	return ctx, model, diagnostics
}
