//go:build acceptance.test

package openwrt_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
	"github.com/ory/dockertest/v3"
	"gotest.tools/v3/assert"
)

var (
	dockerPool *dockertest.Pool
)

func TestMain(m *testing.M) {
	var (
		code     int
		err      error
		tearDown func()
	)
	ctx := context.Background()
	tearDown, dockerPool, err = acceptancetest.Setup(ctx)
	defer func() {
		tearDown()
		os.Exit(code)
	}()
	if err != nil {
		fmt.Printf("Problem setting up tests: %s", err)
		code = 1
		return
	}

	log.Printf("Running tests")
	code = m.Run()
}

func TestOpenWrtProviderConfigureConnectsWithoutError(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	openWrtProvider := openwrt.New("test", os.LookupEnv)
	schemaReq := provider.SchemaRequest{}
	schemaRes := &provider.SchemaResponse{}
	openWrtProvider.Schema(ctx, schemaReq, schemaRes)
	config := tfsdk.Config{
		Schema: schemaRes.Schema,
		Raw: tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"hostname": tftypes.String,
					"password": tftypes.String,
					"port":     tftypes.Number,
					"scheme":   tftypes.String,
					"username": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"hostname": tftypes.NewValue(tftypes.String, openWrtServer.Hostname),
				"password": tftypes.NewValue(tftypes.String, openWrtServer.Password),
				"port":     tftypes.NewValue(tftypes.Number, openWrtServer.HTTPPort),
				"scheme":   tftypes.NewValue(tftypes.String, openWrtServer.Scheme),
				"username": tftypes.NewValue(tftypes.String, openWrtServer.Username),
			},
		),
	}
	req := provider.ConfigureRequest{
		Config: config,
	}
	res := &provider.ConfigureResponse{}

	// When
	openWrtProvider.Configure(ctx, req, res)

	// Then
	assert.DeepEqual(t, res.Diagnostics, diag.Diagnostics{})
}

func TestOpenWrtProviderConfigureConnectsWithoutErrorWithEnvironmentVariables(t *testing.T) {
	// Given
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	env := map[string]string{
		"OPENWRT_HOSTNAME": openWrtServer.Hostname,
		"OPENWRT_PASSWORD": openWrtServer.Password,
		"OPENWRT_PORT":     strconv.Itoa(int(openWrtServer.HTTPPort)),
		"OPENWRT_SCHEME":   openWrtServer.Scheme,
		"OPENWRT_USERNAME": openWrtServer.Username,
	}
	lookupEnv := func(key string) (string, bool) {
		value, ok := env[key]
		return value, ok
	}
	openWrtProvider := openwrt.New("test", lookupEnv)
	schemaReq := provider.SchemaRequest{}
	schemaRes := &provider.SchemaResponse{}
	openWrtProvider.Schema(ctx, schemaReq, schemaRes)
	config := tfsdk.Config{
		Schema: schemaRes.Schema,
		Raw: tftypes.NewValue(
			tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"hostname": tftypes.String,
					"password": tftypes.String,
					"port":     tftypes.Number,
					"scheme":   tftypes.String,
					"username": tftypes.String,
				},
			},
			map[string]tftypes.Value{
				"hostname": tftypes.NewValue(tftypes.String, nil),
				"password": tftypes.NewValue(tftypes.String, nil),
				"port":     tftypes.NewValue(tftypes.Number, nil),
				"scheme":   tftypes.NewValue(tftypes.String, nil),
				"username": tftypes.NewValue(tftypes.String, nil),
			},
		),
	}
	req := provider.ConfigureRequest{
		Config: config,
	}
	res := &provider.ConfigureResponse{}

	// When
	openWrtProvider.Configure(ctx, req, res)

	// Then
	assert.DeepEqual(t, res.Diagnostics, diag.Diagnostics{})
}
