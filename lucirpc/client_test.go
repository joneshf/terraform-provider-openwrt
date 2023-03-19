package lucirpc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestClientCreateSection(t *testing.T) {
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
		_, err := client.CreateSection(
			ctx,
			"",
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "problem sending request to create section")
	})

	t.Run("makes a request to correct endpoint", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/cgi-bin/luci/rpc/uci":
				fmt.Fprintf(w, `{
					"result": true
				}`)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.CreateSection(
			ctx,
			"",
			"",
			"",
			lucirpc.Options{},
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
		_, err := client.CreateSection(
			ctx,
			"",
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "expected create section to respond with a 200")
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
		_, err := client.CreateSection(
			ctx,
			"",
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "unable to parse create section response")
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
		_, err := client.CreateSection(
			ctx,
			"",
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "unable to create section")
	})

	t.Run("does not handle unknown stuff in result", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": 31
			}`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.CreateSection(
			ctx,
			"",
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "unable to parse create section response")
	})

	t.Run("returns section data when successful", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": true
			}`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

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
		want := true
		assert.DeepEqual(t, got, want)
	})

	t.Run("commits changes", func(t *testing.T) {
		// Given
		ctx := context.Background()
		var committed bool
		handle := func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/cgi-bin/luci/rpc/uci":
				decoder := json.NewDecoder(r.Body)
				var body map[string]json.RawMessage
				err := decoder.Decode(&body)
				assert.NilError(t, err)
				method, ok := body["method"]
				assert.Check(t, ok)
				switch string(method) {
				case `"commit"`:
					committed = true
				}

				fmt.Fprintf(w, `{
					"result": true
				}`)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		client.CreateSection(
			ctx,
			"",
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.Check(t, committed)
	})
}

func TestClientDeleteSection(t *testing.T) {
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
		_, err := client.DeleteSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "problem sending request to delete section")
	})

	t.Run("makes a request to correct endpoint", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/cgi-bin/luci/rpc/uci":
				fmt.Fprintf(w, `{
					"result": true
				}`)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.DeleteSection(
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
		_, err := client.DeleteSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "expected delete section to respond with a 200")
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
		_, err := client.DeleteSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "unable to parse delete section response")
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
		_, err := client.DeleteSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "unable to delete section")
	})

	t.Run("does not handle unknown stuff in result", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": 31
			}`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.DeleteSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.ErrorContains(t, err, "unable to parse delete section response")
	})

	t.Run("returns true when successful", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": true
			}`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		got, err := client.DeleteSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.NilError(t, err)
		want := true
		assert.DeepEqual(t, got, want)
	})

	t.Run("commits changes", func(t *testing.T) {
		// Given
		ctx := context.Background()
		var committed bool
		handle := func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/cgi-bin/luci/rpc/uci":
				decoder := json.NewDecoder(r.Body)
				var body map[string]json.RawMessage
				err := decoder.Decode(&body)
				assert.NilError(t, err)
				method, ok := body["method"]
				assert.Check(t, ok)
				switch string(method) {
				case `"commit"`:
					committed = true
				}

				fmt.Fprintf(w, `{
					"result": true
				}`)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		client.DeleteSection(
			ctx,
			"",
			"",
		)

		// Then
		assert.Check(t, committed)
	})
}

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
		assert.ErrorContains(t, err, "problem sending request to get section")
	})

	t.Run("makes a request to correct endpoint", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/cgi-bin/luci/rpc/uci":
				fmt.Fprintf(w, `{
					"result": {}
				}`)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
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
		assert.ErrorContains(t, err, "expected get section to respond with a 200")
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
		assert.ErrorContains(t, err, "unable to parse get section response")
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

	t.Run("handles errors in result", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": [
					false,
					"Invalid argument"
				]
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
		assert.ErrorContains(t, err, `incorrect config ("") and/or section (""): result from LuCI`)
	})

	t.Run("does not handle unknown stuff in result", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": 31
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
		assert.ErrorContains(t, err, "unable to parse get section response")
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
		assert.ErrorContains(t, err, "could not find section")
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
		want := lucirpc.Options{
			".name": lucirpc.String("section-name"),
			"baz":   lucirpc.Boolean(true),
			"foo":   lucirpc.String("bar"),
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
			switch r.URL.Path {
			case "/cgi-bin/luci/rpc/auth":
				fmt.Fprintf(w, `{
					"result": "correct path"
				}`)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
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
		assert.ErrorContains(t, err, "expected login to respond with a 200")
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
		assert.ErrorContains(t, err, "unable to parse login response")
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

func TestClientUpdateSection(t *testing.T) {
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
		_, err := client.UpdateSection(
			ctx,
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "problem sending request to update section")
	})

	t.Run("makes a request to correct endpoint", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/cgi-bin/luci/rpc/uci":
				fmt.Fprintf(w, `{
					"result": true
				}`)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.UpdateSection(
			ctx,
			"",
			"",
			lucirpc.Options{},
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
		_, err := client.UpdateSection(
			ctx,
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "expected update section to respond with a 200")
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
		_, err := client.UpdateSection(
			ctx,
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "unable to parse update section response")
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
		_, err := client.UpdateSection(
			ctx,
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "unable to update section")
	})

	t.Run("does not handle unknown stuff in result", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": 31
			}`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		_, err := client.UpdateSection(
			ctx,
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.ErrorContains(t, err, "unable to parse update section response")
	})

	t.Run("returns true when successful", func(t *testing.T) {
		// Given
		ctx := context.Background()
		handle := func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{
				"result": true
			}`)
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		got, err := client.UpdateSection(
			ctx,
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.NilError(t, err)
		want := true
		assert.DeepEqual(t, got, want)
	})

	t.Run("commits changes", func(t *testing.T) {
		// Given
		ctx := context.Background()
		var committed bool
		handle := func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/cgi-bin/luci/rpc/uci":
				decoder := json.NewDecoder(r.Body)
				var body map[string]json.RawMessage
				err := decoder.Decode(&body)
				assert.NilError(t, err)
				method, ok := body["method"]
				assert.Check(t, ok)
				switch string(method) {
				case `"commit"`:
					committed = true
				}

				fmt.Fprintf(w, `{
					"result": true
				}`)

			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
		client, close := authenticatedClient(
			t,
			ctx,
			http.HandlerFunc(handle),
		)
		defer close()

		// When
		client.UpdateSection(
			ctx,
			"",
			"",
			lucirpc.Options{},
		)

		// Then
		assert.Check(t, committed)
	})
}

func authenticatedClient(
	t *testing.T,
	ctx context.Context,
	handler http.Handler,
) (*lucirpc.Client, func()) {
	t.Helper()
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
	if err != nil {
		close()
		assert.NilError(t, err)
	}

	return client, close
}

func newServer(
	t *testing.T,
	handler http.Handler,
) (*url.URL, int, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	address, err := url.Parse(server.URL)
	if err != nil {
		server.Close()
		assert.NilError(t, err)
	}

	port, err := strconv.Atoi(address.Port())
	if err != nil {
		server.Close()
		assert.NilError(t, err)
	}

	return address, port, server.Close
}
