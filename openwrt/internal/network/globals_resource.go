package network

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

func NewGlobalsResource() resource.Resource {
	return lucirpcglue.NewResource(
		globalsModelGetId,
		globalsSchemaAttributes,
		globalsSchemaDescription,
		globalsUCIConfig,
		globalsUCIGlobalsType,
	)
}
