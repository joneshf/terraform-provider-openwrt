package lucirpcglue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

// CreateSection attempts to create a new section.
// Any diagnostic information found in the process (including errors) is returned.
func CreateSection(
	ctx context.Context,
	client lucirpc.Client,
	config string,
	sectionType string,
	section string,
	options lucirpc.Options,
) diag.Diagnostics {
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
		return diagnostics
	}

	if !result {
		diagnostics.AddError(
			fmt.Sprintf("Could not create %s.%s section", config, section),
			"It is not currently known why this happens. It is unclear if this is a problem with the provider. Please double check the values provided are acceptable.",
		)
		return diagnostics
	}

	return diagnostics
}

// DeleteSection attempts to delete an existing section.
// Any diagnostic information found in the process (including errors) is returned.
func DeleteSection(
	ctx context.Context,
	client lucirpc.Client,
	config string,
	section string,
) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}
	result, err := client.DeleteSection(
		ctx,
		config,
		section,
	)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("problem deleting %s.%s section", config, section),
			err.Error(),
		)
		return diagnostics
	}

	if !result {
		diagnostics.AddError(
			fmt.Sprintf("Could not delete %s.%s section", config, section),
			"It is not currently known why this happens. It is unclear if this is a problem with the provider. Please double check the values provided are acceptable.",
		)
		return diagnostics
	}

	return diagnostics
}

// GetMetadataString attempts to parse the given metadata key from the section.
// Any diagnostic information found in the process (including errors) is returned.
func GetSection(
	ctx context.Context,
	client lucirpc.Client,
	config string,
	section string,
) (lucirpc.Options, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	result, err := client.GetSection(ctx, config, section)
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("problem getting %s.%s section", config, section),
			err.Error(),
		)
		return lucirpc.Options{}, diagnostics
	}

	return result, diagnostics
}

// UpdateSection attempts to update an existing section.
// Any diagnostic information found in the process (including errors) is returned.
func UpdateSection(
	ctx context.Context,
	client lucirpc.Client,
	config string,
	section string,
	options lucirpc.Options,
) diag.Diagnostics {
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
		return diagnostics
	}

	if !result {
		diagnostics.AddError(
			fmt.Sprintf("Could not update %s.%s section", config, section),
			"It is not currently known why this happens. It is unclear if this is a problem with the provider. Please double check the values provided are acceptable.",
		)
		return diagnostics
	}

	return diagnostics
}
