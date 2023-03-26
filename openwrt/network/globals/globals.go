package globals

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	packetSteeringAttribute = "packet_steering"
	packetSteeringUCIOption = "packet_steering"

	schemaDescription = "Contains interface-independent options affecting the network configuration in general."

	uciConfig = "network"
	uciType   = "globals"

	ulaPrefixAttribute = "ula_prefix"
	ulaPrefixUCIOption = "ula_prefix"
)

var (
	packetSteeringSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "Use every CPU to handle packet traffic.",
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetPacketSteering, packetSteeringAttribute, packetSteeringUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetPacketSteering, packetSteeringAttribute, packetSteeringUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		ulaPrefixAttribute:      ulaPrefixSchemaAttribute,
		packetSteeringAttribute: packetSteeringSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
	}

	ulaPrefixSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "IPv6 ULA prefix for this device.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetULAPrefix, ulaPrefixAttribute, ulaPrefixUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetULAPrefix, ulaPrefixAttribute, ulaPrefixUCIOption),
	}
)

func NewDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

func NewResource() resource.Resource {
	return lucirpcglue.NewResource(
		modelGetId,
		schemaAttributes,
		schemaDescription,
		uciConfig,
		uciType,
	)
}

type model struct {
	Id             types.String `tfsdk:"id"`
	PacketSteering types.Bool   `tfsdk:"packet_steering"`
	ULAPrefix      types.String `tfsdk:"ula_prefix"`
}

func modelGetId(m model) types.String           { return m.Id }
func modelGetPacketSteering(m model) types.Bool { return m.PacketSteering }
func modelGetULAPrefix(m model) types.String    { return m.ULAPrefix }

func modelSetId(m *model, value types.String)           { m.Id = value }
func modelSetPacketSteering(m *model, value types.Bool) { m.PacketSteering = value }
func modelSetULAPrefix(m *model, value types.String)    { m.ULAPrefix = value }
