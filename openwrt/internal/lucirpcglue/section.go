package lucirpcglue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

// GetMetadataString attempts to parse the given metadata key from the section.
// Any diagnostic information found in the process (including errors) is returned.
func GetSection(
	ctx context.Context,
	client lucirpc.Client,
	config string,
	section string,
) (map[string]json.RawMessage, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result, err := client.GetSection(ctx, config, section)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("problem getting %s.%s section", config, section),
			err.Error(),
		)
		return map[string]json.RawMessage{}, diagnostics
	}

	return result, diagnostics
}
