package network

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	globalsPacketSteeringAttribute = "packet_steering"
	globalsPacketSteeringUCIOption = "packet_steering"

	globalsSchemaDescription = "Contains interface-independent options affecting the network configuration in general."

	globalsUCIConfig = "network"
	globalsUCIType   = "globals"

	globalsULAPrefixAttribute = "ula_prefix"
	globalsULAPrefixUCIOption = "ula_prefix"
)

var (
	globalsSchemaAttributes = map[string]lucirpcglue.SchemaAttribute[globalsModel, lucirpc.Options, lucirpc.Options]{
		globalsULAPrefixAttribute:      globalsULAPrefixSchemaAttribute,
		globalsPacketSteeringAttribute: globalsPacketSteeringSchemaAttribute,
		lucirpcglue.IdAttribute:        lucirpcglue.IdSchemaAttribute(globalsModelGetId, globalsModelSetId),
	}

	globalsULAPrefixSchemaAttribute = lucirpcglue.StringSchemaAttribute[globalsModel, lucirpc.Options, lucirpc.Options]{
		Description:       "IPv6 ULA prefix for this device.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(globalsModelSetULAPrefix, globalsULAPrefixAttribute, globalsULAPrefixUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(globalsModelGetULAPrefix, globalsULAPrefixAttribute, globalsULAPrefixUCIOption),
	}

	globalsPacketSteeringSchemaAttribute = lucirpcglue.BoolSchemaAttribute[globalsModel, lucirpc.Options, lucirpc.Options]{
		Description:       "Use every CPU to handle packet traffic.",
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(globalsModelSetPacketSteering, globalsPacketSteeringAttribute, globalsPacketSteeringUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(globalsModelGetPacketSteering, globalsPacketSteeringAttribute, globalsPacketSteeringUCIOption),
	}
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

func NewGlobalsResource() resource.Resource {
	return lucirpcglue.NewResource(
		globalsModelGetId,
		globalsSchemaAttributes,
		globalsSchemaDescription,
		globalsUCIConfig,
		globalsUCIType,
	)
}

type globalsModel struct {
	Id             types.String `tfsdk:"id"`
	PacketSteering types.Bool   `tfsdk:"packet_steering"`
	ULAPrefix      types.String `tfsdk:"ula_prefix"`
}

func globalsModelGetId(model globalsModel) types.String           { return model.Id }
func globalsModelGetPacketSteering(model globalsModel) types.Bool { return model.PacketSteering }
func globalsModelGetULAPrefix(model globalsModel) types.String    { return model.ULAPrefix }

func globalsModelSetId(model *globalsModel, value types.String) { model.Id = value }
func globalsModelSetPacketSteering(model *globalsModel, value types.Bool) {
	model.PacketSteering = value
}
func globalsModelSetULAPrefix(model *globalsModel, value types.String) { model.ULAPrefix = value }
