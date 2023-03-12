package network

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

func NewDeviceResource() resource.Resource {
	return lucirpcglue.NewResource(
		deviceModelGetId,
		deviceSchemaAttributes,
		deviceSchemaDescription,
		deviceUCIConfig,
		deviceUCIType,
	)
}
