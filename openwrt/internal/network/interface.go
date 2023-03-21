package network

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	interfaceBringUpOnBootAttribute            = "auto"
	interfaceBringUpOnBootAttributeDescription = "Specifies whether to bring up this interface on boot."
	interfaceBringUpOnBootUCIOption            = "auto"

	interfaceDeviceAttribute            = "device"
	interfaceDeviceAttributeDescription = "Name of the (physical or virtual) device. This name is what the device is known as in LuCI or the `name` field in Terraform. This is not the UCI config name."
	interfaceDeviceUCIOption            = "device"

	interfaceDisabledAttribute            = "disabled"
	interfaceDisabledAttributeDescription = "Disables this interface."
	interfaceDisabledUCIOption            = "disabled"

	interfaceDNSAttribute            = "dns"
	interfaceDNSAttributeDescription = "DNS servers"
	interfaceDNSUCIOption            = "dns"

	interfaceGatewayAttribute            = "gateway"
	interfaceGatewayAttributeDescription = "Gateway of the interface"
	interfaceGatewayUCIOption            = "gateway"

	interfaceIP6AssignAttribute            = "ip6assign"
	interfaceIP6AssignAttributeDescription = "Delegate a prefix of given length to this interface"
	interfaceIP6AssignUCIOption            = "ip6assign"

	interfaceIPAddressAttribute            = "ipaddr"
	interfaceIPAddressAttributeDescription = "IP address of the interface"
	interfaceIPAddressUCIOption            = "ipaddr"

	interfaceMacAddressAttribute            = "macaddr"
	interfaceMacAddressAttributeDescription = "Override the MAC Address of this interface."
	interfaceMacAddressUCIOption            = "macaddr"

	interfaceMTUAttribute            = "mtu"
	interfaceMTUAttributeDescription = "Override the default MTU on this interface."
	interfaceMTUUCIOption            = "mtu"

	interfaceNetmaskAttribute            = "netmask"
	interfaceNetmaskAttributeDescription = "Netmask of the interface"
	interfaceNetmaskUCIOption            = "netmask"

	interfacePeerDNSAttribute            = "peerdns"
	interfacePeerDNSAttributeDescription = "Use DHCP-provided DNS servers."
	interfacePeerDNSUCIOption            = "peerdns"

	interfaceProtocolAttribute            = "proto"
	interfaceProtocolAttributeDescription = `The protocol type of the interface. Currently, only "dhcp, and "static" are supported.`
	interfaceProtocolDHCP                 = "dhcp"
	interfaceProtocolDHCPV6               = "dhcpv6"
	interfaceProtocolStatic               = "static"
	interfaceProtocolUCIOption            = "proto"

	interfaceRequestingAddressAttribute            = "reqaddress"
	interfaceRequestingAddressAttributeDescription = `Behavior for requesting address. Can only be one of "force", "try", or "none".`
	interfaceRequestingAddressForce                = "force"
	interfaceRequestingAddressNone                 = "none"
	interfaceRequestingAddressTry                  = "try"
	interfaceRequestingAddressUCIOption            = "reqaddress"

	// The fact we can only support `"auto"` is because we haven't figured out how to represent unions.
	// Once we do,
	// we can support `"auto"`, `no`, or 0-64.
	interfaceRequestingPrefixAttribute            = "reqprefix"
	interfaceRequestingPrefixAttributeDescription = `Behavior for requesting prefixes. Currently, only "auto" is supported.`
	interfaceRequestingPrefixAuto                 = "auto"
	interfaceRequestingPrefixUCIOption            = "reqprefix"

	interfaceSchemaDescription = "A logic network."

	interfaceUCIConfig = "network"
	interfaceUCIType   = "interface"
)

var (
	interfaceBringUpOnBootSchemaAttribute = lucirpcglue.BoolSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceBringUpOnBootAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(interfaceModelSetBringUpOnBoot, interfaceBringUpOnBootAttribute, interfaceBringUpOnBootUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(interfaceModelGetBringUpOnBoot, interfaceBringUpOnBootAttribute, interfaceBringUpOnBootUCIOption),
	}

	interfaceDeviceSchemaAttribute = lucirpcglue.StringSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceDeviceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(interfaceModelSetDevice, interfaceDeviceAttribute, interfaceDeviceUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(interfaceModelGetDevice, interfaceDeviceAttribute, interfaceDeviceUCIOption),
	}

	interfaceDisabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceDisabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(interfaceModelSetDisabled, interfaceDisabledAttribute, interfaceDisabledUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(interfaceModelGetDisabled, interfaceDisabledAttribute, interfaceDisabledUCIOption),
	}

	interfaceDNSSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceDNSAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(interfaceModelSetDNS, interfaceDNSAttribute, interfaceDNSUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(interfaceModelGetDNS, interfaceDNSAttribute, interfaceDNSUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.AnyWithAllWarnings(
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(interfaceProtocolAttribute),
					interfaceProtocolDHCP,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(interfaceProtocolAttribute),
					interfaceProtocolDHCPV6,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(interfaceProtocolAttribute),
					interfaceProtocolStatic,
				),
			),
		},
	}

	interfaceGatewaySchemaAttribute = lucirpcglue.StringSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceGatewayAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(interfaceModelSetGateway, interfaceGatewayAttribute, interfaceGatewayUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(interfaceModelGetGateway, interfaceGatewayAttribute, interfaceGatewayUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
				`must be a valid gateway (e.g. "192.168.1.1")`,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(interfaceProtocolAttribute),
				interfaceProtocolStatic,
			),
		},
	}

	interfaceIP6AssignSchemaAttribute = lucirpcglue.Int64SchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceIP6AssignAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(interfaceModelSetIP6Assign, interfaceIP6AssignAttribute, interfaceIP6AssignUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(interfaceModelGetIP6Assign, interfaceIP6AssignAttribute, interfaceIP6AssignUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(0, 64),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(interfaceProtocolAttribute),
				interfaceProtocolStatic,
			),
		},
	}

	interfaceIPAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceIPAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(interfaceModelSetIPAddress, interfaceIPAddressAttribute, interfaceIPAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(interfaceModelGetIPAddress, interfaceIPAddressAttribute, interfaceIPAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
				`must be a valid IP address (e.g. "192.168.3.1")`,
			),
		},
	}

	interfaceMacAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceMacAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(interfaceModelSetMacAddress, interfaceMacAddressAttribute, interfaceMacAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(interfaceModelGetMacAddress, interfaceMacAddressAttribute, interfaceMacAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:xdigit:]][[:xdigit:]]:){5}[[:xdigit:]][[:xdigit:]]$"),
				`must be a valid MAC address (e.g. "12:34:56:78:90:ab")`,
			),
		},
	}

	interfaceMTUSchemaAttribute = lucirpcglue.Int64SchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceMTUAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(interfaceModelSetMTU, interfaceMTUAttribute, interfaceMTUUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(interfaceModelGetMTU, interfaceMTUAttribute, interfaceMTUUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(576, 9200),
		},
	}

	interfaceNetmaskSchemaAttribute = lucirpcglue.StringSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceNetmaskAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(interfaceModelSetNetmask, interfaceNetmaskAttribute, interfaceNetmaskUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(interfaceModelGetNetmask, interfaceNetmaskAttribute, interfaceNetmaskUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
				`must be a valid netmask (e.g. "255.255.255.0")`,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(interfaceProtocolAttribute),
				interfaceProtocolStatic,
			),
		},
	}

	interfacePeerDNSSchemaAttribute = lucirpcglue.BoolSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfacePeerDNSAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(interfaceModelSetPeerDNS, interfacePeerDNSAttribute, interfacePeerDNSUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(interfaceModelGetPeerDNS, interfacePeerDNSAttribute, interfacePeerDNSUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(interfaceProtocolAttribute),
				interfaceProtocolDHCP,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(interfaceProtocolAttribute),
				interfaceProtocolDHCPV6,
			),
		},
	}

	interfaceProtocolSchemaAttribute = lucirpcglue.StringSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceProtocolAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(interfaceModelSetProtocol, interfaceProtocolAttribute, interfaceProtocolUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(interfaceModelGetProtocol, interfaceProtocolAttribute, interfaceProtocolUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				interfaceProtocolDHCP,
				interfaceProtocolDHCPV6,
				interfaceProtocolStatic,
			),
		},
	}

	interfaceRequestingAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceRequestingAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(interfaceModelSetRequestingAddress, interfaceRequestingAddressAttribute, interfaceRequestingAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(interfaceModelGetRequestingAddress, interfaceRequestingAddressAttribute, interfaceRequestingAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				interfaceRequestingAddressForce,
				interfaceRequestingAddressNone,
				interfaceRequestingAddressTry,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(interfaceProtocolAttribute),
				interfaceProtocolDHCPV6,
			),
		},
	}

	interfaceRequestingPrefixSchemaAttribute = lucirpcglue.StringSchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceRequestingPrefixAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(interfaceModelSetRequestingPrefix, interfaceRequestingPrefixAttribute, interfaceRequestingPrefixUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(interfaceModelGetRequestingPrefix, interfaceRequestingPrefixAttribute, interfaceRequestingPrefixUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				interfaceRequestingPrefixAuto,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(interfaceProtocolAttribute),
				interfaceProtocolDHCPV6,
			),
		},
	}

	interfaceSchemaAttributes = map[string]lucirpcglue.SchemaAttribute[interfaceModel, lucirpc.Options, lucirpc.Options]{
		interfaceBringUpOnBootAttribute:     interfaceBringUpOnBootSchemaAttribute,
		interfaceDeviceAttribute:            interfaceDeviceSchemaAttribute,
		interfaceDisabledAttribute:          interfaceDisabledSchemaAttribute,
		interfaceDNSAttribute:               interfaceDNSSchemaAttribute,
		interfaceGatewayAttribute:           interfaceGatewaySchemaAttribute,
		interfaceIP6AssignAttribute:         interfaceIP6AssignSchemaAttribute,
		interfaceIPAddressAttribute:         interfaceIPAddressSchemaAttribute,
		interfaceMacAddressAttribute:        interfaceMacAddressSchemaAttribute,
		interfaceMTUAttribute:               interfaceMTUSchemaAttribute,
		interfaceNetmaskAttribute:           interfaceNetmaskSchemaAttribute,
		interfacePeerDNSAttribute:           interfacePeerDNSSchemaAttribute,
		interfaceProtocolAttribute:          interfaceProtocolSchemaAttribute,
		interfaceRequestingAddressAttribute: interfaceRequestingAddressSchemaAttribute,
		interfaceRequestingPrefixAttribute:  interfaceRequestingPrefixSchemaAttribute,
		lucirpcglue.IdAttribute:             lucirpcglue.IdSchemaAttribute(interfaceModelGetId, interfaceModelSetId),
	}
)

func NewInterfaceDataSource() datasource.DataSource {
	return lucirpcglue.NewDataSource(
		interfaceModelGetId,
		interfaceSchemaAttributes,
		interfaceSchemaDescription,
		interfaceUCIConfig,
		interfaceUCIType,
	)
}

func NewInterfaceResource() resource.Resource {
	return lucirpcglue.NewResource(
		interfaceModelGetId,
		interfaceSchemaAttributes,
		interfaceSchemaDescription,
		interfaceUCIConfig,
		interfaceUCIType,
	)
}

type interfaceModel struct {
	BringUpOnBoot     types.Bool   `tfsdk:"auto"`
	Device            types.String `tfsdk:"device"`
	Disabled          types.Bool   `tfsdk:"disabled"`
	DNS               types.List   `tfsdk:"dns"`
	Gateway           types.String `tfsdk:"gateway"`
	Id                types.String `tfsdk:"id"`
	IP6Assign         types.Int64  `tfsdk:"ip6assign"`
	IPAddress         types.String `tfsdk:"ipaddr"`
	MacAddress        types.String `tfsdk:"macaddr"`
	MTU               types.Int64  `tfsdk:"mtu"`
	Netmask           types.String `tfsdk:"netmask"`
	PeerDNS           types.Bool   `tfsdk:"peerdns"`
	Protocol          types.String `tfsdk:"proto"`
	RequestingAddress types.String `tfsdk:"reqaddress"`
	RequestingPrefix  types.String `tfsdk:"reqprefix"`
}

func interfaceModelGetBringUpOnBoot(model interfaceModel) types.Bool { return model.BringUpOnBoot }
func interfaceModelGetDevice(model interfaceModel) types.String      { return model.Device }
func interfaceModelGetDisabled(model interfaceModel) types.Bool      { return model.Disabled }
func interfaceModelGetDNS(model interfaceModel) types.List           { return model.DNS }
func interfaceModelGetGateway(model interfaceModel) types.String     { return model.Gateway }
func interfaceModelGetId(model interfaceModel) types.String          { return model.Id }
func interfaceModelGetIP6Assign(model interfaceModel) types.Int64    { return model.IP6Assign }
func interfaceModelGetIPAddress(model interfaceModel) types.String   { return model.IPAddress }
func interfaceModelGetMacAddress(model interfaceModel) types.String  { return model.MacAddress }
func interfaceModelGetMTU(model interfaceModel) types.Int64          { return model.MTU }
func interfaceModelGetNetmask(model interfaceModel) types.String     { return model.Netmask }
func interfaceModelGetPeerDNS(model interfaceModel) types.Bool       { return model.PeerDNS }
func interfaceModelGetProtocol(model interfaceModel) types.String    { return model.Protocol }
func interfaceModelGetRequestingAddress(model interfaceModel) types.String {
	return model.RequestingAddress
}
func interfaceModelGetRequestingPrefix(model interfaceModel) types.String {
	return model.RequestingPrefix
}

func interfaceModelSetBringUpOnBoot(model *interfaceModel, value types.Bool) {
	model.BringUpOnBoot = value
}
func interfaceModelSetDevice(model *interfaceModel, value types.String)     { model.Device = value }
func interfaceModelSetDisabled(model *interfaceModel, value types.Bool)     { model.Disabled = value }
func interfaceModelSetDNS(model *interfaceModel, value types.List)          { model.DNS = value }
func interfaceModelSetGateway(model *interfaceModel, value types.String)    { model.Gateway = value }
func interfaceModelSetId(model *interfaceModel, value types.String)         { model.Id = value }
func interfaceModelSetIP6Assign(model *interfaceModel, value types.Int64)   { model.IP6Assign = value }
func interfaceModelSetIPAddress(model *interfaceModel, value types.String)  { model.IPAddress = value }
func interfaceModelSetMacAddress(model *interfaceModel, value types.String) { model.MacAddress = value }
func interfaceModelSetMTU(model *interfaceModel, value types.Int64)         { model.MTU = value }
func interfaceModelSetNetmask(model *interfaceModel, value types.String)    { model.Netmask = value }
func interfaceModelSetPeerDNS(model *interfaceModel, value types.Bool)      { model.PeerDNS = value }
func interfaceModelSetProtocol(model *interfaceModel, value types.String)   { model.Protocol = value }
func interfaceModelSetRequestingAddress(model *interfaceModel, value types.String) {
	model.RequestingAddress = value
}
func interfaceModelSetRequestingPrefix(model *interfaceModel, value types.String) {
	model.RequestingPrefix = value
}
