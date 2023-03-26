package networkswitch

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	enableMirrorReceivedAttribute            = "enable_mirror_rx"
	enableMirrorReceivedAttributeDescription = "Mirror received packets from the `mirror_source_port` to the `mirror_monitor_port`."
	enableMirrorReceivedUCIOption            = "enable_mirror_rx"

	enableMirrorTransmittedAttribute            = "enable_mirror_tx"
	enableMirrorTransmittedAttributeDescription = "Mirror transmitted packets from the `mirror_source_port` to the `mirror_monitor_port`."
	enableMirrorTransmittedUCIOption            = "enable_mirror_tx"

	enableVLANAttribute            = "enable_vlan"
	enableVLANAttributeDescription = "Enables VLAN functionality."
	enableVLANUCIOption            = "enable_vlan"

	mirrorMonitorPortAttribute            = "mirror_monitor_port"
	mirrorMonitorPortAttributeDescription = "Switch port to which packets are mirrored."
	mirrorMonitorPortUCIOption            = "mirror_monitor_port"

	mirrorSourcePortAttribute            = "mirror_source_port"
	mirrorSourcePortAttributeDescription = "Switch port from which packets are mirrored."
	mirrorSourcePortUCIOption            = "mirror_source_port"

	nameAttribute            = "name"
	nameAttributeDescription = "Name of the switch. This name is what is shown in LuCI or the `name` field in Terraform. This is not the UCI config name."
	nameUCIOption            = "name"

	resetAttribute            = "reset"
	resetAttributeDescription = "Reset the switch."
	resetUCIOption            = "reset"

	schemaDescription = "Legacy `swconfig` configuration"

	uciConfig = "network"
	uciType   = "switch"
)

var (
	enableMirrorReceivedSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enableMirrorReceivedAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnableMirrorReceived, enableMirrorReceivedAttribute, enableMirrorReceivedUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnableMirrorReceived, enableMirrorReceivedAttribute, enableMirrorReceivedUCIOption),
	}

	enableMirrorTransmittedSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enableMirrorTransmittedAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnableMirrorTransmitted, enableMirrorTransmittedAttribute, enableMirrorTransmittedUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnableMirrorTransmitted, enableMirrorTransmittedAttribute, enableMirrorTransmittedUCIOption),
	}

	enableVLANSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       enableVLANAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetEnableVLAN, enableVLANAttribute, enableVLANUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetEnableVLAN, enableVLANAttribute, enableVLANUCIOption),
	}

	mirrorMonitorPortSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mirrorMonitorPortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMirrorMonitorPort, mirrorMonitorPortAttribute, mirrorMonitorPortUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMirrorMonitorPort, mirrorMonitorPortAttribute, mirrorMonitorPortUCIOption),
		Validators: []validator.Int64{
			int64validator.Any(
				lucirpcglue.RequiresAttributeEqualBool(
					path.MatchRoot(enableMirrorReceivedAttribute),
					true,
				),
				lucirpcglue.RequiresAttributeEqualBool(
					path.MatchRoot(enableMirrorTransmittedAttribute),
					true,
				),
			),
		},
	}

	mirrorSourcePortSchemaAttribute = lucirpcglue.Int64SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       mirrorSourcePortAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(modelSetMirrorSourcePort, mirrorSourcePortAttribute, mirrorSourcePortUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(modelGetMirrorSourcePort, mirrorSourcePortAttribute, mirrorSourcePortUCIOption),
		Validators: []validator.Int64{
			int64validator.Any(
				lucirpcglue.RequiresAttributeEqualBool(
					path.MatchRoot(enableMirrorReceivedAttribute),
					true,
				),
				lucirpcglue.RequiresAttributeEqualBool(
					path.MatchRoot(enableMirrorTransmittedAttribute),
					true,
				),
			),
		},
	}

	nameSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       nameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetName, nameAttribute, nameUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetName, nameAttribute, nameUCIOption),
	}

	resetSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       resetAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetReset, resetAttribute, resetUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetReset, resetAttribute, resetUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		enableMirrorReceivedAttribute:    enableMirrorReceivedSchemaAttribute,
		enableMirrorTransmittedAttribute: enableMirrorTransmittedSchemaAttribute,
		enableVLANAttribute:              enableVLANSchemaAttribute,
		lucirpcglue.IdAttribute:          lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		mirrorMonitorPortAttribute:       mirrorMonitorPortSchemaAttribute,
		mirrorSourcePortAttribute:        mirrorSourcePortSchemaAttribute,
		nameAttribute:                    nameSchemaAttribute,
		resetAttribute:                   resetSchemaAttribute,
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
	EnableMirrorReceived    types.Bool   `tfsdk:"enable_mirror_rx"`
	EnableMirrorTransmitted types.Bool   `tfsdk:"enable_mirror_tx"`
	EnableVLAN              types.Bool   `tfsdk:"enable_vlan"`
	Id                      types.String `tfsdk:"id"`
	MirrorMonitorPort       types.Int64  `tfsdk:"mirror_monitor_port"`
	MirrorSourcePort        types.Int64  `tfsdk:"mirror_source_port"`
	Name                    types.String `tfsdk:"name"`
	Reset                   types.Bool   `tfsdk:"reset"`
}

func modelGetEnableMirrorReceived(m model) types.Bool    { return m.EnableMirrorReceived }
func modelGetEnableMirrorTransmitted(m model) types.Bool { return m.EnableMirrorTransmitted }
func modelGetEnableVLAN(m model) types.Bool              { return m.EnableVLAN }
func modelGetId(m model) types.String                    { return m.Id }
func modelGetMirrorMonitorPort(m model) types.Int64      { return m.MirrorMonitorPort }
func modelGetMirrorSourcePort(m model) types.Int64       { return m.MirrorSourcePort }
func modelGetName(m model) types.String                  { return m.Name }
func modelGetReset(m model) types.Bool                   { return m.Reset }

func modelSetEnableMirrorReceived(m *model, value types.Bool)    { m.EnableMirrorReceived = value }
func modelSetEnableMirrorTransmitted(m *model, value types.Bool) { m.EnableMirrorTransmitted = value }
func modelSetEnableVLAN(m *model, value types.Bool)              { m.EnableVLAN = value }
func modelSetId(m *model, value types.String)                    { m.Id = value }
func modelSetMirrorMonitorPort(m *model, value types.Int64)      { m.MirrorMonitorPort = value }
func modelSetMirrorSourcePort(m *model, value types.Int64)       { m.MirrorSourcePort = value }
func modelSetName(m *model, value types.String)                  { m.Name = value }
func modelSetReset(m *model, value types.Bool)                   { m.Reset = value }
