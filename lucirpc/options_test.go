package lucirpc_test

import (
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
