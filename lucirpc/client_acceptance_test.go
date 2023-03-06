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

func TestClientCreateSectionAcceptance(t *testing.T) {
	t.Parallel()

	t.Run("returns false when unsuccessful", func(t *testing.T) {
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
		got, err := client.CreateSection(
			ctx,
			"",
			"",
			"",
			map[string]json.RawMessage{},
		)

		// Then
		assert.NilError(t, err)
		assert.Check(t, !got)
	})

	t.Run("returns true when successful", func(t *testing.T) {
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
		got, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			map[string]json.RawMessage{},
		)

		// Then
		assert.NilError(t, err)
		assert.Check(t, got)
	})

	t.Run("creates the section", func(t *testing.T) {
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
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			map[string]json.RawMessage{
				"option_1": json.RawMessage("true"),
				"option_2": json.RawMessage("31"),
				"option_3": json.RawMessage(`["foo", "bar", "baz"]`),
			},
		)

		// Then
		assert.NilError(t, err)
		got, err := client.GetSection(
			ctx,
			"network",
			"testing",
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, got, map[string]json.RawMessage{
			".anonymous": json.RawMessage("false"),
			".name":      json.RawMessage(`"testing"`),
			".type":      json.RawMessage(`"interface"`),
			"option_1":   json.RawMessage(`"1"`),
			"option_2":   json.RawMessage(`"31"`),
			"option_3":   json.RawMessage(`["foo","bar","baz"]`),
		})
	})

	t.Run("does not leave pending changes when successful", func(t *testing.T) {
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
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			map[string]json.RawMessage{},
		)

		// Then
		assert.NilError(t, err)
		got, err := client.ShowChanges(
			ctx,
			"network",
		)
		assert.NilError(t, err)
		assert.DeepEqual(t, got, [][]string{})
	})
}

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
