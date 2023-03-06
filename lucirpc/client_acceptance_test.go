//go:build acceptance.test

package lucirpc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/ory/dockertest/v3"
	"gotest.tools/v3/assert"
)

var (
	dockerPool *dockertest.Pool
)

func TestClientGetSectionAcceptance(t *testing.T) {
	t.Parallel()

	t.Run("returns error when get section fails", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrt, client := acceptancetest.AuthenticatedClient(
			ctx,
			*dockerPool,
			t,
		)
		defer openWrt.Close()

		// When
		_, err := client.GetSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, `incorrect config ("") and/or section (""): result from LuCI`)
	})

	t.Run("returns system section data", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrt, client := acceptancetest.AuthenticatedClient(
			ctx,
			*dockerPool,
			t,
		)
		defer openWrt.Close()

		// When
		got, err := client.GetSection(
			ctx,
			"system",
			"@system[0]",
		)

		// Then
		assert.NilError(t, err)
		want := map[string]json.RawMessage{
			".anonymous":   json.RawMessage("true"),
			".name":        json.RawMessage(`"cfg01e48a"`),
			".type":        json.RawMessage(`"system"`),
			"hostname":     json.RawMessage(`"OpenWrt"`),
			"log_size":     json.RawMessage(`"64"`),
			"timezone":     json.RawMessage(`"UTC"`),
			"ttylogin":     json.RawMessage(`"0"`),
			"urandom_seed": json.RawMessage(`"0"`),
		}
		assert.DeepEqual(t, got, want)
	})
}

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

func TestNewClientAcceptance(t *testing.T) {
	t.Parallel()

	t.Run("does not error when authentication succeeds", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrt, hostname, port := acceptancetest.RunOpenWrtServer(ctx, *dockerPool, t)
		defer openWrt.Close()

		// When
		_, err := lucirpc.NewClient(
			ctx,
			acceptancetest.Scheme,
			hostname,
			port,
			acceptancetest.Username,
			acceptancetest.Password,
		)

		// Then
		assert.NilError(t, err)
	})
}
