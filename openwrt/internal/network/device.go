package network

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	deviceBridgePortsAttribute            = "ports"
	deviceBridgePortsAttributeDescription = "Specifies the wired ports to attach to this bridge."
	deviceBridgePortsUCIOption            = "ports"

	deviceBringUpEmptyBridgeAttribute            = "bridge_empty"
	deviceBringUpEmptyBridgeAttributeDescription = "Bring up the bridge device even if no ports are attached"
	deviceBringUpEmptyBridgeUCIOption            = "bridge_empty"

	deviceDADTransmitsAttribute            = "dadtransmits"
	deviceDADTransmitsAttributeDescription = "Amount of Duplicate Address Detection probes to send"
	deviceDADTransmitsUCIOption            = "dadtransmits"

	deviceEnableIPv6Attribute            = "ipv6"
	deviceEnableIPv6AttributeDescription = "Enable IPv6 for the device."
	deviceEnableIPv6UCIOption            = "ipv6"

	deviceIdAttribute            = "id"
	deviceIdAttributeDescription = "Name of the section. This name is only used when interacting with UCI directly."
	deviceIdUCISection           = ".name"

	deviceMacAddressAttribute            = "macaddr"
	deviceMacAddressAttributeDescription = "MAC Address of the device."
	deviceMacAddressUCIOption            = "macaddr"

	deviceMTUAttribute            = "mtu"
	deviceMTUAttributeDescription = "Maximum Transmissible Unit."
	deviceMTUUCIOption            = "mtu"

	deviceMTU6Attribute            = "mtu6"
	deviceMTU6AttributeDescription = "Maximum Transmissible Unit for IPv6."
	deviceMTU6UCIOption            = "mtu6"

	deviceNameAttribute            = "name"
	deviceNameAttributeDescription = "Name of the device. This name is referenced in other network configuration."
	deviceNameUCIOption            = "name"

	deviceSchemaDescription = `A physical or virtual "device" in OpenWrt jargon. Commonly referred to as an "interface" in other networking jargon.`

	deviceTXQueueLengthAttribute            = "txqueuelen"
	deviceTXQueueLengthAttributeDescription = "Transmission queue length."
	deviceTXQueueLengthUCIOption            = "txqueuelen"

	deviceTypeAttribute            = "type"
	deviceTypeAttributeDescription = `The type of device. Currently, only "bridge" is supported.`
	deviceTypeBridge               = "bridge"
	deviceTypeUCIOption            = "type"

	deviceTypeName  = "network_device"
	deviceUCIConfig = "network"
	deviceUCIType   = "device"
)

var (
	deviceBridgePortsSchemaAttribute = lucirpcglue.SetStringSchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceBridgePortsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionSetString(deviceModelSetBridgePorts, deviceBridgePortsAttribute, deviceBridgePortsUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionSetString(deviceModelGetBridgePorts, deviceBridgePortsAttribute, deviceBridgePortsUCIOption),
		Validators: []validator.Set{
			setvalidator.SizeAtLeast(1),
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(deviceTypeAttribute),
				deviceTypeBridge,
			),
		},
	}

	deviceBringUpEmptyBridgeSchemaAttribute = lucirpcglue.BoolSchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceBringUpEmptyBridgeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(deviceModelSetBringUpEmptyBridge, deviceBringUpEmptyBridgeAttribute, deviceBringUpEmptyBridgeUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(deviceModelGetBringUpEmptyBridge, deviceBringUpEmptyBridgeAttribute, deviceBringUpEmptyBridgeUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(deviceTypeAttribute),
				deviceTypeBridge,
			),
		},
	}

	deviceDADTransmitsSchemaAttribute = lucirpcglue.Int64SchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceDADTransmitsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(deviceModelSetDADTransmits, deviceDADTransmitsAttribute, deviceDADTransmitsUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(deviceModelGetDADTransmits, deviceDADTransmitsAttribute, deviceDADTransmitsUCIOption),
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
			lucirpcglue.RequiresAttributeEqualBool(
				path.MatchRoot(deviceEnableIPv6Attribute),
				true,
			),
		},
	}

	deviceEnableIPv6SchemaAttribute = lucirpcglue.BoolSchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceEnableIPv6AttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(deviceModelSetEnableIPv6, deviceEnableIPv6Attribute, deviceEnableIPv6UCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(deviceModelGetEnableIPv6, deviceEnableIPv6Attribute, deviceEnableIPv6UCIOption),
	}

	deviceIdSchemaAttribute = lucirpcglue.StringSchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		DataSourceExistence: lucirpcglue.Required,
		Description:         deviceIdAttributeDescription,
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model deviceModel,
		) (context.Context, deviceModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetMetadataString(ctx, fullTypeName, terraformType, section, deviceIdUCISection)
			model.Id = value
			return ctx, model, diagnostics
		},
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model deviceModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			ctx = logger.SetFieldString(ctx, fullTypeName, lucirpcglue.ResourceTerraformType, deviceIdAttribute, model.Id)
			return ctx, options, diag.Diagnostics{}
		},
	}

	deviceMacAddressSchemaAttribute = lucirpcglue.StringSchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceMacAddressAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(deviceModelSetMacAddress, deviceMacAddressAttribute, deviceMacAddressUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(deviceModelGetMacAddress, deviceMacAddressAttribute, deviceMacAddressUCIOption),
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile("^([[:xdigit:]][[:xdigit:]]:){5}[[:xdigit:]][[:xdigit:]]$"),
				`must be a valid MAC address (e.g. "12:34:56:78:90:ab")`,
			),
		},
	}

	deviceMTUSchemaAttribute = lucirpcglue.Int64SchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceMTUAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(deviceModelSetMTU, deviceMTUAttribute, deviceMTUUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(deviceModelGetMTU, deviceMTUAttribute, deviceMTUUCIOption),
		Validators: []validator.Int64{
			int64validator.Between(576, 9200),
		},
	}

	deviceMTU6SchemaAttribute = lucirpcglue.Int64SchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceMTU6AttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(deviceModelSetMTU6, deviceMTU6Attribute, deviceMTU6UCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(deviceModelGetMTU6, deviceMTU6Attribute, deviceMTU6UCIOption),
		Validators: []validator.Int64{
			int64validator.Between(576, 9200),
			lucirpcglue.RequiresAttributeEqualBool(
				path.MatchRoot(deviceEnableIPv6Attribute),
				true,
			),
		},
	}

	deviceNameSchemaAttribute = lucirpcglue.StringSchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceNameAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(deviceModelSetName, deviceNameAttribute, deviceNameUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(deviceModelGetName, deviceNameAttribute, deviceNameUCIOption),
	}

	deviceSchemaAttributes = map[string]lucirpcglue.SchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		deviceBridgePortsAttribute:        deviceBridgePortsSchemaAttribute,
		deviceBringUpEmptyBridgeAttribute: deviceBringUpEmptyBridgeSchemaAttribute,
		deviceDADTransmitsAttribute:       deviceDADTransmitsSchemaAttribute,
		deviceEnableIPv6Attribute:         deviceEnableIPv6SchemaAttribute,
		deviceIdAttribute:                 deviceIdSchemaAttribute,
		deviceMacAddressAttribute:         deviceMacAddressSchemaAttribute,
		deviceMTUAttribute:                deviceMTUSchemaAttribute,
		deviceMTU6Attribute:               deviceMTU6SchemaAttribute,
		deviceNameAttribute:               deviceNameSchemaAttribute,
		deviceTXQueueLengthAttribute:      deviceTXQueueLengthSchemaAttribute,
		deviceTypeAttribute:               deviceTypeSchemaAttribute,
	}

	deviceTXQueueLengthSchemaAttribute = lucirpcglue.Int64SchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceTXQueueLengthAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionInt64(deviceModelSetTXQueueLength, deviceTXQueueLengthAttribute, deviceTXQueueLengthUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionInt64(deviceModelGetTXQueueLength, deviceTXQueueLengthAttribute, deviceTXQueueLengthUCIOption),
		Validators: []validator.Int64{
			int64validator.AtLeast(1),
		},
	}

	deviceTypeSchemaAttribute = lucirpcglue.StringSchemaAttribute[deviceModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       deviceTypeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(deviceModelSetType, deviceTypeAttribute, deviceTypeUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(deviceModelGetType, deviceTypeAttribute, deviceTypeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				deviceTypeBridge,
			),
		},
	}
)

type deviceModel struct {
	BridgePorts        types.Set    `tfsdk:"ports"`
	BringUpEmptyBridge types.Bool   `tfsdk:"bridge_empty"`
	DADTransmits       types.Int64  `tfsdk:"dadtransmits"`
	EnableIPv6         types.Bool   `tfsdk:"ipv6"`
	Id                 types.String `tfsdk:"id"`
	MacAddress         types.String `tfsdk:"macaddr"`
	MTU                types.Int64  `tfsdk:"mtu"`
	MTU6               types.Int64  `tfsdk:"mtu6"`
	Name               types.String `tfsdk:"name"`
	TXQueueLength      types.Int64  `tfsdk:"txqueuelen"`
	Type               types.String `tfsdk:"type"`
}

func (m deviceModel) generateAPIBody(
	ctx context.Context,
	fullTypeName string,
) (context.Context, string, map[string]json.RawMessage, diag.Diagnostics) {
	tflog.Info(ctx, "Generating API request body")
	var diagnostics diag.Diagnostics
	allDiagnostics := diag.Diagnostics{}
	options := map[string]json.RawMessage{}

	tflog.Debug(ctx, "Handling attributes")
	id := m.Id.ValueString()
	for _, attribute := range deviceSchemaAttributes {
		ctx, options, diagnostics = attribute.Upsert(ctx, fullTypeName, lucirpcglue.ResourceTerraformType, options, m)
		allDiagnostics.Append(diagnostics...)
	}

	return ctx, id, options, allDiagnostics
}

func deviceModelGetBridgePorts(model deviceModel) types.Set         { return model.BridgePorts }
func deviceModelGetBringUpEmptyBridge(model deviceModel) types.Bool { return model.BringUpEmptyBridge }
func deviceModelGetDADTransmits(model deviceModel) types.Int64      { return model.DADTransmits }
func deviceModelGetEnableIPv6(model deviceModel) types.Bool         { return model.EnableIPv6 }
func deviceModelGetMacAddress(model deviceModel) types.String       { return model.MacAddress }
func deviceModelGetMTU(model deviceModel) types.Int64               { return model.MTU }
func deviceModelGetMTU6(model deviceModel) types.Int64              { return model.MTU6 }
func deviceModelGetName(model deviceModel) types.String             { return model.Name }
func deviceModelGetTXQueueLength(model deviceModel) types.Int64     { return model.TXQueueLength }
func deviceModelGetType(model deviceModel) types.String             { return model.Type }

func deviceModelSetBridgePorts(model *deviceModel, value types.Set) { model.BridgePorts = value }
func deviceModelSetBringUpEmptyBridge(model *deviceModel, value types.Bool) {
	model.BringUpEmptyBridge = value
}
func deviceModelSetDADTransmits(model *deviceModel, value types.Int64)  { model.DADTransmits = value }
func deviceModelSetEnableIPv6(model *deviceModel, value types.Bool)     { model.EnableIPv6 = value }
func deviceModelSetMacAddress(model *deviceModel, value types.String)   { model.MacAddress = value }
func deviceModelSetMTU(model *deviceModel, value types.Int64)           { model.MTU = value }
func deviceModelSetMTU6(model *deviceModel, value types.Int64)          { model.MTU6 = value }
func deviceModelSetName(model *deviceModel, value types.String)         { model.Name = value }
func deviceModelSetTXQueueLength(model *deviceModel, value types.Int64) { model.TXQueueLength = value }
func deviceModelSetType(model *deviceModel, value types.String)         { model.Type = value }

func readDeviceModel(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	client lucirpc.Client,
	sectionName string,
) (context.Context, deviceModel, diag.Diagnostics) {
	tflog.Info(ctx, "Reading device model")
	var (
		allDiagnostics diag.Diagnostics
		model          deviceModel
	)

	section, diagnostics := lucirpcglue.GetSection(ctx, client, deviceUCIConfig, sectionName)
	allDiagnostics.Append(diagnostics...)
	if allDiagnostics.HasError() {
		return ctx, model, allDiagnostics
	}

	for _, attribute := range deviceSchemaAttributes {
		ctx, model, diagnostics = attribute.Read(ctx, fullTypeName, terraformType, section, model)
		allDiagnostics.Append(diagnostics...)
	}

	return ctx, model, diagnostics
}
