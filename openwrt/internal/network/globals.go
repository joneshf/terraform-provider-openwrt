package network

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/logger"
	"github.com/joneshf/terraform-provider-openwrt/openwrt/internal/lucirpcglue"
)

const (
	globalsIdAttribute  = "id"
	globalsIdUCISection = ".name"

	globalsPacketSteeringAttribute = "packet_steering"
	globalsPacketSteeringUCIOption = "packet_steering"

	globalsSchemaDescription = "Contains interface-independent options affecting the network configuration in general."

	globalsTypeName       = "network_globals"
	globalsUCIConfig      = "network"
	globalsUCIGlobalsType = "globals"

	globalsULAPrefixAttribute = "ula_prefix"
	globalsULAPrefixUCIOption = "ula_prefix"
)

var (
	globalsIdSchemaAttribute = lucirpcglue.StringSchemaAttribute[globalsModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		DataSourceExistence: lucirpcglue.Required,
		Description:         "Name of the section.",
		ReadResponse: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			section map[string]json.RawMessage,
			model globalsModel,
		) (context.Context, globalsModel, diag.Diagnostics) {
			ctx, value, diagnostics := lucirpcglue.GetMetadataString(ctx, fullTypeName, terraformType, section, globalsIdUCISection)
			model.Id = value
			return ctx, model, diagnostics
		},
		ResourceExistence: lucirpcglue.Required,
		UpsertRequest: func(
			ctx context.Context,
			fullTypeName string,
			terraformType string,
			options map[string]json.RawMessage,
			model globalsModel,
		) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
			ctx = logger.SetFieldString(ctx, fullTypeName, lucirpcglue.ResourceTerraformType, globalsIdAttribute, model.Id)
			return ctx, options, diag.Diagnostics{}
		},
	}

	globalsSchemaAttributes = map[string]lucirpcglue.SchemaAttribute[globalsModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		globalsIdAttribute:             globalsIdSchemaAttribute,
		globalsULAPrefixAttribute:      globalsULAPrefixSchemaAttribute,
		globalsPacketSteeringAttribute: globalsPacketSteeringSchemaAttribute,
	}

	globalsULAPrefixSchemaAttribute = lucirpcglue.StringSchemaAttribute[globalsModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "IPv6 ULA prefix for this device.",
		ReadResponse:      lucirpcglue.ReadResponseOptionString(globalsModelSetULAPrefix, globalsULAPrefixAttribute, globalsULAPrefixUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionString(globalsModelGetULAPrefix, globalsULAPrefixAttribute, globalsULAPrefixUCIOption),
	}

	globalsPacketSteeringSchemaAttribute = lucirpcglue.BoolSchemaAttribute[globalsModel, map[string]json.RawMessage, map[string]json.RawMessage]{
		Description:       "Use every CPU to handle packet traffic.",
		ReadResponse:      lucirpcglue.ReadResponseOptionBool(globalsModelSetPacketSteering, globalsPacketSteeringAttribute, globalsPacketSteeringUCIOption),
		ResourceExistence: lucirpcglue.NoValidation,
		UpsertRequest:     lucirpcglue.UpsertRequestOptionBool(globalsModelGetPacketSteering, globalsPacketSteeringAttribute, globalsPacketSteeringUCIOption),
	}
)

type globalsModel struct {
	Id             types.String `tfsdk:"id"`
	PacketSteering types.Bool   `tfsdk:"packet_steering"`
	ULAPrefix      types.String `tfsdk:"ula_prefix"`
}

func (m globalsModel) generateAPIBody(
	ctx context.Context,
	fullTypeName string,
) (context.Context, map[string]json.RawMessage, diag.Diagnostics) {
	tflog.Info(ctx, "Generating API request body")
	var diagnostics diag.Diagnostics
	allDiagnostics := diag.Diagnostics{}
	options := map[string]json.RawMessage{}

	tflog.Debug(ctx, "Handling attributes")
	for _, attribute := range globalsSchemaAttributes {
		ctx, options, diagnostics = attribute.Upsert(ctx, fullTypeName, lucirpcglue.ResourceTerraformType, options, m)
		allDiagnostics.Append(diagnostics...)
	}

	return ctx, options, allDiagnostics
}

func globalsModelGetPacketSteering(model globalsModel) types.Bool { return model.PacketSteering }
func globalsModelGetULAPrefix(model globalsModel) types.String    { return model.ULAPrefix }

func globalsModelSetPacketSteering(model *globalsModel, value types.Bool) {
	model.PacketSteering = value
}
func globalsModelSetULAPrefix(model *globalsModel, value types.String) { model.ULAPrefix = value }
