package lucirpcglue

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
)

// GetOptionBool attempts to parse the given option from the section as a bool.
// Any diagnostic information found in the process (including errors) is returned.
func GetOptionBool(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section map[string]json.RawMessage,
	attribute path.Path,
	option string,
) (context.Context, types.Bool, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.BoolNull()
	raw, ok := section[option]
	if !ok {
		return ctx, result, diagnostics
	}

	// Booleans in UCI can be any number of things:
	// - True: "1", "yes", "on", "true", "enabled"
	// - False: "0", "no", "off", "false", "disabled"
	// We try to parse on of these out of the string.
	var boolish string
	err := json.Unmarshal(raw, &boolish)
	if err != nil {
		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	switch boolish {
	case "1", "yes", "on", "true", "enabled":
		result = types.BoolValue(true)

	case "0", "no", "off", "false", "disabled":
		result = types.BoolValue(false)

	default:
		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("Unexpected value for option: %q", option),
			fmt.Sprintf(`expected one of "1", "yes", "on", "true", "enabled", "0", "no", "off", "false", or "disabled"; got: %q`, boolish),
		)
		return ctx, result, diagnostics
	}

	ctx = logger.SetFieldBool(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, diagnostics
}

// GetOptionInt64 attempts to parse the given option from the section as an int64.
// Any diagnostic information found in the process (including errors) is returned.
func GetOptionInt64(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section map[string]json.RawMessage,
	attribute path.Path,
	option string,
) (context.Context, types.Int64, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.Int64Null()
	raw, ok := section[option]
	if !ok {
		return ctx, result, diagnostics
	}

	// Integers in UCI are stored as strtings.
	// We have to unmarshall first, then parse the string.
	var intish string
	err := json.Unmarshal(raw, &intish)
	if err != nil {
		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to parse option: %q", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	value, err := strconv.Atoi(intish)
	if err != nil {
		diagnostics.AddAttributeError(
			attribute,
			fmt.Sprintf("unable to convert option: %q to a string", option),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	result = types.Int64Value(int64(value))
	ctx = logger.SetFieldInt64(ctx, fullTypeName, terraformType, option, result)
	return ctx, result, diagnostics
}

// GetOptionString attempts to parse the given option from the section as a string.
// Any diagnostic information found in the process (including errors) is returned.
func GetOptionString(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section map[string]json.RawMessage,
	attribute path.Path,
	option string,
) (context.Context, types.String, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.StringNull()
	raw, ok := section[option]
	if !ok {
		return ctx, result, diagnostics
	}

	var value string
	err := json.Unmarshal(raw, &value)
	if err != nil {
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
