package lucirpcglue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

func ReadModel[Model any](
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	client lucirpc.Client,
	attributes map[string]SchemaAttribute[Model, map[string]json.RawMessage, map[string]json.RawMessage],
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
