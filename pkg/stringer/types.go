package stringer

import (
	"reflect"

	"go.yaml.in/yaml/v3"
)

type StringHook struct {
	From   reflect.Type
	String func(val any) (string, error)
}

func NewStringHookFor[T any](f func(val any) (string, error)) StringHook {
	t := reflect.TypeFor[T]()
	return NewStringHook(t, f)
}

func NewStringHook(from reflect.Type, f func(val any) (string, error)) StringHook {
	return StringHook{
		From:   from,
		String: f,
	}
}

func ToYaml[T any]() StringHook {
	return StringHook{
		From: reflect.TypeFor[T](),
		String: func(val any) (string, error) {
			data, err := yaml.Marshal(val)
			if err != nil {
				return "", err
			}

			return string(data), nil
		},
	}
}
