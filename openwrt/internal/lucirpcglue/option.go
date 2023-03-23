package lucirpcglue

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
)

// GetOptionBool attempts to parse the given option from the section as a bool.
// Any diagnostic information found in the process (including errors) is returned.
func GetOptionBool(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section lucirpc.Options,
	attribute path.Path,
	option string,
) (context.Context, types.Bool, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.BoolNull()
	value, err := section.GetBoolean(option)
	if err != nil {
		if errors.As(err, &lucirpc.OptionNotFoundError{}) {
			return ctx, result, diagnostics
		}

		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	result = types.BoolValue(value)
	ctx = logger.SetFieldBool(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, diagnostics
}

// GetOptionInt64 attempts to parse the given option from the section as an int64.
// Any diagnostic information found in the process (including errors) is returned.
func GetOptionInt64(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section lucirpc.Options,
	attribute path.Path,
	option string,
) (context.Context, types.Int64, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.Int64Null()
	value, err := section.GetInteger(option)
	if err != nil {
		if errors.As(err, &lucirpc.OptionNotFoundError{}) {
			return ctx, result, diagnostics
		}

		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	result = types.Int64Value(int64(value))
	ctx = logger.SetFieldInt64(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, diagnostics
}

// GetOptionListString attempts to parse the given option from the section as a []string.
// Any diagnostic information found in the process (including errors) is returned.
func GetOptionListString(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section lucirpc.Options,
	attribute path.Path,
	option string,
) (context.Context, types.List, diag.Diagnostics) {
	allDiagnostics := diag.Diagnostics{}
	result := types.ListNull(types.StringType)
	values, err := section.GetListString(option)
	if err != nil {
		if errors.As(err, &lucirpc.OptionNotFoundError{}) {
			return ctx, result, allDiagnostics
		}

		allDiagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, allDiagnostics
	}

	var attrValues []attr.Value
	for _, value := range values {
		var attrValue attr.Value
		diagnostics := tfsdk.ValueFrom(ctx, value, types.StringType, &attrValue)
		allDiagnostics.Append(diagnostics...)
		if allDiagnostics.HasError() {
			// We don't want to exit early.
			// We want to continue to accumulate diagnostics.
			continue
		}

		attrValues = append(attrValues, attrValue)
	}

	if allDiagnostics.HasError() {
		return ctx, result, allDiagnostics
	}

	result, allDiagnostics = types.ListValue(types.StringType, attrValues)
	if allDiagnostics.HasError() {
		return ctx, result, allDiagnostics
	}

	ctx = logger.SetFieldListString(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, allDiagnostics
}

// GetOptionSetString attempts to parse the given option from the section as a []string.
// Any diagnostic information found in the process (including errors) is returned.
func GetOptionSetString(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section lucirpc.Options,
	attribute path.Path,
	option string,
) (context.Context, types.Set, diag.Diagnostics) {
	allDiagnostics := diag.Diagnostics{}
	result := types.SetNull(types.StringType)
	values, err := section.GetListString(option)
	if err != nil {
		if errors.As(err, &lucirpc.OptionNotFoundError{}) {
			return ctx, result, allDiagnostics
		}

		allDiagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, allDiagnostics
	}

	var attrValues []attr.Value
	for _, value := range values {
		var attrValue attr.Value
		diagnostics := tfsdk.ValueFrom(ctx, value, types.StringType, &attrValue)
		allDiagnostics.Append(diagnostics...)
		if allDiagnostics.HasError() {
			// We don't want to exit early.
			// We want to continue to accumulate diagnostics.
			continue
		}

		attrValues = append(attrValues, attrValue)
	}

	if allDiagnostics.HasError() {
		return ctx, result, allDiagnostics
	}

	result, allDiagnostics = types.SetValue(types.StringType, attrValues)
	if allDiagnostics.HasError() {
		return ctx, result, allDiagnostics
	}

	ctx = logger.SetFieldSetString(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, allDiagnostics
}

// GetOptionString attempts to parse the given option from the section as a string.
// Any diagnostic information found in the process (including errors) is returned.
func GetOptionString(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section lucirpc.Options,
	attribute path.Path,
	option string,
) (context.Context, types.String, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.StringNull()
	value, err := section.GetString(option)
	if err != nil {
		if errors.As(err, &lucirpc.OptionNotFoundError{}) {
			return ctx, result, diagnostics
		}

		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	result = types.StringValue(value)
	ctx = logger.SetFieldString(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, diagnostics
}
