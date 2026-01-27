// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package stringer_test

import (
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zauberhaus/random/pkg/stringer"
)

// --- Helper types for testing interfaces ---

func Ptr[T any](v T) *T {
	return &v
}

type textMarshalerImpl struct {
	Value string
}

func (tm textMarshalerImpl) MarshalText() ([]byte, error) {
	return []byte("text:" + tm.Value), nil
}

type textMarshalerErrImpl struct{}

func (tm textMarshalerErrImpl) MarshalText() ([]byte, error) {
	return nil, errors.New("text marshal error")
}

type yamlMarshalerImpl struct {
	Value string
}

func (ym yamlMarshalerImpl) MarshalYAML() (interface{}, error) {
	return "yaml:" + ym.Value, nil
}

type yamlMarshalerErrImpl struct{}

func (ym yamlMarshalerErrImpl) MarshalYAML() (interface{}, error) {
	return nil, errors.New("yaml marshal error")
}

type stringerImpl struct {
	Value string
}

func (s stringerImpl) String() string {
	return "stringer:" + s.Value
}

type jsonStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type yamlMarshalerNonString struct{}

func (ym yamlMarshalerNonString) MarshalYAML() (interface{}, error) {
	return 12345, nil
}

func TestString_Values(t *testing.T) {
	t.Parallel()

	tmValue := textMarshalerImpl{Value: "hello"}
	ymValue := yamlMarshalerImpl{Value: "world"}
	sValue := stringerImpl{Value: "go"}

	loc := time.FixedZone("Test", int((1 * time.Hour).Seconds()))

	tests := []struct {
		name          string
		input         any
		expected      string
		expectErr     bool
		expectedErrAs error
	}{
		// Basic types
		{"string", "hello", "hello", false, nil},
		{"pointer to string", Ptr("world"), "world", false, nil},
		{"int", 123, "123", false, nil},
		{"pointer to int", Ptr(123), "123", false, nil},
		{"bool", true, "true", false, nil},
		{"pointer to bool", Ptr(true), "true", false, nil},
		{"float", 3.14, "3.14", false, nil},
		{"pointer to float", Ptr(0.123), "0.123", false, nil},

		// Nil values
		{"nil interface", nil, "", false, nil},
		{"nil pointer to string", (*string)(nil), "", false, nil},
		{"nil slice", []int(nil), "", false, nil},
		{"nil map", map[string]int(nil), "", false, nil},

		// Interface implementations (value receiver)
		{"TextMarshaler value", tmValue, "text:hello", false, nil},
		{"YamlMarshaler value", ymValue, "yaml:world", false, nil},
		{"Stringer value", sValue, "stringer:go", false, nil},

		// Interface implementations (pointer receiver)
		{"TextMarshaler pointer", &tmValue, "text:hello", false, nil},
		{"YamlMarshaler pointer", &ymValue, "yaml:world", false, nil},
		{"Stringer pointer", &sValue, "stringer:go", false, nil},

		// Interface error cases
		{"TextMarshaler error", textMarshalerErrImpl{}, "", true, errors.New("text marshal error")},
		{"TextMarshaler pointer error", &textMarshalerErrImpl{}, "", true, errors.New("text marshal error")},
		{"YamlMarshaler error", yamlMarshalerErrImpl{}, "", true, errors.New("yaml marshal error")},
		{"YamlMarshaler pointer error", &yamlMarshalerErrImpl{}, "", true, errors.New("yaml marshal error")},

		// Slices and Arrays
		{"slice of ints", []int{1, 2, 3}, "1,2,3", false, nil},
		{"slice of strings", []string{"a", "b", "c"}, "a,b,c", false, nil},
		{"array of ints", [3]int{4, 5, 6}, "4,5,6", false, nil},
		{"slice with interface value", []any{1, "two", true}, "1,two,true", false, nil},
		{"slice with TextMarshaler", []textMarshalerImpl{tmValue, {Value: "another"}}, "text:hello,text:another", false, nil},
		{"slice of struct", []jsonStruct{{Name: "Gemini", Age: 1}, {Name: "Code Assist", Age: 2}}, "{\"name\":\"Gemini\",\"age\":1},{\"name\":\"Code Assist\",\"age\":2}", false, nil},
		{"array of struct", [2]jsonStruct{{Name: "Gemini", Age: 1}, {Name: "Code Assist", Age: 2}}, "{\"name\":\"Gemini\",\"age\":1},{\"name\":\"Code Assist\",\"age\":2}", false, nil},

		// Maps
		{"map string to int", map[string]int{"one": 1, "two": 2}, "one=1,two=2", false, nil}, // Note: map order is not guaranteed
		{"map int to string", map[int]string{1: "one", 2: "two"}, "1=one,2=two", false, nil},
		{"map with Stringer key", map[stringerImpl]int{{Value: "key"}: 1}, "stringer:key=1", false, nil},

		// Structs (JSON marshaling)
		{"simple struct", jsonStruct{Name: "Gemini", Age: 1}, `{"name":"Gemini","age":1}`, false, nil},
		{"pointer to struct", &jsonStruct{Name: "Code Assist", Age: 2}, `{"name":"Code Assist","age":2}`, false, nil},

		// Other complex types
		{"time.Time", time.Date(2023, 10, 27, 10, 0, 0, 0, loc), "2023-10-27T10:00:00+01:00", false, nil},
	}

	// Helper to sort map string output for consistent comparison
	sortMapString := func(s string) string {
		parts := strings.Split(s, ",")
		slices.Sort(parts)
		return strings.Join(parts, ",")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := stringer.String(tt.input)

			if tt.expectErr {
				require.Error(t, err)
				if tt.expectedErrAs != nil {
					// The target for ErrorAs must be a pointer to a variable of an interface or error type.
					// We create a new variable of the expected error's type and pass its address.
					target := reflect.New(reflect.TypeOf(tt.expectedErrAs).Elem()).Interface()
					assert.ErrorAs(t, err, &target)
					assert.EqualError(t, err, tt.expectedErrAs.Error())
				}
			} else {
				require.NoError(t, err)
				expected := tt.expected
				// Sort map output for comparison
				if strings.Contains(tt.name, "map") {
					result = sortMapString(result)
					expected = sortMapString(expected)
				}
				assert.Equal(t, expected, result)
			}
		})
	}
}

func TestString_NonStringerYAML(t *testing.T) {
	// Special case for yaml.Marshaler that returns a non-string type
	result, err := stringer.String(yamlMarshalerNonString{}) // This type is now defined at the package level
	require.NoError(t, err)
	assert.Equal(t, "12345", result)
}

func TestString_JSONMarshalError(t *testing.T) {
	// Special case for json.Marshal that returns an error
	type jsonErrorStruct struct {
		C chan int
	}

	_, err := stringer.String(jsonErrorStruct{C: make(chan int)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "json: unsupported type: chan int")
}

func TestString_SliceWithErrorElement(t *testing.T) {
	// Special case for a slice containing an element that errors on string conversion
	slice := []any{1, "two", textMarshalerErrImpl{}}
	_, err := stringer.String(slice)
	require.Error(t, err)
	assert.ErrorContains(t, err, "text marshal error")
}

func TestString_MapWithErrorElement(t *testing.T) {
	// Special case for a map containing an element that errors on string conversion
	m := map[string]any{"good": 1, "bad": textMarshalerErrImpl{}}
	_, err := stringer.String(m)
	require.Error(t, err)
	assert.ErrorContains(t, err, "text marshal error")
}

func TestString_MapWithErrorKey(t *testing.T) {
	// Special case for a map containing a key that errors on string conversion
	m := map[any]int{textMarshalerErrImpl{}: 1}
	_, err := stringer.String(m)
	require.Error(t, err)
	assert.ErrorContains(t, err, "text marshal error")
}

func TestString_WithHook(t *testing.T) {
	type customStringHookType struct {
		Data string
	}

	f := func(val any) (string, error) {
		if v, ok := val.(customStringHookType); ok {
			return "hooked:" + v.Data, nil
		}
		if v, ok := val.(*customStringHookType); ok {
			return "hooked_ptr:" + v.Data, nil
		}
		return "", errors.New("unexpected type for customStringer")
	}

	valueHook := stringer.NewStringHook(reflect.TypeOf(customStringHookType{}), f)
	pointerHook := stringer.NewStringHook(reflect.TypeOf(&customStringHookType{}), f)

	t.Run("with value type hook", func(t *testing.T) {
		val := customStringHookType{Data: "test"}
		result, err := stringer.String(val, valueHook)
		require.NoError(t, err)
		assert.Equal(t, "hooked:test", result)
	})

	t.Run("with pointer type hook", func(t *testing.T) {
		val := &customStringHookType{Data: "pointer_test"}
		result, err := stringer.String(val, pointerHook)
		require.NoError(t, err)
		assert.Equal(t, "hooked_ptr:pointer_test", result)
	})
}

func TestString_WithHookFor(t *testing.T) {
	type customHookType struct {
		Value string
	}

	t.Run("with value type hook", func(t *testing.T) {
		hook := stringer.NewStringHookFor[customHookType](func(v any) (string, error) {
			return "custom:" + v.(customHookType).Value, nil
		})

		val := customHookType{Value: "test"}
		str, err := stringer.String(val, hook)
		require.NoError(t, err)
		assert.Equal(t, "custom:test", str)
	})

	t.Run("with pointer type hook", func(t *testing.T) {
		hook := stringer.NewStringHookFor[*customHookType](func(v any) (string, error) {
			return "custom_ptr:" + v.(*customHookType).Value, nil
		})

		val := &customHookType{Value: "pointer_test"}
		str, err := stringer.String(val, hook)
		require.NoError(t, err)
		assert.Equal(t, "custom_ptr:pointer_test", str)
	})
}
