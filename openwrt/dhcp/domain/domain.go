package domain

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
	hostnameAttribute            = "name"
	hostnameAttributeDescription = "Hostname to assign."
	hostnameUCIOption            = "name"

	ipAddressAttribute            = "ip"
	ipAddressAttributeDescription = "The IP address to be used for this domain."
	ipAddressUCIOption            = "ip"

	schemaDescription = "Binds a domain name to an IP address."

	uciConfig = "dhcp"
	uciType   = "domain"
)

var (
	hostnameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       hostnameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetHostname, hostnameAttribute, hostnameUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetHostname, hostnameAttribute, hostnameUCIOption),
	}

	ipAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ipAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetIPAddress, ipAddressAttribute, ipAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.Any(
				stringvalidator.RegexMatches(
					regexp.MustCompile("^([[:digit:]]{1,3}.){3}[[:digit:]]{1,3}$"),
					`must be a valid IP address (e.g. "192.168.3.1")`,
				),
			),
		},
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		hostnameAttribute:       hostnameSchemaAttribute,
		ipAddressAttribute:      ipAddressSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
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
	Hostname  types.String `tfsdk:"name"`
	Id        types.String `tfsdk:"id"`
	IPAddress types.String `tfsdk:"ip"`
}

func modelGetHostname(m model) types.String  { return m.Hostname }
func modelGetId(m model) types.String        { return m.Id }
func modelGetIPAddress(m model) types.String { return m.IPAddress }

func modelSetHostname(m *model, value types.String)  { m.Hostname = value }
func modelSetId(m *model, value types.String)        { m.Id = value }
func modelSetIPAddress(m *model, value types.String) { m.IPAddress = value }
