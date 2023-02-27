package logger

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SetFieldBool sets a bool field on the logger in the [context.Context].
func SetFieldBool(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	key string,
	value interface{ ValueBool() bool },
) context.Context {
	ctx = tflog.SetField(ctx, fmt.Sprintf("%s_%s_%s", fullTypeName, terraformType, key), value.ValueBool())
	return ctx
}

// SetFieldInt64 sets an int64 field on the logger in the [context.Context].
func SetFieldInt64(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	key string,
	value interface{ ValueInt64() int64 },
) context.Context {
	ctx = tflog.SetField(ctx, fmt.Sprintf("%s_%s_%s", fullTypeName, terraformType, key), value.ValueInt64())
	return ctx
}

// SetFieldString sets a string field on the logger in the [context.Context].
func SetFieldString(
	ctx context.Context,
	fullTypeName string,
	terraformType string,
	key string,
	value interface{ ValueString() string },
) context.Context {
	ctx = tflog.SetField(ctx, fmt.Sprintf("%s_%s_%s", fullTypeName, terraformType, key), value.ValueString())
	return ctx
}
