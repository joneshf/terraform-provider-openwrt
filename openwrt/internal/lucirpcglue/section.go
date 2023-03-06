package lucirpcglue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

// CreateSection attempts to create a new section.
// The bool represents whether or not creation was successful.
// Any diagnostic information found in the process (including errors) is returned.
func CreateSection(
	ctx context.Context,
	client lucirpc.Client,
	config string,
	sectionType string,
	section string,
	options map[string]json.RawMessage,
) (bool, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result, err := client.CreateSection(
		ctx,
		config,
		sectionType,
		section,
		options,
	)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("problem creating %s.%s section", config, section),
			err.Error(),
		)
		return false, diagnostics
	}

	return result, diagnostics
}

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

// UpdateSection attempts to update an existing section.
// The bool represents whether or not updating was successful.
// Any diagnostic information found in the process (including errors) is returned.
func UpdateSection(
	ctx context.Context,
	client lucirpc.Client,
	config string,
	section string,
	options map[string]json.RawMessage,
) (bool, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result, err := client.UpdateSection(
		ctx,
		config,
		section,
		options,
	)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("problem updating %s.%s section", config, section),
			err.Error(),
		)
		return false, diagnostics
	}

	return result, diagnostics
}
