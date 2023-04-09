package host

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	addDNSEntriesAttribute            = "dns"
	addDNSEntriesAttributeDescription = "Add static forward and reverse DNS entries for this host."
	addDNSEntriesUCIOption            = "dns"

	hostnameAttribute            = "name"
	hostnameAttributeDescription = "Hostname to assign."
	hostnameUCIOption            = "name"

	ipAddressAttribute            = "ip"
	ipAddressAttributeDescription = "The IP address to be used for this host, or `ignore` to ignore any DHCP request from this host."
	ipAddressUCIOption            = "ip"

	macAddressAttribute            = "mac"
	macAddressAttributeDescription = "The hardware address(es) of this host, separated by spaces."
	macAddressUCIOption            = "mac"

	schemaDescription = "Assign a fixed IP address to hosts."

	uciConfig = "dhcp"
	uciType   = "host"
)

var (
	addDNSEntriesSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       addDNSEntriesAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetAddDNSEntries, addDNSEntriesAttribute, addDNSEntriesUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetAddDNSEntries, addDNSEntriesAttribute, addDNSEntriesUCIOption),
	}

	hostnameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       hostnameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetHostname, hostnameAttribute, hostnameUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetHostname, hostnameAttribute, hostnameUCIOption),
	}

	ipAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.Any(
				stringvalidator.OneOf(
					"ignore",
				),
				stringvalidator.RegexMatches(
					regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
					`must be a valid IP address (e.g. "192.168.3.1")`,
				),
			),
		},
	}

	macAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       macAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetMACAddress, macAddressAttribute, macAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetMACAddress, macAddressAttribute, macAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:xdigit:]][[:xdigit:]]:){5}[[:xdigit:]][[:xdigit:]]$"),
				`must be a valid MAC address (e.g. "12:34:56:78:90:ab")`,
			),
		},
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		addDNSEntriesAttribute:  addDNSEntriesSchemaAttribute,
		hostnameAttribute:       hostnameSchemaAttribute,
		ipAddressAttribute:      ipAddressSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		macAddressAttribute:     macAddressSchemaAttribute,
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
	AddDNSEntries types.Bool   `tfsdk:"dns"`
	Hostname      types.String `tfsdk:"name"`
	Id            types.String `tfsdk:"id"`
	IPAddress     types.String `tfsdk:"ip"`
	MACAddress    types.String `tfsdk:"mac"`
}

func modelGetAddDNSEntries(m model) types.Bool { return m.AddDNSEntries }
func modelGetHostname(m model) types.String    { return m.Hostname }
func modelGetId(m model) types.String          { return m.Id }
func modelGetIPAddress(m model) types.String   { return m.IPAddress }
func modelGetMACAddress(m model) types.String  { return m.MACAddress }

func modelSetAddDNSEntries(m *model, value types.Bool) { m.AddDNSEntries = value }
func modelSetHostname(m *model, value types.String)    { m.Hostname = value }
func modelSetId(m *model, value types.String)          { m.Id = value }
func modelSetIPAddress(m *model, value types.String)   { m.IPAddress = value }
func modelSetMACAddress(m *model, value types.String)  { m.MACAddress = value }
