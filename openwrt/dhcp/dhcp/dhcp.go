package dhcp

import (
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
	dhcpv4ModeAttribute            = "dhcpv4"
	dhcpv4ModeAttributeDescription = `The mode of the DHCPv4 server. Must be one of: "disabled", "server".`
	dhcpv4ModeDisabled             = "disabled"
	dhcpv4ModeServer               = "server"
	dhcpv4ModeUCIOption            = "dhcpv4"

	dhcpv6ModeAttribute            = "dhcpv6"
	dhcpv6ModeAttributeDescription = `The mode of the DHCPv6 server. Must be one of: "disabled", "relay", "server".`
	dhcpv6ModeDisabled             = "disabled"
	dhcpv6ModeRelay                = "relay"
	dhcpv6ModeServer               = "server"
	dhcpv6ModeUCIOption            = "dhcpv6"

	forceAttribute            = "force"
	forceAttributeDescription = "Forces DHCP serving on the specified interface even if another DHCP server is detected on the same network segment."
	forceUCIOption            = "force"

	ignoreAttribute            = "ignore"
	ignoreAttributeDescription = "Specifies whether dnsmasq should ignore this pool."
	ignoreUCIOption            = "ignore"

	interfaceAttribute            = "interface"
	interfaceAttributeDescription = "The interface associated with this DHCP address pool. This name is what the interface is known as in UCI, or the `id` field in Terraform. Required if `ignore` is not `true`."
	interfaceUCIOption            = "interface"

	leaseTimeAttribute            = "leasetime"
	leaseTimeAttributeDescription = "The lease time of addresses handed out to clients. E.g. `12h`, or `30m`. Required if `ignore` is not `true`."
	leaseTimeUCIOption            = "leasetime"

	limitAttribute            = "limit"
	limitAttributeDescription = "Specifies the size of the address pool. E.g. With start = 100, and limit = 150, the maximum address will be 249. Required if `ignore` is not `true`."
	limitUCIOption            = "limit"

	routerAdvertisementFlagsAttribute            = "ra_flags"
	routerAdvertisementFlagsAttributeDescription = `Router Advertisement flags to include in messages. Must be one of: "home-agent", "managed-config", "none", "other-config".`
	routerAdvertisementFlagsHomeAgent            = "home-agent"
	routerAdvertisementFlagsManagedConfig        = "managed-config"
	routerAdvertisementFlagsNone                 = "none"
	routerAdvertisementFlagsOtherConfig          = "other-config"
	routerAdvertisementFlagsUCIOption            = "ra_flags"

	routerAdvertisementModeAttribute            = "ra"
	routerAdvertisementModeAttributeDescription = `The mode of Router Advertisements. Must be one of: "disabled", "relay", "server".`
	routerAdvertisementModeDisabled             = "disabled"
	routerAdvertisementModeRelay                = "relay"
	routerAdvertisementModeServer               = "server"
	routerAdvertisementModeUCIOption            = "ra"

	schemaDescription = "Per interface lease pools and settings for serving DHCP requests."

	startAttribute            = "start"
	startAttributeDescription = "Specifies the offset from the network address of the underlying interface to calculate the minimum address that may be leased to clients. It may be greater than 255 to span subnets. Required if `ignore` is not `true`."
	startUCIOption            = "start"

	uciConfig = "dhcp"
	uciType   = "dhcp"
)

var (
	dhcpv4ModeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       dhcpv4ModeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDHCPv4Mode, dhcpv4ModeAttribute, dhcpv4ModeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDHCPv4Mode, dhcpv4ModeAttribute, dhcpv4ModeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				dhcpv4ModeDisabled,
				dhcpv4ModeServer,
			),
		},
	}

	dhcpv6ModeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       dhcpv6ModeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDHCPv6Mode, dhcpv6ModeAttribute, dhcpv6ModeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDHCPv6Mode, dhcpv6ModeAttribute, dhcpv6ModeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				dhcpv6ModeDisabled,
				dhcpv6ModeRelay,
				dhcpv6ModeServer,
			),
		},
	}

	forceSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       forceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetForce, forceAttribute, forceUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetForce, forceAttribute, forceUCIOption),
	}

	ignoreSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ignoreAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetIgnore, ignoreAttribute, ignoreUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetIgnore, ignoreAttribute, ignoreUCIOption),
	}

	interfaceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       interfaceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetInterface, interfaceAttribute, interfaceUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetInterface, interfaceAttribute, interfaceUCIOption),
		Validators: []validator.String{
			lucirpcglue.RequiredIfAttributeNotEqualBool(path.MatchRoot(ignoreAttribute), true),
		},
	}

	leaseTimeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       leaseTimeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLeaseTime, leaseTimeAttribute, leaseTimeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLeaseTime, leaseTimeAttribute, leaseTimeUCIOption),
		Validators: []validator.String{
			lucirpcglue.RequiredIfAttributeNotEqualBool(path.MatchRoot(ignoreAttribute), true),
		},
	}

	limitSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       limitAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetLimit, limitAttribute, limitUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetLimit, limitAttribute, limitUCIOption),
		Validators: []validator.Int64{
			lucirpcglue.RequiredIfAttributeNotEqualBool(path.MatchRoot(ignoreAttribute), true),
		},
	}

	routerAdvertisementFlagsSchemaAttribute = lucirpcglue.SetStringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       routerAdvertisementFlagsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionSetString(modelSetRouterAdvertisementFlags, routerAdvertisementFlagsAttribute, routerAdvertisementFlagsUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionSetString(modelGetRouterAdvertisementFlags, routerAdvertisementFlagsAttribute, routerAdvertisementFlagsUCIOption),
		Validators: []validator.Set{
			setvalidator.ValueStringsAre(
				stringvalidator.OneOf(
					routerAdvertisementFlagsHomeAgent,
					routerAdvertisementFlagsManagedConfig,
					routerAdvertisementFlagsNone,
					routerAdvertisementFlagsOtherConfig,
				),
			),
		},
	}

	routerAdvertisementModeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       routerAdvertisementModeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetRouterAdvertisementMode, routerAdvertisementModeAttribute, routerAdvertisementModeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetRouterAdvertisementMode, routerAdvertisementModeAttribute, routerAdvertisementModeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				routerAdvertisementModeDisabled,
				routerAdvertisementModeRelay,
				routerAdvertisementModeServer,
			),
		},
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		dhcpv4ModeAttribute:               dhcpv4ModeSchemaAttribute,
		dhcpv6ModeAttribute:               dhcpv6ModeSchemaAttribute,
		forceAttribute:                    forceSchemaAttribute,
		ignoreAttribute:                   ignoreSchemaAttribute,
		interfaceAttribute:                interfaceSchemaAttribute,
		leaseTimeAttribute:                leaseTimeSchemaAttribute,
		limitAttribute:                    limitSchemaAttribute,
		lucirpcglue.IdAttribute:           lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		routerAdvertisementFlagsAttribute: routerAdvertisementFlagsSchemaAttribute,
		routerAdvertisementModeAttribute:  routerAdvertisementModeSchemaAttribute,
		startAttribute:                    startSchemaAttribute,
	}

	startSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       startAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetStart, startAttribute, startUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetStart, startAttribute, startUCIOption),
		Validators: []validator.Int64{
			lucirpcglue.RequiredIfAttributeNotEqualBool(path.MatchRoot(ignoreAttribute), true),
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
	DHCPv4Mode               types.String `tfsdk:"dhcpv4"`
	DHCPv6Mode               types.String `tfsdk:"dhcpv6"`
	Force                    types.Bool   `tfsdk:"force"`
	Id                       types.String `tfsdk:"id"`
	Ignore                   types.Bool   `tfsdk:"ignore"`
	Interface                types.String `tfsdk:"interface"`
	LeaseTime                types.String `tfsdk:"leasetime"`
	Limit                    types.Int64  `tfsdk:"limit"`
	RouterAdvertisementFlags types.Set    `tfsdk:"ra_flags"`
	RouterAdvertisementMode  types.String `tfsdk:"ra"`
	Start                    types.Int64  `tfsdk:"start"`
}

func modelGetDHCPv4Mode(m model) types.String              { return m.DHCPv4Mode }
func modelGetDHCPv6Mode(m model) types.String              { return m.DHCPv6Mode }
func modelGetForce(m model) types.Bool                     { return m.Force }
func modelGetId(m model) types.String                      { return m.Id }
func modelGetIgnore(m model) types.Bool                    { return m.Ignore }
func modelGetInterface(m model) types.String               { return m.Interface }
func modelGetLeaseTime(m model) types.String               { return m.LeaseTime }
func modelGetLimit(m model) types.Int64                    { return m.Limit }
func modelGetRouterAdvertisementFlags(m model) types.Set   { return m.RouterAdvertisementFlags }
func modelGetRouterAdvertisementMode(m model) types.String { return m.RouterAdvertisementMode }
func modelGetStart(m model) types.Int64                    { return m.Start }

func modelSetDHCPv4Mode(m *model, value types.String)              { m.DHCPv4Mode = value }
func modelSetDHCPv6Mode(m *model, value types.String)              { m.DHCPv6Mode = value }
func modelSetForce(m *model, value types.Bool)                     { m.Force = value }
func modelSetId(m *model, value types.String)                      { m.Id = value }
func modelSetIgnore(m *model, value types.Bool)                    { m.Ignore = value }
func modelSetInterface(m *model, value types.String)               { m.Interface = value }
func modelSetLeaseTime(m *model, value types.String)               { m.LeaseTime = value }
func modelSetLimit(m *model, value types.Int64)                    { m.Limit = value }
func modelSetRouterAdvertisementFlags(m *model, value types.Set)   { m.RouterAdvertisementFlags = value }
func modelSetRouterAdvertisementMode(m *model, value types.String) { m.RouterAdvertisementMode = value }
func modelSetStart(m *model, value types.Int64)                    { m.Start = value }
