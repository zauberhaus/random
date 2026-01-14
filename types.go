package random

import "reflect"

type StringHook struct {
	From   reflect.Type
	String func(val any) (string, error)
}

func NewStringHookFor[T any](f func(val any) (string, error)) StringHook {
	t := reflect.TypeFor[T]()
	return NewStringHook(t, f)
}

func NewStringHook(to reflect.Type, f func(val any) (string, error)) StringHook {
	return StringHook{
		From:   to,
		String: f,
	}
}
