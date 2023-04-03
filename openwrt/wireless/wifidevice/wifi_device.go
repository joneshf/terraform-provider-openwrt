package wifidevice

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	bandAttribute            = "band"
	bandAttributeDescription = `Channel width. Must be one of: "2g", "5g", "6g".`
	band2G                   = "2g"
	band5G                   = "5g"
	band6G                   = "6g"
	bandUCIOption            = "band"

	cellDensityAttribute            = "cell_density"
	cellDensityAttributeDescription = "Configures data rates based on the coverage cell density. Must be one of 0, 1, 2, 3."
	cellDensityDisabled             = 0
	cellDensityHigh                 = 2
	cellDensityNormal               = 1
	cellDensityUCIOption            = "cell_density"
	cellDensityVeryHigh             = 3

	channelAttribute            = "channel"
	channelAttributeDescription = `The wireless channel. Currently, only "auto" is supported.`
	channelAuto                 = "auto"
	channelUCIOption            = "channel"

	countryCodeAttribute            = "country"
	countryCodeAttributeDescription = `Two-digit country code. E.g. "US".`
	countryCodeUCIOption            = "country"

	htModeAttribute            = "htmode"
	htModeAttributeDescription = `Channel width. Must be one of: "HE20", "HE40", "HE80", "HE160", "HT20", "HT40", "HT40-", "HT40+", "NONE", "VHT20", "VHT40", "VHT80", "VHT160".`
	htModeHE160                = "HE160"
	htModeHE20                 = "HE20"
	htModeHE40                 = "HE40"
	htModeHE80                 = "HE80"
	htModeHT20                 = "HT20"
	htModeHT40                 = "HT40"
	htModeHT40Minus            = "HT40-"
	htModeHT40Plus             = "HT40+"
	htModeNone                 = "NONE"
	htModeUCIOption            = "htmode"
	htModeVHT160               = "VHT160"
	htModeVHT20                = "VHT20"
	htModeVHT40                = "VHT40"
	htModeVHT80                = "VHT80"

	pathAttribute            = "path"
	pathAttributeDescription = "Path of the device in `/sys/devices`."
	pathUCIOption            = "path"

	schemaDescription = "The physical radio device."

	typeAttribute            = "type"
	typeAttributeDescription = `The type of device. Currently only "mac80211" is supported.`
	typeMac80211             = "mac80211"
	typeUCIOption            = "type"

	uciConfig = "wireless"
	uciType   = "wifi-device"
)

var (
	bandSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       bandAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetBand, bandAttribute, bandUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetBand, bandAttribute, bandUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				band2G,
				band5G,
				band6G,
			),
		},
	}

	cellDensitySchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       cellDensityAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetCellDensity, cellDensityAttribute, cellDensityUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetCellDensity, cellDensityAttribute, cellDensityUCIOption),
		Validators: []validator.Int64{
			int64validator.OneOf(
				cellDensityDisabled,
				cellDensityNormal,
				cellDensityHigh,
				cellDensityVeryHigh,
			),
		},
	}

	channelSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       channelAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetChannel, channelAttribute, channelUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetChannel, channelAttribute, channelUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				channelAuto,
			),
		},
	}

	countryCodeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       countryCodeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetCountryCode, countryCodeAttribute, countryCodeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetCountryCode, countryCodeAttribute, countryCodeUCIOption),
		Validators: []validator.String{
			stringvalidator.LengthBetween(2, 2),
		},
	}

	htModeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       htModeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetHTMode, htModeAttribute, htModeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetHTMode, htModeAttribute, htModeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				htModeHE160,
				htModeHE20,
				htModeHE40,
				htModeHE80,
				htModeHT20,
				htModeHT40,
				htModeHT40Minus,
				htModeHT40Plus,
				htModeNone,
				htModeVHT160,
				htModeVHT20,
				htModeVHT40,
				htModeVHT80,
			),
		},
	}

	pathSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       pathAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetPath, pathAttribute, pathUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetPath, pathAttribute, pathUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		bandAttribute:           bandSchemaAttribute,
		cellDensityAttribute:    cellDensitySchemaAttribute,
		channelAttribute:        channelSchemaAttribute,
		countryCodeAttribute:    countryCodeSchemaAttribute,
		lucirpcglue.IdAttribute: lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		htModeAttribute:         htModeSchemaAttribute,
		pathAttribute:           pathSchemaAttribute,
		typeAttribute:           typeSchemaAttribute,
	}

	typeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       typeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetType, typeAttribute, typeUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetType, typeAttribute, typeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				typeMac80211,
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
	Band        types.String `tfsdk:"band"`
	CellDensity types.Int64  `tfsdk:"cell_density"`
	Channel     types.String `tfsdk:"channel"`
	CountryCode types.String `tfsdk:"country"`
	HTMode      types.String `tfsdk:"htmode"`
	Id          types.String `tfsdk:"id"`
	Path        types.String `tfsdk:"path"`
	Type        types.String `tfsdk:"type"`
}

func modelGetBand(m model) types.String        { return m.Band }
func modelGetCellDensity(m model) types.Int64  { return m.CellDensity }
func modelGetChannel(m model) types.String     { return m.Channel }
func modelGetCountryCode(m model) types.String { return m.CountryCode }
func modelGetHTMode(m model) types.String      { return m.HTMode }
func modelGetId(m model) types.String          { return m.Id }
func modelGetPath(m model) types.String        { return m.Path }
func modelGetType(m model) types.String        { return m.Type }

func modelSetBand(m *model, value types.String)        { m.Band = value }
func modelSetCellDensity(m *model, value types.Int64)  { m.CellDensity = value }
func modelSetChannel(m *model, value types.String)     { m.Channel = value }
func modelSetCountryCode(m *model, value types.String) { m.CountryCode = value }
func modelSetHTMode(m *model, value types.String)      { m.HTMode = value }
func modelSetId(m *model, value types.String)          { m.Id = value }
func modelSetPath(m *model, value types.String)        { m.Path = value }
func modelSetType(m *model, value types.String)        { m.Type = value }
