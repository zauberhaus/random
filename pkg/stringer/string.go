// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package stringer

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	utils "github.com/zauberhaus/reflect_utils"
	"go.yaml.in/yaml/v3"
)

func String(val any, stringer ...StringHook) (string, error) {
	if utils.IsNil(val) {
		return "", nil
	}

	v := reflect.ValueOf(val)
	t := reflect.TypeOf(val)

	for _, hook := range stringer {
		if hook.From == t {
			return hook.String(val)
		}
	}

	switch o := val.(type) {
	case string:
		return o, nil
	case *string:
		return *o, nil
	}

	isPointer := t.Kind() == reflect.Pointer
	if isPointer {
		if tm, ok := val.(encoding.TextMarshaler); ok {
			data, err := tm.MarshalText()
			if err != nil {
				return "", err
			}

			return string(data), nil
		} else if ym, ok := val.(yaml.Marshaler); ok {
			data, err := ym.MarshalYAML()
			if err != nil {
				return "", err
			}

			if txt, ok := data.(string); ok {
				return txt, nil
			}

			return fmt.Sprintf("%v", data), nil
		} else if str, ok := val.(fmt.Stringer); ok {
			return str.String(), nil
		}
	} else {
		ptr := utils.CopyToHeap(val)

		if tm, ok := ptr.(encoding.TextMarshaler); ok {
			data, err := tm.MarshalText()
			if err != nil {
				return "", err
			}

			return string(data), nil
		} else if ym, ok := ptr.(yaml.Marshaler); ok {
			data, err := ym.MarshalYAML()
			if err != nil {
				return "", err
			}

			if txt, ok := data.(string); ok {
				return txt, nil
			}

			return fmt.Sprintf("%v", data), nil
		} else if str, ok := ptr.(fmt.Stringer); ok {
			return str.String(), nil
		}
	}

	if isPointer {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		var result []string
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			v, err := String(element.Interface())
			if err != nil {
				return "", err
			}

			result = append(result, v)
		}

		return strings.Join(result, ","), nil
	case reflect.Map:
		keys := v.MapKeys()

		var result []string
		for _, key := range keys {
			element := v.MapIndex(key)

			k, err := String(key.Interface())
			if err != nil {
				return "", err
			}

			val, err := String(element.Interface())
			if err != nil {
				return "", err
			}

			result = append(result, fmt.Sprintf("%v=%v", k, val))
		}

		sort.Strings(result)

		return strings.Join(result, ","), nil

	case reflect.Struct:
		data, err := json.Marshal(val)
		if err != nil {
			return "", err
		}

		return string(data), nil
	default:
		return fmt.Sprintf("%v", val), nil
	}
}
