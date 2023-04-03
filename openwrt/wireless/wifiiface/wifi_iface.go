package wifiiface

import (
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
	deviceAttribute            = "device"
	deviceAttributeDescription = "Name of the physical device. This name is what the device is known as in LuCI/UCI, or the `id` field in Terraform."
	deviceUCIOption            = "device"

	encryptionMethodAttribute            = "encryption"
	encryptionMethodAttributeDescription = `Encryption method. Currently, only PSK encryption methods are supported. Must be one of: "none", "psk", "psk2", "psk2+aes", "psk2+ccmp", "psk2+tkip", "psk2+tkip+aes", "psk2+tkip+ccmp", "psk+aes", "psk+ccmp", "psk-mixed", "psk-mixed+aes", "psk-mixed+ccmp", "psk-mixed+tkip", "psk-mixed+tkip+aes", "psk-mixed+tkip+ccmp", "psk+tkip", "psk+tkip+aes", "psk+tkip+ccmp", "sae", "sae-mixed".`
	encryptionMethodNone                 = "none"
	encryptionMethodPSK                  = "psk"
	encryptionMethodPSK2                 = "psk2"
	encryptionMethodPSK2AES              = "psk2+aes"
	encryptionMethodPSK2CCMP             = "psk2+ccmp"
	encryptionMethodPSK2TKIP             = "psk2+tkip"
	encryptionMethodPSK2TKIPAES          = "psk2+tkip+aes"
	encryptionMethodPSK2TKIPCCMP         = "psk2+tkip+ccmp"
	encryptionMethodPSKAES               = "psk+aes"
	encryptionMethodPSKCCMP              = "psk+ccmp"
	encryptionMethodPSKMixed             = "psk-mixed"
	encryptionMethodPSKMixedAES          = "psk-mixed+aes"
	encryptionMethodPSKMixedCCMP         = "psk-mixed+ccmp"
	encryptionMethodPSKMixedTKIP         = "psk-mixed+tkip"
	encryptionMethodPSKMixedTKIPAES      = "psk-mixed+tkip+aes"
	encryptionMethodPSKMixedTKIPCCMP     = "psk-mixed+tkip+ccmp"
	encryptionMethodPSKTKIP              = "psk+tkip"
	encryptionMethodPSKTKIPAES           = "psk+tkip+aes"
	encryptionMethodPSKTKIPCCMP          = "psk+tkip+ccmp"
	encryptionMethodSAE                  = "sae"
	encryptionMethodSAEMixed             = "sae-mixed"
	encryptionMethodUCIOption            = "encryption"

	isolateClientsAttribute            = "isolate"
	isolateClientsAttributeDescription = "Isolate wireless clients from each other."
	isolateClientsUCIOption            = "isolate"

	keyAttribute            = "key"
	keyAttributeDescription = "The pre-shared passphrase from which the pre-shared key will be derived. The clear text key has to be 8-63 characters long."
	keyUCIOption            = "key"

	krackWorkaroundAttribute            = "wpa_disable_eapol_key_retries"
	krackWorkaroundAttributeDescription = "Enable WPA key reinstallation attack (KRACK) workaround. This should be `true` to enable KRACK workaround (you almost surely want this enabled)."
	krackWorkaroundUCIOption            = "wpa_disable_eapol_key_retries"

	modeAP                   = "ap"
	modeAttribute            = "mode"
	modeAttributeDescription = `The operation mode of the wireless network interface controller.. Currently only "ap" is supported.`
	modeUCIOption            = "mode"

	networkAttribute            = "network"
	networkAttributeDescription = "Network interface to attach the wireless network. This name is what the interface is known as in UCI, or the `id` field in Terraform."
	networkUCIOption            = "network"

	schemaDescription = "A wireless network."

	ssidAttribute            = "ssid"
	ssidAttributeDescription = "The broadcasted SSID of the wireless network. This is what actual clients will see the network as."
	ssidUCIOption            = "ssid"

	uciConfig = "wireless"
	uciType   = "wifi-iface"
)

var (
	deviceSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       deviceAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetDevice, deviceAttribute, deviceUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetDevice, deviceAttribute, deviceUCIOption),
	}

	encryptionMethodSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       encryptionMethodAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetEncryptionMethod, encryptionMethodAttribute, encryptionMethodUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetEncryptionMethod, encryptionMethodAttribute, encryptionMethodUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				encryptionMethodNone,
				encryptionMethodPSK,
				encryptionMethodPSK2,
				encryptionMethodPSK2AES,
				encryptionMethodPSK2CCMP,
				encryptionMethodPSK2TKIP,
				encryptionMethodPSK2TKIPAES,
				encryptionMethodPSK2TKIPCCMP,
				encryptionMethodPSKAES,
				encryptionMethodPSKCCMP,
				encryptionMethodPSKMixed,
				encryptionMethodPSKMixedAES,
				encryptionMethodPSKMixedCCMP,
				encryptionMethodPSKMixedTKIP,
				encryptionMethodPSKMixedTKIPAES,
				encryptionMethodPSKMixedTKIPCCMP,
				encryptionMethodPSKTKIP,
				encryptionMethodPSKTKIPAES,
				encryptionMethodPSKTKIPCCMP,
				encryptionMethodSAE,
				encryptionMethodSAEMixed,
			),
		},
	}

	isolateClientsSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       isolateClientsAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetIsolateClients, isolateClientsAttribute, isolateClientsUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetIsolateClients, isolateClientsAttribute, isolateClientsUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(modeAttribute),
				modeAP,
			),
		},
	}

	keySchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       keyAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetKey, keyAttribute, keyUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		Sensitive:         true,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetKey, keyAttribute, keyUCIOption),
		Validators: []validator.String{
			stringvalidator.AlsoRequires(path.MatchRoot(encryptionMethodAttribute)),
			stringvalidator.LengthBetween(8, 63),
		},
	}

	krackWorkaroundSchemaAttribute = lucirpcglue.BoolSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       krackWorkaroundAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(modelSetKRACKWorkaround, krackWorkaroundAttribute, krackWorkaroundUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(modelGetKRACKWorkaround, krackWorkaroundAttribute, krackWorkaroundUCIOption),
		Validators: []validator.Bool{
			lucirpcglue.RequiresAttributeEqualString(
				path.MatchRoot(modeAttribute),
				modeAP,
			),
		},
	}

	modeSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       modeAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetMode, modeAttribute, modeUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetMode, modeAttribute, modeUCIOption),
		Validators: []validator.String{
			stringvalidator.OneOf(
				modeAP,
			),
		},
	}

	networkSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       networkAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetNetwork, networkAttribute, networkUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetNetwork, networkAttribute, networkUCIOption),
	}

	schemaAttributes = map[string]lucirpcglue.SchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		deviceAttribute:           deviceSchemaAttribute,
		encryptionMethodAttribute: encryptionMethodSchemaAttribute,
		isolateClientsAttribute:   isolateClientsSchemaAttribute,
		keyAttribute:              keySchemaAttribute,
		krackWorkaroundAttribute:  krackWorkaroundSchemaAttribute,
		lucirpcglue.IdAttribute:   lucirpcglue.IdSchemaAttribute(modelGetId, modelSetId),
		modeAttribute:             modeSchemaAttribute,
		networkAttribute:          networkSchemaAttribute,
		ssidAttribute:             ssidSchemaAttribute,
	}

	ssidSchemaAttribute = lucirpcglue.StringSchemaAttribute[model, lucirpc.Options, lucirpc.Options]{
		Description:       ssidAttributeDescription,
		ReadResponse:      lucirpcglue.ReadResponseOptionString(modelSetSSID, ssidAttribute, ssidUCIOption),
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(modelGetSSID, ssidAttribute, ssidUCIOption),
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
	Device           types.String `tfsdk:"device"`
	EncryptionMethod types.String `tfsdk:"encryption"`
	Id               types.String `tfsdk:"id"`
	IsolateClients   types.Bool   `tfsdk:"isolate"`
	Key              types.String `tfsdk:"key"`
	KRACKWorkaround  types.Bool   `tfsdk:"wpa_disable_eapol_key_retries"`
	Mode             types.String `tfsdk:"mode"`
	Network          types.String `tfsdk:"network"`
	SSID             types.String `tfsdk:"ssid"`
}

func modelGetDevice(m model) types.String           { return m.Device }
func modelGetEncryptionMethod(m model) types.String { return m.EncryptionMethod }
func modelGetId(m model) types.String               { return m.Id }
func modelGetIsolateClients(m model) types.Bool     { return m.IsolateClients }
func modelGetKey(m model) types.String              { return m.Key }
func modelGetKRACKWorkaround(m model) types.Bool    { return m.KRACKWorkaround }
func modelGetMode(m model) types.String             { return m.Mode }
func modelGetNetwork(m model) types.String          { return m.Network }
func modelGetSSID(m model) types.String             { return m.SSID }

func modelSetDevice(m *model, value types.String)           { m.Device = value }
func modelSetEncryptionMethod(m *model, value types.String) { m.EncryptionMethod = value }
func modelSetId(m *model, value types.String)               { m.Id = value }
func modelSetIsolateClients(m *model, value types.Bool)     { m.IsolateClients = value }
func modelSetKey(m *model, value types.String)              { m.Key = value }
func modelSetKRACKWorkaround(m *model, value types.Bool)    { m.KRACKWorkaround = value }
func modelSetMode(m *model, value types.String)             { m.Mode = value }
func modelSetNetwork(m *model, value types.String)          { m.Network = value }
func modelSetSSID(m *model, value types.String)             { m.SSID = value }
