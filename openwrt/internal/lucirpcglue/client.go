package lucirpcglue

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

// NewClient attempts to construct a new [lucirpc.Client].
// Any diagnostic information found in the process (including errors) is returned.
func NewClient(
	req datasource.ConfigureRequest,
) (*lucirpc.Client, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	client, ok := req.ProviderData.(*lucirpc.Client)
	if !ok {
		diagnostics.AddError(
			"OpenWrt provider not configured correctly",
			"Expected UCI tree, but one was not provided. This is a problem with the provider implementation. Please report this to https://github.com/joneshf/terraform-provider-openwrt",
		)
		return nil, diagnostics
	}

	return client, diagnostics
}
