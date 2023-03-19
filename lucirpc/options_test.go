package lucirpc_test

import (
	"encoding/json"
	"testing"

	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestOptionsGetBoolean(t *testing.T) {
	t.Run("errors with no option", func(t *testing.T) {
		// Given
		options := lucirpc.Options{}

		// When
		_, err := options.GetBoolean("option1")

		// Then
		want := lucirpc.NewOptionNotFoundError("option1", []string{})
		assert.DeepEqual(t, err, want)
	})

	t.Run("errors with wrong type", func(t *testing.T) {
		// Given
		options := lucirpc.Options{
			"option1": lucirpc.Integer(0),
		}

		// When
		_, err := options.GetBoolean("option1")

		// Then
		want := lucirpc.NewOptionTypeMismatchError("a boolean", "an integer")
		assert.DeepEqual(t, err, want)
	})

	t.Run("returns correct option", func(t *testing.T) {
		// Given
		options := lucirpc.Options{
			"option1": lucirpc.Boolean(true),
		}

		// When
		got, err := options.GetBoolean("option1")

		// Then
		want := true
		assert.NilError(t, err)
		assert.DeepEqual(t, got, want)
	})
}

func TestOptionsGetInteger(t *testing.T) {
	t.Run("errors with no option", func(t *testing.T) {
		// Given
		options := lucirpc.Options{}

		// When
		_, err := options.GetInteger("option1")

		// Then
		want := lucirpc.NewOptionNotFoundError("option1", []string{})
		assert.DeepEqual(t, err, want)
	})

	t.Run("errors with wrong type", func(t *testing.T) {
		// Given
		options := lucirpc.Options{
			"option1": lucirpc.Boolean(false),
		}

		// When
		_, err := options.GetInteger("option1")

		// Then
		want := lucirpc.NewOptionTypeMismatchError("an integer", "a boolean")
		assert.DeepEqual(t, err, want)
	})

	t.Run("returns correct option", func(t *testing.T) {
		// Given
		options := lucirpc.Options{
			"option1": lucirpc.Integer(31),
		}

		// When
		got, err := options.GetInteger("option1")

		// Then
		want := 31
		assert.NilError(t, err)
		assert.DeepEqual(t, got, want)
	})
}

func TestOptionsGetListString(t *testing.T) {
	t.Run("errors with no option", func(t *testing.T) {
		// Given
		options := lucirpc.Options{}

		// When
		_, err := options.GetListString("option1")

		// Then
		want := lucirpc.NewOptionNotFoundError("option1", []string{})
		assert.DeepEqual(t, err, want)
	})

	t.Run("errors with wrong type", func(t *testing.T) {
		// Given
		options := lucirpc.Options{
			"option1": lucirpc.Boolean(false),
		}

		// When
		_, err := options.GetListString("option1")

		// Then
		want := lucirpc.NewOptionTypeMismatchError("a list of strings", "a boolean")
		assert.DeepEqual(t, err, want)
	})

	t.Run("returns correct option", func(t *testing.T) {
		// Given
		options := lucirpc.Options{
			"option1": lucirpc.ListString([]string{
				"value1",
				"value2",
				"value3",
			}),
		}

		// When
		got, err := options.GetListString("option1")

		// Then
		want := []string{
			"value1",
			"value2",
			"value3",
		}
		assert.NilError(t, err)
		assert.DeepEqual(t, got, want)
	})
}

func TestOptionsGetString(t *testing.T) {
	t.Run("errors with no option", func(t *testing.T) {
		// Given
		options := lucirpc.Options{}

		// When
		_, err := options.GetString("option1")

		// Then
		want := lucirpc.NewOptionNotFoundError("option1", []string{})
		assert.DeepEqual(t, err, want)
	})

	t.Run("errors with wrong type", func(t *testing.T) {
		// Given
		options := lucirpc.Options{
			"option1": lucirpc.Boolean(false),
		}

		// When
		_, err := options.GetString("option1")

		// Then
		want := lucirpc.NewOptionTypeMismatchError("a string", "a boolean")
		assert.DeepEqual(t, err, want)
	})

	t.Run("returns correct option", func(t *testing.T) {
		// Given
		options := lucirpc.Options{
			"option1": lucirpc.String("hello"),
		}

		// When
		got, err := options.GetString("option1")

		// Then
		want := "hello"
		assert.NilError(t, err)
		assert.DeepEqual(t, got, want)
	})
}

func TestOptionsMarshalJSON(t *testing.T) {
	t.Run("works for all types", func(t *testing.T) {
		// Given
		options := lucirpc.Options{
			"option1": lucirpc.Boolean(true),
			"option2": lucirpc.Boolean(false),
			"option3": lucirpc.Integer(31),
			"option4": lucirpc.String("hello"),
			"option5": lucirpc.ListString([]string{
				"value1",
				"value2",
				"value3",
			}),
		}

		// When
		got, err := json.MarshalIndent(options, "\t\t", "\t")

		// Then
		want := []byte(`{
			"option1": true,
			"option2": false,
			"option3": 31,
			"option4": "hello",
			"option5": [
				"value1",
				"value2",
				"value3"
			]
		}`)
		assert.NilError(t, err)
		assert.DeepEqual(t, got, want)
	})
}

func TestOptionsUnmarshalJSON(t *testing.T) {
	t.Run("works for all types", func(t *testing.T) {
		// Given
		var options lucirpc.Options
		rawJSON := `{
			"option1": true,
			"option2": "1",
			"option3": "yes",
			"option4": "on",
			"option5": "true",
			"option6": "enabled",
			"option7": false,
			"option8": "0",
			"option9": "no",
			"option10": "off",
			"option11": "false",
			"option12": "disabled",
			"option13": "31",
			"option14": "hello",
			"option15": [
				"value1",
				"value2",
				"value3"
			]
		}`

		// When
		err := json.Unmarshal([]byte(rawJSON), &options)

		// Then
		want := lucirpc.Options{
			"option1":  lucirpc.Boolean(true),
			"option2":  lucirpc.Boolean(true),
			"option3":  lucirpc.Boolean(true),
			"option4":  lucirpc.Boolean(true),
			"option5":  lucirpc.Boolean(true),
			"option6":  lucirpc.Boolean(true),
			"option7":  lucirpc.Boolean(false),
			"option8":  lucirpc.Boolean(false),
			"option9":  lucirpc.Boolean(false),
			"option10": lucirpc.Boolean(false),
			"option11": lucirpc.Boolean(false),
			"option12": lucirpc.Boolean(false),
			"option13": lucirpc.Integer(31),
			"option14": lucirpc.String("hello"),
			"option15": lucirpc.ListString([]string{
				"value1",
				"value2",
				"value3",
			}),
		}
		assert.NilError(t, err)
		assert.DeepEqual(t, options, want)
	})
}
