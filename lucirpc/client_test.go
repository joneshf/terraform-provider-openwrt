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

func TestClientGetSection(t *testing.T) {
	t.Run("handles server not existing", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		close()

		// When
		_, err := client.GetSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "problem sending request to uci")
	})

	t.Run("makes a request to correct endpoint", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/cgi-bin/luci/rpc/uci" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			fmt.Fprintf(w, `{
				"result": {}
			}`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.GetSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.NilError(t, err)
	})

	t.Run("expects a 200 response", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.GetSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "expected 200 response")
	})

	t.Run("expects a valid JSON-RPC response", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `[]`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.GetSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "unable to process uci response")
	})

	t.Run("returns error when get section fails", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"error": ""
			}`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.GetSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "unable to get section")
	})

	t.Run("handles invalid response", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "{}")
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.GetSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "invalid uci response")
	})

	t.Run("returns section data when successful", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": {
					".name": "section-name",
					"baz": "1",
					"foo": "bar"
				}
			}`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		got, err := client.GetSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.NilError(t, err)
		want := map[string]string{
			".name": "section-name",
			"baz":   "1",
			"foo":   "bar",
		}
		assert.DeepEqual(t, got, want)
	})
}

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

	t.Run("expects a valid JSON-RPC response", func(t *testing.T) {
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

func authenticatedClient(
	t *testing.T,
	ctx context.Context,
	handler http.Handler,
) (*lucirpc.Client, func()) {
	handleWithAuth := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cgi-bin/luci/rpc/auth":
			fmt.Fprintf(w, `{
					"result": "abc123"
				}`)
		default:
			handler.ServeHTTP(w, r)
		}
	}
	address, port, close := newServer(
		t,
		http.HandlerFunc(handleWithAuth),
	)
	client, err := lucirpc.NewClient(
		ctx,
		address.Scheme,
		address.Hostname(),
		uint16(port),
		"root",
		"",
	)
	assert.NilError(t, err)

	return client, close
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
