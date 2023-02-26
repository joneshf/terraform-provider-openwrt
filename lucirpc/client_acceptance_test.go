//go:build acceptance.test

package lucirpc_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

const (
	acceptanceTestHost   = "localhost"
	acceptanceTestPort   = uint16(8080)
	acceptanceTestScheme = "http"
)

func TestClientGetSectionAcceptance(t *testing.T) {
	t.Run("returns error when get section fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		client := authenticatedClientAcceptance(
			t,
			ctx,
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
		// Given
		ctx := context.Background()
		client := authenticatedClientAcceptance(
			t,
			ctx,
		)

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

func TestNewClientAcceptance(t *testing.T) {
	t.Run("does not error when authentication succeeds", func(t *testing.T) {
		// Given
		ctx := context.Background()

		// When
		_, err := lucirpc.NewClient(
			ctx,
			acceptanceTestScheme,
			acceptanceTestHost,
			acceptanceTestPort,
			"root",
			"",
		)

		// Then
		assert.NilError(t, err)
	})
}

func authenticatedClientAcceptance(
	t *testing.T,
	ctx context.Context,
) *lucirpc.Client {
	t.Helper()
	client, err := lucirpc.NewClient(
		ctx,
		acceptanceTestScheme,
		acceptanceTestHost,
		acceptanceTestPort,
		"root",
		"",
	)
	if err != nil {
		assert.NilError(t, err)
	}

	return client
}
