package lucirpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// NewOptionNotFoundError constructs a new [OptionNotFoundError].
// The `option` should be the option that wasn't found.
// The `options` should be the options that were searched.
func NewOptionNotFoundError(
	option string,
	options []string,
) OptionNotFoundError {
	return OptionNotFoundError{
		option:  option,
		options: options,
	}
}

// NewOptionTypeMismatchError constructs a new [OptionTypeMismatchError].
// It requires both the `actual` and `expected` types.
func NewOptionTypeMismatchError(
	expected string,
	actual string,
) OptionTypeMismatchError {
	return OptionTypeMismatchError{
		actual:   actual,
		expected: expected,
	}
}

// Option represents the different types of options that UCI supports.
// They can be parsed from/serialized to JSON using [json.Unmarshal]/[json.Marshal], respectively.
// They can also attempt to be serialized to the underlying Go value using the `As*` methods.
//
// To construct an [Option],
// use one of [Boolean], [Integer], [ListString], or [String].
type Option interface {
	json.Marshaler
	json.Unmarshaler
	AsBoolean() (bool, error)
	AsInteger() (int, error)
	AsListString() ([]string, error)
	AsString() (string, error)
}

// OptionNotFoundError represents an error finding the specified option.
// It also communicates which options were searched.
type OptionNotFoundError struct {
	option  string
	options []string
}

func (e OptionNotFoundError) Equal(other OptionNotFoundError) bool {
	if e.option != other.option {
		return false
	}

	slices.Sort(e.options)
	slices.Sort(other.options)
	return slices.Equal(e.options, e.options)
}

func (e OptionNotFoundError) Error() string {
	options := []string{}
	for _, option := range e.options {
		options = append(options, fmt.Sprintf("%q", option))
	}

	combinedOptions := strings.Join(options, ", ")
	return fmt.Sprintf("could not find option: %q; available options: %s", e.option, combinedOptions)
}

// OptionTypeMismatchError represents an error where the option was not the expected type.
type OptionTypeMismatchError struct {
	actual   string
	expected string
}

func (e OptionTypeMismatchError) Equal(other OptionTypeMismatchError) bool {
	return e.actual == other.actual &&
		e.expected == other.expected
}

func (e OptionTypeMismatchError) Error() string {
	return fmt.Sprintf("expected %s, but option is %s", e.expected, e.actual)
}

// Options are the actual UCI options for each section.
// The values can be booleans, integers, lists, and strings.
type Options map[string]Option

// GetBoolean attempts to find the bool for the given option.
//
// The error could either be [NewOptionNotFoundError],
// or one of the standard JSON errors.
func (os Options) GetBoolean(option string) (bool, error) {
	value, ok := os[option]
	if !ok {
		return false, NewOptionNotFoundError(option, maps.Keys(os))
	}

	return value.AsBoolean()
}

// GetInteger attempts to find the bool for the given option.
//
// The error could either be [NewOptionNotFoundError],
// or one of the standard JSON errors.
func (os Options) GetInteger(option string) (int, error) {
	value, ok := os[option]
	if !ok {
		return 0, NewOptionNotFoundError(option, maps.Keys(os))
	}

	return value.AsInteger()
}

// GetListString attempts to find the bool for the given option.
//
// The error could either be [NewOptionNotFoundError],
// or one of the standard JSON errors.
func (os Options) GetListString(option string) ([]string, error) {
	value, ok := os[option]
	if !ok {
		return nil, NewOptionNotFoundError(option, maps.Keys(os))
	}

	return value.AsListString()
}

// GetString attempts to find the bool for the given option.
//
// The error could either be [NewOptionNotFoundError],
// or one of the standard JSON errors.
func (os Options) GetString(option string) (string, error) {
	value, ok := os[option]
	if !ok {
		return "", NewOptionNotFoundError(option, maps.Keys(os))
	}

	return value.AsString()
}

func (os *Options) UnmarshalJSON(raw []byte) error {
	var options map[string]json.RawMessage
	err := json.Unmarshal(raw, &options)
	if err != nil {
		return err
	}

	*os = map[string]Option{}
	for option, rawValue := range options {
		booleanValue := new(optionBoolean)
		errBoolean := json.Unmarshal(rawValue, booleanValue)
		if errBoolean == nil {
			(*os)[option] = booleanValue
			continue
		}

		integerValue := new(optionInteger)
		errInteger := json.Unmarshal(rawValue, integerValue)
		if errInteger == nil {
			(*os)[option] = integerValue
			continue
		}

		listStringValue := new(optionListString)
		errListString := json.Unmarshal(rawValue, listStringValue)
		if errListString == nil {
			(*os)[option] = listStringValue
			continue
		}

		stringValue := new(optionString)
		errString := json.Unmarshal(rawValue, stringValue)
		if errString == nil {
			(*os)[option] = stringValue
			continue
		}

		errAll := errors.Join(errBoolean, errInteger, errListString, errString)
		err = errors.Join(err, fmt.Errorf("could not parse option %q: %w", option, errAll))
	}

	return err
}

type optionBoolean struct {
	original string
	value    bool
}

func (o *optionBoolean) AsBoolean() (bool, error) {
	return o.value, nil
}

func (o *optionBoolean) AsInteger() (int, error) {
	switch o.original {
	case "0":
		return 0, nil

	case "1":
		return 1, nil

	default:
		return 0, NewOptionTypeMismatchError("an integer", "a boolean")
	}
}

func (o *optionBoolean) AsListString() ([]string, error) {
	return nil, NewOptionTypeMismatchError("a list of strings", "a boolean")
}

func (o *optionBoolean) AsString() (string, error) {
	switch o.original {
	case "0", "no", "off", "false", "disabled":
		return o.original, nil

	case "1", "yes", "on", "true", "enabled":
		return o.original, nil

	default:
		return "", NewOptionTypeMismatchError("a string", "a boolean")
	}
}

func (o *optionBoolean) Equal(other *optionBoolean) bool {
	return o.value == other.value
}

func (o *optionBoolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

// UnmarshalJSON attempts to convert raw JSON into a boolean.
//
// Boolean options in UCI can be any number of things:
//   - True: "1", "yes", "on", "true", "enabled"
//   - False: "0", "no", "off", "false", "disabled"
//
// Boolean options are stored in UCI as a string.
// We try to parse one of these out of the raw JSON by first making sure it's a valid string.
//
// However, boolean metadata from LuCI's JSON-RPC API is returned as a JSON boolean.
// We first try to parse the value as a normal JSON boolean,
// in case it happens to be metadata.
func (o *optionBoolean) UnmarshalJSON(raw []byte) error {
	// First try to parse as a JSON boolean.
	// We could be dealing with metadata.
	var value bool
	err := json.Unmarshal(raw, &value)
	if err == nil {
		o.original = string(raw)
		o.value = value
		return nil
	}

	// If that fails,
	// Try to parse as a UCI boolean.
	var boolish string
	err = json.Unmarshal(raw, &boolish)
	if err != nil {
		return fmt.Errorf("could not convert to a string: %w", err)
	}

	switch boolish {
	case "1", "yes", "on", "true", "enabled":
		o.original = boolish
		o.value = true
		return nil

	case "0", "no", "off", "false", "disabled":
		o.original = boolish
		o.value = false
		return nil

	default:
		return fmt.Errorf(`expected one of "1", "yes", "on", "true", "enabled", "0", "no", "off", "false", or "disabled"; got: %q`, boolish)
	}
}

// Boolean constructs a new [Option].
func Boolean(value bool) *optionBoolean {
	return &optionBoolean{
		value: value,
	}
}

type optionInteger struct {
	value int
}

func (o *optionInteger) AsBoolean() (bool, error) {
	return false, NewOptionTypeMismatchError("a boolean", "an integer")
}

func (o *optionInteger) AsInteger() (int, error) {
	return o.value, nil
}

func (o *optionInteger) AsListString() ([]string, error) {
	return nil, NewOptionTypeMismatchError("a list of strings", "an integer")
}

func (o *optionInteger) AsString() (string, error) {
	return strconv.Itoa(o.value), nil
}

func (o *optionInteger) Equal(other *optionInteger) bool {
	return o.value == other.value
}

func (o *optionInteger) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

// UnmarshalJSON attempts to convert raw JSON into an integer.
//
// Integers are stored in UCI as a string.
// We try to parse one of these out of the raw JSON by first making sure it's a valid string.
func (o *optionInteger) UnmarshalJSON(raw []byte) error {
	var intish string
	err := json.Unmarshal(raw, &intish)
	if err != nil {
		return fmt.Errorf("could not convert to a string: %w", err)
	}

	value, err := strconv.Atoi(intish)
	if err != nil {
		return fmt.Errorf("unable to parse as an integer: %w", err)
	}

	o.value = value
	return nil
}

// Integer constructs a new [Option].
func Integer(value int) *optionInteger {
	return &optionInteger{
		value: value,
	}
}

type optionListString struct {
	value []string
}

func (o *optionListString) AsBoolean() (bool, error) {
	return false, NewOptionTypeMismatchError("a boolean", "a list of strings")
}

func (o *optionListString) AsInteger() (int, error) {
	return 0, NewOptionTypeMismatchError("an integer", "a list of strings")
}

func (o *optionListString) AsListString() ([]string, error) {
	return o.value, nil
}

func (o *optionListString) AsString() (string, error) {
	return "", NewOptionTypeMismatchError("a string", "a list of strings")
}

func (o *optionListString) Equal(other *optionListString) bool {
	if len(o.value) != len(other.value) {
		return false
	}

	for i, value := range o.value {
		if value != other.value[i] {
			return false
		}
	}

	return true
}

func (o *optionListString) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func (o *optionListString) UnmarshalJSON(raw []byte) error {
	return json.Unmarshal(raw, &o.value)
}

// ListString constructs a new [Option].
func ListString(value []string) *optionListString {
	return &optionListString{
		value: value,
	}
}

type optionString struct {
	value string
}

func (o *optionString) AsBoolean() (bool, error) {
	return false, NewOptionTypeMismatchError("a boolean", "a string")
}

func (o *optionString) AsInteger() (int, error) {
	return 0, NewOptionTypeMismatchError("an integer", "a string")
}

func (o *optionString) AsListString() ([]string, error) {
	return nil, NewOptionTypeMismatchError("a list of strings", "a string")
}

func (o *optionString) AsString() (string, error) {
	return o.value, nil
}

func (o *optionString) Equal(other *optionString) bool {
	return o.value == other.value
}

func (o *optionString) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func (o *optionString) UnmarshalJSON(raw []byte) error {
	return json.Unmarshal(raw, &o.value)
}

// String constructs a new [Option].
func String(value string) *optionString {
	return &optionString{
		value: value,
	}
}
