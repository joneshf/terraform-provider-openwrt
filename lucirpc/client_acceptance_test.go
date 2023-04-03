//go:build acceptance.test

package lucirpc_test

import (
	"context"
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
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)

		// When
		got, err := client.CreateSection(
			ctx,
			"",
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.NilError(t, err)
		assert.Check(t, !got)
	})

	t.Run("returns true when successful", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)

		// When
		got, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{},
		)

		// Then
		assert.NilError(t, err)
		assert.Check(t, got)
	})

	t.Run("creates the section", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)

		// When
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{
				"option_1": lucirpc.Boolean(true),
				"option_2": lucirpc.Integer(31),
				"option_3": lucirpc.ListString([]string{"foo", "bar", "baz"}),
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
		assert.DeepEqual(t, got, lucirpc.Options{
			".anonymous": lucirpc.Boolean(false),
			".name":      lucirpc.String("testing"),
			".type":      lucirpc.String("interface"),
			"option_1":   lucirpc.Boolean(true),
			"option_2":   lucirpc.Integer(31),
			"option_3":   lucirpc.ListString([]string{"foo", "bar", "baz"}),
		})
	})

	t.Run("does not leave pending changes when successful", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)

		// When
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{},
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

func TestClientDeleteSectionAcceptance(t *testing.T) {
	t.Parallel()

	t.Run("returns false when unsuccessful", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)

		// When
		got, err := client.DeleteSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.NilError(t, err)
		assert.Check(t, !got)
	})

	t.Run("returns true when successful", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{},
		)
		assert.NilError(t, err)

		// When
		got, err := client.DeleteSection(
			ctx,
			"network",
			"testing",
		)

		// Then
		assert.NilError(t, err)
		assert.Check(t, got)
	})

	t.Run("deletes the section", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{},
		)
		assert.NilError(t, err)

		// When
		_, err = client.DeleteSection(
			ctx,
			"network",
			"testing",
		)

		// Then
		assert.NilError(t, err)
		_, err = client.GetSection(
			ctx,
			"network",
			"testing",
		)
		assert.ErrorContains(t, err, "")
	})

	t.Run("does not leave pending changes when successful", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{},
		)
		assert.NilError(t, err)

		// When
		_, err = client.DeleteSection(
			ctx,
			"network",
			"testing",
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
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)

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
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)

		// When
		got, err := client.GetSection(
			ctx,
			"system",
			"@system[0]",
		)

		// Then
		assert.NilError(t, err)
		want := lucirpc.Options{
			".anonymous":   lucirpc.Boolean(true),
			".name":        lucirpc.String("cfg01e48a"),
			".type":        lucirpc.String("system"),
			"hostname":     lucirpc.String("OpenWrt"),
			"log_size":     lucirpc.Integer(64),
			"timezone":     lucirpc.String("UTC"),
			"ttylogin":     lucirpc.Boolean(false),
			"urandom_seed": lucirpc.Boolean(false),
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

func TestNewClientAcceptance(t *testing.T) {
	t.Parallel()

	t.Run("does not error when authentication succeeds", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(ctx, *dockerPool, t)

		// When
		_, err := lucirpc.NewClient(
			ctx,
			openWrtServer.Scheme,
			openWrtServer.Hostname,
			openWrtServer.HTTPPort,
			openWrtServer.Username,
			openWrtServer.Password,
		)

		// Then
		assert.NilError(t, err)
	})
}

func TestClientUpdateSectionAcceptance(t *testing.T) {
	t.Parallel()

	t.Run("returns false when unsuccessful", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)

		// When
		got, err := client.UpdateSection(
			ctx,
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.NilError(t, err)
		assert.Check(t, !got)
	})

	t.Run("fails with no options", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{},
		)
		assert.NilError(t, err)

		// When
		got, err := client.UpdateSection(
			ctx,
			"network",
			"testing",
			lucirpc.Options{},
		)

		// Then
		assert.NilError(t, err)
		assert.Check(t, !got)
	})

	t.Run("returns true when successful", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{},
		)
		assert.NilError(t, err)

		// When
		got, err := client.UpdateSection(
			ctx,
			"network",
			"testing",
			lucirpc.Options{
				"foo": lucirpc.Boolean(true),
			},
		)

		// Then
		assert.NilError(t, err)
		assert.Check(t, got)
	})

	t.Run("updates the section", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{},
		)
		assert.NilError(t, err)

		// When
		_, err = client.UpdateSection(
			ctx,
			"network",
			"testing",
			lucirpc.Options{
				"option_1": lucirpc.Boolean(true),
				"option_2": lucirpc.Integer(31),
				"option_3": lucirpc.ListString([]string{"foo", "bar", "baz"}),
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
		assert.DeepEqual(t, got, lucirpc.Options{
			".anonymous": lucirpc.Boolean(false),
			".name":      lucirpc.String("testing"),
			".type":      lucirpc.String("interface"),
			"option_1":   lucirpc.Boolean(true),
			"option_2":   lucirpc.Integer(31),
			"option_3":   lucirpc.ListString([]string{"foo", "bar", "baz"}),
		})
	})

	t.Run("does not leave pending changes when successful", func(t *testing.T) {
		t.Parallel()

		// Given
		ctx := context.Background()
		openWrtServer := acceptancetest.RunOpenWrtServer(
			ctx,
			*dockerPool,
			t,
		)
		client := openWrtServer.LuCIRPCClient(
			ctx,
			t,
		)
		_, err := client.CreateSection(
			ctx,
			"network",
			"interface",
			"testing",
			lucirpc.Options{},
		)
		assert.NilError(t, err)

		// When
		_, err = client.UpdateSection(
			ctx,
			"network",
			"testing",
			lucirpc.Options{
				"option_1": lucirpc.Boolean(true),
				"option_2": lucirpc.Integer(31),
				"option_3": lucirpc.ListString([]string{"foo", "bar", "baz"}),
			},
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
