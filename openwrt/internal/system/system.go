package system

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	systemConLogLevelAttribute = "conloglevel"
	systemConLogLevelUCIOption = "conloglevel"

	systemCronLogLevelAttribute = "cronloglevel"
	systemCronLogLevelUCIOption = "cronloglevel"

	systemDescriptionAttribute = "description"
	systemDescriptionUCIOption = "description"

	systemHostnameAttribute = "hostname"
	systemHostnameUCIOption = "hostname"

	systemIdAttribute  = "id"
	systemIdUCISection = ".name"

	systemLogSizeAttribute = "log_size"
	systemLogSizeUCIOption = "log_size"

	systemNotesAttribute = "notes"
	systemNotesUCIOption = "notes"

	systemTimezoneAttribute = "timezone"
	systemTimezoneUCIOption = "timezone"

	systemTTYLoginAttribute = "ttylogin"
	systemTTYLoginUCIOption = "ttylogin"

	systemTypeName   = "system_system"
	systemUCIConfig  = "system"
	systemUCISection = "@system[0]"

	systemZonenameAttribute = "zonename"
	systemZonenameUCIOption = "zonename"
)

type systemModel struct {
	ConLogLevel  types.Int64  `tfsdk:"conloglevel"`
	CronLogLevel types.Int64  `tfsdk:"cronloglevel"`
	Description  types.String `tfsdk:"description"`
	Hostname     types.String `tfsdk:"hostname"`
	Id           types.String `tfsdk:"id"`
	LogSize      types.Int64  `tfsdk:"log_size"`
	Notes        types.String `tfsdk:"notes"`
	Timezone     types.String `tfsdk:"timezone"`
	TTYLogin     types.Bool   `tfsdk:"ttylogin"`
	Zonename     types.String `tfsdk:"zonename"`
}

func ReadModel(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	client lucirpc.Client,
) (context.Context, systemModel, diag.Diagnostics) {
	tflog.Info(ctx, "Reading system model")
	var (
		allDiagnostics diag.Diagnostics
		model          systemModel
	)

	section, diagnostics := lucirpcglue.GetSection(ctx, client, systemUCIConfig, systemUCISection)
	allDiagnostics.Append(diagnostics...)
	if allDiagnostics.HasError() {
		return ctx, model, allDiagnostics
	}

	ctx, model.ConLogLevel, diagnostics = lucirpcglue.GetOptionInt64(ctx, fullTypeName, dataSourceTerraformType, section, path.Root(systemConLogLevelAttribute), systemConLogLevelUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.CronLogLevel, diagnostics = lucirpcglue.GetOptionInt64(ctx, fullTypeName, dataSourceTerraformType, section, path.Root(systemCronLogLevelAttribute), systemCronLogLevelUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Description, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, dataSourceTerraformType, section, path.Root(systemDescriptionAttribute), systemDescriptionUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Hostname, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, dataSourceTerraformType, section, path.Root(systemHostnameAttribute), systemHostnameUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.LogSize, diagnostics = lucirpcglue.GetOptionInt64(ctx, fullTypeName, dataSourceTerraformType, section, path.Root(systemLogSizeAttribute), systemLogSizeUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Notes, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, dataSourceTerraformType, section, path.Root(systemNotesAttribute), systemNotesUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Timezone, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, dataSourceTerraformType, section, path.Root(systemTimezoneAttribute), systemTimezoneUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.TTYLogin, diagnostics = lucirpcglue.GetOptionBool(ctx, fullTypeName, dataSourceTerraformType, section, path.Root(systemTTYLoginAttribute), systemTTYLoginUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Zonename, diagnostics = lucirpcglue.GetOptionString(ctx, fullTypeName, dataSourceTerraformType, section, path.Root(systemZonenameAttribute), systemZonenameUCIOption)
	allDiagnostics.Append(diagnostics...)
	ctx, model.Id, diagnostics = lucirpcglue.GetMetadataString(ctx, fullTypeName, dataSourceTerraformType, section, systemIdUCISection)
	allDiagnostics.Append(diagnostics...)

	return ctx, model, diagnostics
}
