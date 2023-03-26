package networkinterface

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
	bringUpOnBootAttribute            = "auto"
	bringUpOnBootAttributeDescription = "Specifies whether to bring up this interface on boot."
	bringUpOnBootUCIOption            = "auto"

	deviceAttribute            = "device"
	deviceAttributeDescription = "Name of the (physical or virtual) device. This name is what the device is known as in LuCI or the `name` field in Terraform. This is not the UCI config name."
	deviceUCIOption            = "device"

	disabledAttribute            = "disabled"
	disabledAttributeDescription = "Disables this interface."
	disabledUCIOption            = "disabled"

	dnsAttribute            = "dns"
	dnsAttributeDescription = "DNS servers"
	dnsUCIOption            = "dns"

	gatewayAttribute            = "gateway"
	gatewayAttributeDescription = "Gateway of the interface"
	gatewayUCIOption            = "gateway"

	ip6AssignAttribute            = "ip6assign"
	ip6AssignAttributeDescription = "Delegate a prefix of given length to this interface"
	ip6AssignUCIOption            = "ip6assign"

	ipAddressAttribute            = "ipaddr"
	ipAddressAttributeDescription = "IP address of the interface"
	ipAddressUCIOption            = "ipaddr"

	macAddressAttribute            = "macaddr"
	macAddressAttributeDescription = "Override the MAC Address of this interface."
	macAddressUCIOption            = "macaddr"

	mtuAttribute            = "mtu"
	mtuAttributeDescription = "Override the default MTU on this interface."
	mtuUCIOption            = "mtu"

	netmaskAttribute            = "netmask"
	netmaskAttributeDescription = "Netmask of the interface"
	netmaskUCIOption            = "netmask"

	peerDNSAttribute            = "peerdns"
	peerDNSAttributeDescription = "Use DHCP-provided DNS servers."
	peerDNSUCIOption            = "peerdns"

	protocolAttribute            = "proto"
	protocolAttributeDescription = `The protocol type of the interface. Currently, only "dhcp, and "static" are supported.`
	protocolDHCP                 = "dhcp"
	protocolDHCPV6               = "dhcpv6"
	protocolStatic               = "static"
	protocolUCIOption            = "proto"

	requestingAddressAttribute            = "reqaddress"
	requestingAddressAttributeDescription = `Behavior for requesting address. Can only be one of "force", "try", or "none".`
	requestingAddressForce                = "force"
	requestingAddressNone                 = "none"
	requestingAddressTry                  = "try"
	requestingAddressUCIOption            = "reqaddress"

	// The fact we can only support `"auto"` is because we haven't figured out how to represent unions.
	// Once we do,
	// we can support `"auto"`, `no`, or 0-64.
	requestingPrefixAttribute            = "reqprefix"
	requestingPrefixAttributeDescription = `Behavior for requesting prefixes. Currently, only "auto" is supported.`
	requestingPrefixAuto                 = "auto"
	requestingPrefixUCIOption            = "reqprefix"

	schemaDescription = "A logic network."

	uciConfig = "network"
	uciType   = "interface"
)

var (
	bringUpOnBootSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       bringUpOnBootAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetBringUpOnBoot, bringUpOnBootAttribute, bringUpOnBootUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetBringUpOnBoot, bringUpOnBootAttribute, bringUpOnBootUCIOption),
	}

	deviceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       deviceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDevice, deviceAttribute, deviceUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDevice, deviceAttribute, deviceUCIOption),
	}

	disabledSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       disabledAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetDisabled, disabledAttribute, disabledUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetDisabled, disabledAttribute, disabledUCIOption),
	}

	dnsSchemaAttribute = lucirpcglue.ListStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       dnsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionListString(modelSetDNS, dnsAttribute, dnsUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionListString(modelGetDNS, dnsAttribute, dnsUCIOption),
		Validators: []validator.List{
			listvalidator.SizeAtLeast(1),
			listvalidator.AnyWithAllWarnings(
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCP,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCPV6,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolStatic,
				),
			),
		},
	}

	gatewaySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       gatewayAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetGateway, gatewayAttribute, gatewayUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetGateway, gatewayAttribute, gatewayUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
				`must be a valid gateway (e.g. "192.168.1.1")`,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolStatic,
			),
		},
	}

	ip6AssignSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ip6AssignAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetIP6Assign, ip6AssignAttribute, ip6AssignUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetIP6Assign, ip6AssignAttribute, ip6AssignUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(0, 64),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolStatic,
			),
		},
	}

	ipAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
				`must be a valid IP address (e.g. "192.168.3.1")`,
			),
		},
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

	netmaskSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       netmaskAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetNetmask, netmaskAttribute, netmaskUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetNetmask, netmaskAttribute, netmaskUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
				`must be a valid netmask (e.g. "255.255.255.0")`,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolStatic,
			),
		},
	}

	peerDNSSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       peerDNSAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetPeerDNS, peerDNSAttribute, peerDNSUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetPeerDNS, peerDNSAttribute, peerDNSUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.AnyBool(
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCP,
				),
				lucirpcglue.RequiresAttributeEqualString(
					path.MatchRoot(protocolAttribute),
					protocolDHCPV6,
				),
			),
		},
	}

	protocolSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       protocolAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetProtocol, protocolAttribute, protocolUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetProtocol, protocolAttribute, protocolUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				protocolDHCP,
				protocolDHCPV6,
				protocolStatic,
			),
		},
	}

	requestingAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       requestingAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetRequestingAddress, requestingAddressAttribute, requestingAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetRequestingAddress, requestingAddressAttribute, requestingAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				requestingAddressForce,
				requestingAddressNone,
				requestingAddressTry,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolDHCPV6,
			),
		},
	}

	requestingPrefixSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       requestingPrefixAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetRequestingPrefix, requestingPrefixAttribute, requestingPrefixUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetRequestingPrefix, requestingPrefixAttribute, requestingPrefixUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				requestingPrefixAuto,
			),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(protocolAttribute),
				protocolDHCPV6,
			),
		},
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		bringUpOnBootAttribute:     bringUpOnBootSchemaAttribute,
		deviceAttribute:            deviceSchemaAttribute,
		disabledAttribute:          disabledSchemaAttribute,
		dnsAttribute:               dnsSchemaAttribute,
		gatewayAttribute:           gatewaySchemaAttribute,
		ip6AssignAttribute:         ip6AssignSchemaAttribute,
		ipAddressAttribute:         ipAddressSchemaAttribute,
		macAddressAttribute:        macAddressSchemaAttribute,
		mtuAttribute:               mtuSchemaAttribute,
		netmaskAttribute:           netmaskSchemaAttribute,
		peerDNSAttribute:           peerDNSSchemaAttribute,
		protocolAttribute:          protocolSchemaAttribute,
		requestingAddressAttribute: requestingAddressSchemaAttribute,
		requestingPrefixAttribute:  requestingPrefixSchemaAttribute,
		lucirpcglue.IdAttribute:    lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
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

func modelGetBringUpOnBoot(m model) types.Bool       { return m.BringUpOnBoot }
func modelGetDevice(m model) types.String            { return m.Device }
func modelGetDisabled(m model) types.Bool            { return m.Disabled }
func modelGetDNS(m model) types.List                 { return m.DNS }
func modelGetGateway(m model) types.String           { return m.Gateway }
func modelGetId(m model) types.String                { return m.Id }
func modelGetIP6Assign(m model) types.Int64          { return m.IP6Assign }
func modelGetIPAddress(m model) types.String         { return m.IPAddress }
func modelGetMacAddress(m model) types.String        { return m.MacAddress }
func modelGetMTU(m model) types.Int64                { return m.MTU }
func modelGetNetmask(m model) types.String           { return m.Netmask }
func modelGetPeerDNS(m model) types.Bool             { return m.PeerDNS }
func modelGetProtocol(m model) types.String          { return m.Protocol }
func modelGetRequestingAddress(m model) types.String { return m.RequestingAddress }
func modelGetRequestingPrefix(m model) types.String  { return m.RequestingPrefix }

func modelSetBringUpOnBoot(m *model, value types.Bool)       { m.BringUpOnBoot = value }
func modelSetDevice(m *model, value types.String)            { m.Device = value }
func modelSetDisabled(m *model, value types.Bool)            { m.Disabled = value }
func modelSetDNS(m *model, value types.List)                 { m.DNS = value }
func modelSetGateway(m *model, value types.String)           { m.Gateway = value }
func modelSetId(m *model, value types.String)                { m.Id = value }
func modelSetIP6Assign(m *model, value types.Int64)          { m.IP6Assign = value }
func modelSetIPAddress(m *model, value types.String)         { m.IPAddress = value }
func modelSetMacAddress(m *model, value types.String)        { m.MacAddress = value }
func modelSetMTU(m *model, value types.Int64)                { m.MTU = value }
func modelSetNetmask(m *model, value types.String)           { m.Netmask = value }
func modelSetPeerDNS(m *model, value types.Bool)             { m.PeerDNS = value }
func modelSetProtocol(m *model, value types.String)          { m.Protocol = value }
func modelSetRequestingAddress(m *model, value types.String) { m.RequestingAddress = value }
func modelSetRequestingPrefix(m *model, value types.String)  { m.RequestingPrefix = value }
