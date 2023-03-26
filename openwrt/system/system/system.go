package system

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	conLogLevelAttribute = "conloglevel"
	conLogLevelUCIOption = "conloglevel"

	cronLogLevelAttribute = "cronloglevel"
	cronLogLevelUCIOption = "cronloglevel"

	descriptionAttribute = "description"
	descriptionUCIOption = "description"

	hostnameAttribute = "hostname"
	hostnameUCIOption = "hostname"

	logSizeAttribute = "log_size"
	logSizeUCIOption = "log_size"

	notesAttribute = "notes"
	notesUCIOption = "notes"

	schemaDescription = "Provides system data about an OpenWrt device"

	timezoneAttribute = "timezone"
	timezoneUCIOption = "timezone"

	ttyLoginAttribute = "ttylogin"
	ttyLoginUCIOption = "ttylogin"

	uciConfig = "system"
	uciType   = "system"

	zonenameAttribute = "zonename"
	zonenameUCIOption = "zonename"
)

var (
	conLogLevelSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "The maximum log level for kernel messages to be logged to the console.",
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetConLogLevel, conLogLevelAttribute, conLogLevelUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetConLogLevel, conLogLevelAttribute, conLogLevelUCIOption),
	}

	cronLogLevelSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "The minimum level for cron messages to be logged to syslog.",
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetCronLogLevel, cronLogLevelAttribute, cronLogLevelUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetCronLogLevel, cronLogLevelAttribute, cronLogLevelUCIOption),
	}

	descriptionSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "The hostname for the system.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDescription, descriptionAttribute, descriptionUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDescription, descriptionAttribute, descriptionUCIOption),
	}

	hostnameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "A short single-line description for the system.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetHostname, hostnameAttribute, hostnameUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetHostname, hostnameAttribute, hostnameUCIOption),
	}

	logSizeSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "Size of the file based log buffer in KiB.",
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetLogSize, logSizeAttribute, logSizeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetLogSize, logSizeAttribute, logSizeUCIOption),
	}

	notesSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "Multi-line free-form text about the system.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetNotes, notesAttribute, notesUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetNotes, notesAttribute, notesUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		conLogLevelAttribute:    conLogLevelSchemaAttribute,
		cronLogLevelAttribute:   cronLogLevelSchemaAttribute,
		descriptionAttribute:    descriptionSchemaAttribute,
		hostnameAttribute:       hostnameSchemaAttribute,
		logSizeAttribute:        logSizeSchemaAttribute,
		notesAttribute:          notesSchemaAttribute,
		timezoneAttribute:       timezoneSchemaAttribute,
		ttyLoginAttribute:       ttyLoginSchemaAttribute,
		zonenameAttribute:       zonenameSchemaAttribute,
	}

	timezoneSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "The POSIX.1 time zone string. This has no corresponding value in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetTimezone, timezoneAttribute, timezoneUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetTimezone, timezoneAttribute, timezoneUCIOption),
	}

	ttyLoginSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "Require authentication for local users to log in the system.",
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetTTYLogin, ttyLoginAttribute, ttyLoginUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetTTYLogin, ttyLoginAttribute, ttyLoginUCIOption),
	}

	zonenameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       "The IANA/Olson time zone string. This corresponds to \"Timezone\" in LuCI. See: https://github.com/openwrt/luci/blob/cd82ccacef78d3bb8b8af6b87dabb9e892e2b2aa/modules/luci-base/luasrc/sys/zoneinfo/tzdata.lua.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetZonename, zonenameAttribute, zonenameUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetZonename, zonenameAttribute, zonenameUCIOption),
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
	ConLogLevel  types.Int64  `tfsdk:"conloglevel"`
	CronLogLevel types.Int64  `tfsdk:"cronloglevel"`
	Description  types.String `tfsdk:"description"`
	Hostname     types.String `tfsdk:"hostname"`
	Id           types.String `tfsdk:"id"`
	LogSize      types.Int64  `tfsdk:"log_size"`
	Notes        types.String `tfsdk:"notes"`
	Timezone     types.String `tfsdk:"timezone"`
	TTYLogin     types.Bool   `tfsdk:"ttylogin"`
	Zonename     types.String `tfsdk:"zonename"`
}

func modelGetConLogLevel(m model) types.Int64  { return m.ConLogLevel }
func modelGetCronLogLevel(m model) types.Int64 { return m.CronLogLevel }
func modelGetDescription(m model) types.String { return m.Description }
func modelGetHostname(m model) types.String    { return m.Hostname }
func modelGetId(m model) types.String          { return m.Id }
func modelGetLogSize(m model) types.Int64      { return m.LogSize }
func modelGetNotes(m model) types.String       { return m.Notes }
func modelGetTimezone(m model) types.String    { return m.Timezone }
func modelGetTTYLogin(m model) types.Bool      { return m.TTYLogin }
func modelGetZonename(m model) types.String    { return m.Zonename }

func modelSetConLogLevel(m *model, value types.Int64)  { m.ConLogLevel = value }
func modelSetCronLogLevel(m *model, value types.Int64) { m.CronLogLevel = value }
func modelSetDescription(m *model, value types.String) { m.Description = value }
func modelSetHostname(m *model, value types.String)    { m.Hostname = value }
func modelSetId(m *model, value types.String)          { m.Id = value }
func modelSetLogSize(m *model, value types.Int64)      { m.LogSize = value }
func modelSetNotes(m *model, value types.String)       { m.Notes = value }
func modelSetTimezone(m *model, value types.String)    { m.Timezone = value }
func modelSetTTYLogin(m *model, value types.Bool)      { m.TTYLogin = value }
func modelSetZonename(m *model, value types.String)    { m.Zonename = value }
