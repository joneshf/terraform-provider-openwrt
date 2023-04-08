package dnsmasq

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	authoritativeModeAttribute            = "authoritative"
	authoritativeModeAttributeDescription = "Force dnsmasq into authoritative mode. This speeds up DHCP leasing. Used if this is the only server on the network."
	authoritativeModeUCIOption            = "authoritative"

	domainAttribute            = "domain"
	domainAttributeDescription = "DNS domain handed out to DHCP clients."
	domainUCIOption            = "domain"

	domainNeededAttribute            = "domainneeded"
	domainNeededAttributeDescription = "Never forward queries for plain names, without dots or domain parts, to upstream nameservers."
	domainNeededUCIOption            = "domainneeded"

	ednsPacketMaxAttribute            = "ednspacket_max"
	ednsPacketMaxAttributeDescription = "Specify the largest EDNS.0 UDP packet which is supported by the DNS forwarder."
	ednsPacketMaxUCIOption            = "ednspacket_max"

	expandHostsAttribute            = "expandhosts"
	expandHostsAttributeDescription = "Never forward queries for plain names, without dots or domain parts, to upstream nameservers."
	expandHostsUCIOption            = "expandhosts"

	leaseFileAttribute            = "leasefile"
	leaseFileAttributeDescription = "Store DHCP leases in this file."
	leaseFileUCIOption            = "leasefile"

	localizeQueriesAttribute            = "localise_queries"
	localizeQueriesAttributeDescription = "Choose IP address to match the incoming interface if multiple addresses are assigned to a host name in `/etc/hosts`."
	localizeQueriesUCIOption            = "localise_queries"

	localLookupAttribute            = "local"
	localLookupAttributeDescription = "Look up DNS entries for this domain from `/etc/hosts`."
	localLookupUCIOption            = "local"

	localServiceAttribute            = "localservice"
	localServiceAttributeDescription = "Accept DNS queries only from hosts whose address is on a local subnet."
	localServiceUCIOption            = "localservice"

	readEthersAttribute            = "readethers"
	readEthersAttributeDescription = "Read static lease entries from `/etc/ethers`, re-read on SIGHUP."
	readEthersUCIOption            = "readethers"

	rebindLocalhostAttribute            = "rebind_localhost"
	rebindLocalhostAttributeDescription = "Allows upstream 127.0.0.0/8 responses, required for DNS based blocklist services. Only takes effect if rebind protection is enabled."
	rebindLocalhostUCIOption            = "rebind_localhost"

	rebindProtectionAttribute            = "rebind_protection"
	rebindProtectionAttributeDescription = "Enables DNS rebind attack protection by discarding upstream RFC1918 responses."
	rebindProtectionUCIOption            = "rebind_protection"

	resolvFileAttribute            = "resolvfile"
	resolvFileAttributeDescription = "Specifies an alternative resolv file."
	resolvFileUCIOption            = "resolvfile"

	schemaDescription = "A lightweight DHCP and caching DNS server."

	uciConfig = "dhcp"
	uciType   = "dnsmasq"
)

var (
	authoritativeModeSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       authoritativeModeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetAuthoritativeMode, authoritativeModeAttribute, authoritativeModeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetAuthoritativeMode, authoritativeModeAttribute, authoritativeModeUCIOption),
	}

	domainSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       domainAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDomain, domainAttribute, domainUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDomain, domainAttribute, domainUCIOption),
	}

	domainNeededSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       domainNeededAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetDomainNeeded, domainNeededAttribute, domainNeededUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetDomainNeeded, domainNeededAttribute, domainNeededUCIOption),
	}

	ednsPacketMaxSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ednsPacketMaxAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetEDNSPacketMax, ednsPacketMaxAttribute, ednsPacketMaxUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetEDNSPacketMax, ednsPacketMaxAttribute, ednsPacketMaxUCIOption),
	}

	expandHostsSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       expandHostsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetExpandHosts, expandHostsAttribute, expandHostsUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetExpandHosts, expandHostsAttribute, expandHostsUCIOption),
	}

	leaseFileSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       leaseFileAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLeaseFile, leaseFileAttribute, leaseFileUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLeaseFile, leaseFileAttribute, leaseFileUCIOption),
	}

	localizeQueriesSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       localizeQueriesAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetLocalizeQueries, localizeQueriesAttribute, localizeQueriesUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetLocalizeQueries, localizeQueriesAttribute, localizeQueriesUCIOption),
	}

	localLookupSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       localLookupAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLocalLookup, localLookupAttribute, localLookupUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLocalLookup, localLookupAttribute, localLookupUCIOption),
	}

	localServiceSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       localServiceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetLocalService, localServiceAttribute, localServiceUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetLocalService, localServiceAttribute, localServiceUCIOption),
	}

	readEthersSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       readEthersAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetReadEthers, readEthersAttribute, readEthersUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetReadEthers, readEthersAttribute, readEthersUCIOption),
	}

	rebindLocalhostSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       rebindLocalhostAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetRebindLocalhost, rebindLocalhostAttribute, rebindLocalhostUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetRebindLocalhost, rebindLocalhostAttribute, rebindLocalhostUCIOption),
	}

	rebindProtectionSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       rebindProtectionAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetRebindProtection, rebindProtectionAttribute, rebindProtectionUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetRebindProtection, rebindProtectionAttribute, rebindProtectionUCIOption),
	}

	resolvFileSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       resolvFileAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetResolvFile, resolvFileAttribute, resolvFileUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetResolvFile, resolvFileAttribute, resolvFileUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		authoritativeModeAttribute: authoritativeModeSchemaAttribute,
		domainAttribute:            domainSchemaAttribute,
		domainNeededAttribute:      domainNeededSchemaAttribute,
		ednsPacketMaxAttribute:     ednsPacketMaxSchemaAttribute,
		expandHostsAttribute:       expandHostsSchemaAttribute,
		leaseFileAttribute:         leaseFileSchemaAttribute,
		localizeQueriesAttribute:   localizeQueriesSchemaAttribute,
		localLookupAttribute:       localLookupSchemaAttribute,
		localServiceAttribute:      localServiceSchemaAttribute,
		lucirpcglue.IdAttribute:    lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		readEthersAttribute:        readEthersSchemaAttribute,
		rebindLocalhostAttribute:   rebindLocalhostSchemaAttribute,
		rebindProtectionAttribute:  rebindProtectionSchemaAttribute,
		resolvFileAttribute:        resolvFileSchemaAttribute,
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
	AuthoritativeMode types.Bool   `tfsdk:"authoritative"`
	Domain            types.String `tfsdk:"domain"`
	DomainNeeded      types.Bool   `tfsdk:"domainneeded"`
	EDNSPacketMax     types.Int64  `tfsdk:"ednspacket_max"`
	ExpandHosts       types.Bool   `tfsdk:"expandhosts"`
	Id                types.String `tfsdk:"id"`
	LeaseFile         types.String `tfsdk:"leasefile"`
	LocalizeQueries   types.Bool   `tfsdk:"localise_queries"`
	LocalLookup       types.String `tfsdk:"local"`
	LocalService      types.Bool   `tfsdk:"localservice"`
	ReadEthers        types.Bool   `tfsdk:"readethers"`
	RebindLocalhost   types.Bool   `tfsdk:"rebind_localhost"`
	RebindProtection  types.Bool   `tfsdk:"rebind_protection"`
	ResolvFile        types.String `tfsdk:"resolvfile"`
}

func modelGetAuthoritativeMode(m model) types.Bool { return m.AuthoritativeMode }
func modelGetDomain(m model) types.String          { return m.Domain }
func modelGetDomainNeeded(m model) types.Bool      { return m.DomainNeeded }
func modelGetEDNSPacketMax(m model) types.Int64    { return m.EDNSPacketMax }
func modelGetExpandHosts(m model) types.Bool       { return m.ExpandHosts }
func modelGetId(m model) types.String              { return m.Id }
func modelGetLeaseFile(m model) types.String       { return m.LeaseFile }
func modelGetLocalizeQueries(m model) types.Bool   { return m.LocalizeQueries }
func modelGetLocalLookup(m model) types.String     { return m.LocalLookup }
func modelGetLocalService(m model) types.Bool      { return m.LocalService }
func modelGetReadEthers(m model) types.Bool        { return m.ReadEthers }
func modelGetRebindLocalhost(m model) types.Bool   { return m.RebindLocalhost }
func modelGetRebindProtection(m model) types.Bool  { return m.RebindProtection }
func modelGetResolvFile(m model) types.String      { return m.ResolvFile }

func modelSetAuthoritativeMode(m *model, value types.Bool) { m.AuthoritativeMode = value }
func modelSetDomain(m *model, value types.String)          { m.Domain = value }
func modelSetDomainNeeded(m *model, value types.Bool)      { m.DomainNeeded = value }
func modelSetEDNSPacketMax(m *model, value types.Int64)    { m.EDNSPacketMax = value }
func modelSetExpandHosts(m *model, value types.Bool)       { m.ExpandHosts = value }
func modelSetId(m *model, value types.String)              { m.Id = value }
func modelSetLeaseFile(m *model, value types.String)       { m.LeaseFile = value }
func modelSetLocalizeQueries(m *model, value types.Bool)   { m.LocalizeQueries = value }
func modelSetLocalLookup(m *model, value types.String)     { m.LocalLookup = value }
func modelSetLocalService(m *model, value types.Bool)      { m.LocalService = value }
func modelSetReadEthers(m *model, value types.Bool)        { m.ReadEthers = value }
func modelSetRebindLocalhost(m *model, value types.Bool)   { m.RebindLocalhost = value }
func modelSetRebindProtection(m *model, value types.Bool)  { m.RebindProtection = value }
func modelSetResolvFile(m *model, value types.String)      { m.ResolvFile = value }
