package device

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	bridgePortsAttribute            = "ports"
	bridgePortsAttributeDescription = "Specifies the wired ports to attach to this bridge."
	bridgePortsUCIOption            = "ports"

	bringUpEmptyBridgeAttribute            = "bridge_empty"
	bringUpEmptyBridgeAttributeDescription = "Bring up the bridge device even if no ports are attached"
	bringUpEmptyBridgeUCIOption            = "bridge_empty"

	dadTransmitsAttribute            = "dadtransmits"
	dadTransmitsAttributeDescription = "Amount of Duplicate Address Detection probes to send"
	dadTransmitsUCIOption            = "dadtransmits"

	enableIPv6Attribute            = "ipv6"
	enableIPv6AttributeDescription = "Enable IPv6 for the device."
	enableIPv6UCIOption            = "ipv6"

	macAddressAttribute            = "macaddr"
	macAddressAttributeDescription = "MAC Address of the device."
	macAddressUCIOption            = "macaddr"

	mtuAttribute            = "mtu"
	mtuAttributeDescription = "Maximum Transmissible Unit."
	mtuUCIOption            = "mtu"

	mtu6Attribute            = "mtu6"
	mtu6AttributeDescription = "Maximum Transmissible Unit for IPv6."
	mtu6UCIOption            = "mtu6"

	nameAttribute            = "name"
	nameAttributeDescription = "Name of the device. This name is referenced in other network configuration."
	nameUCIOption            = "name"

	schemaDescription = `A physical or virtual "device" in OpenWrt jargon. Commonly referred to as an "interface" in other networking jargon.`

	txQueueLengthAttribute            = "txqueuelen"
	txQueueLengthAttributeDescription = "Transmission queue length."
	txQueueLengthUCIOption            = "txqueuelen"

	typeAttribute            = "type"
	typeAttributeDescription = `The type of device. Currently, only "bridge" is supported.`
	typeBridge               = "bridge"
	typeUCIOption            = "type"

	uciConfig = "network"
	uciType   = "device"
)

var (
	bridgePortsSchemaAttribute = lucirpcglue.SetStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       bridgePortsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionSetString(modelSetBridgePorts, bridgePortsAttribute, bridgePortsUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionSetString(modelGetBridgePorts, bridgePortsAttribute, bridgePortsUCIOption),
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(typeAttribute),
				typeBridge,
			),
		},
	}

	bringUpEmptyBridgeSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       bringUpEmptyBridgeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetBringUpEmptyBridge, bringUpEmptyBridgeAttribute, bringUpEmptyBridgeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetBringUpEmptyBridge, bringUpEmptyBridgeAttribute, bringUpEmptyBridgeUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(typeAttribute),
				typeBridge,
			),
		},
	}

	dadTransmitsSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       dadTransmitsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetDADTransmits, dadTransmitsAttribute, dadTransmitsUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetDADTransmits, dadTransmitsAttribute, dadTransmitsUCIOption),
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
			lucirpcglue.RequiresAttributeEqualBool(
				path.MatchRoot(enableIPv6Attribute),
				true,
			),
		},
	}

	enableIPv6SchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enableIPv6AttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnableIPv6, enableIPv6Attribute, enableIPv6UCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnableIPv6, enableIPv6Attribute, enableIPv6UCIOption),
	}

	macAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       macAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetMacAddress, macAddressAttribute, macAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetMacAddress, macAddressAttribute, macAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:xdigit:]][[:xdigit:]]:){5}[[:xdigit:]][[:xdigit:]]$"),
				`must be a valid MAC address (e.g. "12:34:56:78:90:ab")`,
			),
		},
	}

	mtuSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mtuAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMTU, mtuAttribute, mtuUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMTU, mtuAttribute, mtuUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(576, 9200),
		},
	}

	mtu6SchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mtu6AttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMTU6, mtu6Attribute, mtu6UCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMTU6, mtu6Attribute, mtu6UCIOption),
		Validators: []validator.Int64{
			int64validator.Between(576, 9200),
			lucirpcglue.RequiresAttributeEqualBool(
				path.MatchRoot(enableIPv6Attribute),
				true,
			),
		},
	}

	nameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       nameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetName, nameAttribute, nameUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetName, nameAttribute, nameUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		bridgePortsAttribute:        bridgePortsSchemaAttribute,
		bringUpEmptyBridgeAttribute: bringUpEmptyBridgeSchemaAttribute,
		dadTransmitsAttribute:       dadTransmitsSchemaAttribute,
		enableIPv6Attribute:         enableIPv6SchemaAttribute,
		macAddressAttribute:         macAddressSchemaAttribute,
		mtuAttribute:                mtuSchemaAttribute,
		mtu6Attribute:               mtu6SchemaAttribute,
		nameAttribute:               nameSchemaAttribute,
		txQueueLengthAttribute:      txQueueLengthSchemaAttribute,
		typeAttribute:               typeSchemaAttribute,
		lucirpcglue.IdAttribute:     lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
	}

	txQueueLengthSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       txQueueLengthAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetTXQueueLength, txQueueLengthAttribute, txQueueLengthUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetTXQueueLength, txQueueLengthAttribute, txQueueLengthUCIOption),
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
		},
	}

	typeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       typeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetType, typeAttribute, typeUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetType, typeAttribute, typeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				typeBridge,
			),
		},
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
	BridgePorts        types.Set    `tfsdk:"ports"`
	BringUpEmptyBridge types.Bool   `tfsdk:"bridge_empty"`
	DADTransmits       types.Int64  `tfsdk:"dadtransmits"`
	EnableIPv6         types.Bool   `tfsdk:"ipv6"`
	Id                 types.String `tfsdk:"id"`
	MacAddress         types.String `tfsdk:"macaddr"`
	MTU                types.Int64  `tfsdk:"mtu"`
	MTU6               types.Int64  `tfsdk:"mtu6"`
	Name               types.String `tfsdk:"name"`
	TXQueueLength      types.Int64  `tfsdk:"txqueuelen"`
	Type               types.String `tfsdk:"type"`
}

func modelGetBridgePorts(m model) types.Set         { return m.BridgePorts }
func modelGetBringUpEmptyBridge(m model) types.Bool { return m.BringUpEmptyBridge }
func modelGetDADTransmits(m model) types.Int64      { return m.DADTransmits }
func modelGetEnableIPv6(m model) types.Bool         { return m.EnableIPv6 }
func modelGetId(m model) types.String               { return m.Id }
func modelGetMacAddress(m model) types.String       { return m.MacAddress }
func modelGetMTU(m model) types.Int64               { return m.MTU }
func modelGetMTU6(m model) types.Int64              { return m.MTU6 }
func modelGetName(m model) types.String             { return m.Name }
func modelGetTXQueueLength(m model) types.Int64     { return m.TXQueueLength }
func modelGetType(m model) types.String             { return m.Type }

func modelSetBridgePorts(m *model, value types.Set)         { m.BridgePorts = value }
func modelSetBringUpEmptyBridge(m *model, value types.Bool) { m.BringUpEmptyBridge = value }
func modelSetDADTransmits(m *model, value types.Int64)      { m.DADTransmits = value }
func modelSetEnableIPv6(m *model, value types.Bool)         { m.EnableIPv6 = value }
func modelSetId(m *model, value types.String)               { m.Id = value }
func modelSetMacAddress(m *model, value types.String)       { m.MacAddress = value }
func modelSetMTU(m *model, value types.Int64)               { m.MTU = value }
func modelSetMTU6(m *model, value types.Int64)              { m.MTU6 = value }
func modelSetName(m *model, value types.String)             { m.Name = value }
func modelSetTXQueueLength(m *model, value types.Int64)     { m.TXQueueLength = value }
func modelSetType(m *model, value types.String)             { m.Type = value }
