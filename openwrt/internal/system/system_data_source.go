package system

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

func NewSystemDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		systemModelGetId,
		systemSchemaAttributes,
		systemSchemaDescription,
		systemUCIConfig,
		systemUCIType,
	)
}
