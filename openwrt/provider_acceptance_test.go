//go:build acceptance.test

package openwrt_test

import (
	"context"
	"fmt"
	"log"
	"os"
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
	tearDown, dockerPool, err = acceptancetest.Setup(ctx, m)
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
	t.Parallel()

	// Given
	ctx := context.Background()
	openWrt, hostname, port := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	defer openWrt.Close()
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
				"hostname": tftypes.NewValue(tftypes.String, hostname),
				"password": tftypes.NewValue(tftypes.String, acceptancetest.Password),
				"port":     tftypes.NewValue(tftypes.Number, port),
				"scheme":   tftypes.NewValue(tftypes.String, acceptancetest.Scheme),
				"username": tftypes.NewValue(tftypes.String, acceptancetest.Username),
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
