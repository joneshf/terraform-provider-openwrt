package lucirpc_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestNewClient(t *testing.T) {
	t.Run("handles server not existing", func(t *testing.T) {
		// Given
		ctx := context.Background()

		// When
		_, err := lucirpc.NewClient(
			ctx,
			"http",
			"non.existent",
			80,
			"root",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "problem sending request to login")
	})

	t.Run("expects a 200 response", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		}
		address, port, close := newServer(t, http.HandlerFunc(handle))
		defer close()

		// When
		_, err := lucirpc.NewClient(
			ctx,
			address.Scheme,
			address.Hostname(),
			uint16(port),
			"root",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "expected 200 response")
	})

	t.Run("expects a valid JSONRPC response", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "[]")
		}
		address, port, close := newServer(t, http.HandlerFunc(handle))
		defer close()

		// When
		_, err := lucirpc.NewClient(
			ctx,
			address.Scheme,
			address.Hostname(),
			uint16(port),
			"root",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "unable to process login response")
	})

	t.Run("returns error when authentication fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"error": "invalid password"
			}`)
		}
		address, port, close := newServer(t, http.HandlerFunc(handle))
		defer close()

		// When
		_, err := lucirpc.NewClient(
			ctx,
			address.Scheme,
			address.Hostname(),
			uint16(port),
			"root",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "unable to login")
	})

	t.Run("handles invalid responses", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "{}")
		}
		address, port, close := newServer(t, http.HandlerFunc(handle))
		defer close()

		// When
		_, err := lucirpc.NewClient(
			ctx,
			address.Scheme,
			address.Hostname(),
			uint16(port),
			"root",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "invalid login response")
	})

	t.Run("makes request to correct endpoint", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/cgi-bin/luci/rpc/auth" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			fmt.Fprintf(w, `{
				"result": "correct path"
			}`)
		}
		address, port, close := newServer(t, http.HandlerFunc(handle))
		defer close()

		// When
		_, err := lucirpc.NewClient(
			ctx,
			address.Scheme,
			address.Hostname(),
			uint16(port),
			"root",
			"",
		)

		// Then
		assert.NilError(t, err)
	})

	t.Run("does not error when authentication succeeds", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": "abc123"
			}`)
		}
		address, port, close := newServer(t, http.HandlerFunc(handle))
		defer close()

		// When
		_, err := lucirpc.NewClient(
			ctx,
			address.Scheme,
			address.Hostname(),
			uint16(port),
			"root",
			"",
		)

		// Then
		assert.NilError(t, err)
	})
}

func newServer(
	t *testing.T,
	handler http.Handler,
) (*url.URL, int, func()) {
	server := httptest.NewServer(handler)
	address, err := url.Parse(server.URL)
	assert.NilError(t, err)
	port, err := strconv.Atoi(address.Port())
	assert.NilError(t, err)

	return address, port, server.Close
}
