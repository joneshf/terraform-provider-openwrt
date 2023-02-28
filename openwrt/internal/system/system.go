package system

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type systemModel struct {
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
