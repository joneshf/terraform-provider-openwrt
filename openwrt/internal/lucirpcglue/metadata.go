package lucirpcglue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
)

// GetMetadataString attempts to parse the given metadata key from the section as a string.
// Any diagnostic information found in the process (including errors) is returned.
func GetMetadataString(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	section lucirpc.Options,
	key string,
) (context.Context, types.String, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result := types.StringNull()
	value, err := section.GetString(key)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("unable to parse metadata: %q", key),
			err.Error(),
		)
		return ctx, result, diagnostics
	}

	result = types.StringValue(value)
	ctx = logger.SetFieldString(ctx, fullTypeName, terraformType, key, result)
	return ctx, result, diagnostics
}
