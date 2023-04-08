package odhcpd

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	leaseFileAttribute            = "leasefile"
	leaseFileAttributeDescription = "Location of the lease/hostfile for DHCPv4 and DHCPv6."
	leaseFileUCIOption            = "leasefile"

	leaseTriggerAttribute            = "leasetrigger"
	leaseTriggerAttributeDescription = "Location of the lease trigger script."
	leaseTriggerUCIOption            = "leasetrigger"

	legacyAttribute            = "legacy"
	legacyAttributeDescription = "Enable DHCPv4 if the 'dhcp' section constains a `start` option, but no `dhcpv4` option set."
	legacyUCIOption            = "legacy"

	logLevelAttribute            = "loglevel"
	logLevelAttributeDescription = "Syslog level priority (0-7)."
	logLevelUCIOption            = "loglevel"

	mainDHCPAttribute            = "maindhcp"
	mainDHCPAttributeDescription = "Use odhcpd as the main DHCPv4 service."
	mainDHCPUCIOption            = "maindhcp"

	schemaDescription = "An embedded DHCP/DHCPv6/RA server & NDP relay."

	uciConfig = "dhcp"
	uciType   = "odhcpd"
)

var (
	leaseFileSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       leaseFileAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLeaseFile, leaseFileAttribute, leaseFileUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLeaseFile, leaseFileAttribute, leaseFileUCIOption),
	}

	leaseTriggerSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       leaseTriggerAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetLeaseTrigger, leaseTriggerAttribute, leaseTriggerUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetLeaseTrigger, leaseTriggerAttribute, leaseTriggerUCIOption),
	}

	legacySchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       legacyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetLegacy, legacyAttribute, legacyUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetLegacy, legacyAttribute, legacyUCIOption),
	}

	logLevelSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       logLevelAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetLogLevel, logLevelAttribute, logLevelUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetLogLevel, logLevelAttribute, logLevelUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(0, 7),
		},
	}

	mainDHCPSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mainDHCPAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetMainDHCP, mainDHCPAttribute, mainDHCPUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetMainDHCP, mainDHCPAttribute, mainDHCPUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		leaseFileAttribute:      leaseFileSchemaAttribute,
		leaseTriggerAttribute:   leaseTriggerSchemaAttribute,
		legacyAttribute:         legacySchemaAttribute,
		logLevelAttribute:       logLevelSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		mainDHCPAttribute:       mainDHCPSchemaAttribute,
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
	Id           types.String `tfsdk:"id"`
	LeaseFile    types.String `tfsdk:"leasefile"`
	LeaseTrigger types.String `tfsdk:"leasetrigger"`
	Legacy       types.Bool   `tfsdk:"legacy"`
	LogLevel     types.Int64  `tfsdk:"loglevel"`
	MainDHCP     types.Bool   `tfsdk:"maindhcp"`
}

func modelGetId(m model) types.String           { return m.Id }
func modelGetLeaseFile(m model) types.String    { return m.LeaseFile }
func modelGetLeaseTrigger(m model) types.String { return m.LeaseTrigger }
func modelGetLegacy(m model) types.Bool         { return m.Legacy }
func modelGetLogLevel(m model) types.Int64      { return m.LogLevel }
func modelGetMainDHCP(m model) types.Bool       { return m.MainDHCP }

func modelSetId(m *model, value types.String)           { m.Id = value }
func modelSetLeaseFile(m *model, value types.String)    { m.LeaseFile = value }
func modelSetLeaseTrigger(m *model, value types.String) { m.LeaseTrigger = value }
func modelSetLegacy(m *model, value types.Bool)         { m.Legacy = value }
func modelSetLogLevel(m *model, value types.Int64)      { m.LogLevel = value }
func modelSetMainDHCP(m *model, value types.Bool)       { m.MainDHCP = value }
