package network

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

func NewDeviceDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		deviceModelGetId,
		deviceSchemaAttributes,
		deviceSchemaDescription,
		deviceUCIConfig,
		deviceUCIType,
	)
}
