package system

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

func NewSystemResource() resource.Resource {
	return lucirpcglue.NewResource(
		systemModelGetId,
		systemSchemaAttributes,
		systemSchemaDescription,
		systemUCIConfig,
		systemUCIType,
	)
}
