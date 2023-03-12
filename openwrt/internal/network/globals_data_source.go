package network

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

func NewGlobalsDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		globalsModelGetId,
		globalsSchemaAttributes,
		globalsSchemaDescription,
		globalsUCIConfig,
		globalsUCIType,
	)
}
