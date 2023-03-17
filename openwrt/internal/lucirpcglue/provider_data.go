package lucirpcglue

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
)

type ConfigureRequest struct {
	ProviderData any
}

func NewProviderData(
	client lucirpc.Client,
	typeName string,
) ProviderData {
	return ProviderData{
		Client:   client,
		TypeName: typeName,
	}
}

// ParseProviderData attempts to extract a [ProviderData] from the given [ConfigureRequest].
// Any diagnostic information found in the process (including errors) is returned.
func ParseProviderData(
	req ConfigureRequest,
) (ProviderData, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}
	providerData, ok := req.ProviderData.(ProviderData)
	if !ok {
		diagnostics.AddError(
			"OpenWrt provider not configured correctly",
			"Expected the provider data to be of a given structure, but it was not. This is a problem with the provider implementation. Please report this to https://github.com/joneshf/terraform-provider-openwrt",
		)
		return ProviderData{}, diagnostics
	}

	return providerData, diagnostics
}

type ProviderData struct {
	Client   lucirpc.Client
	TypeName string
}
