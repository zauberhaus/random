// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package random

import (
	"fmt"
	"maps"
	"math"
	"math/rand"
	"reflect"
	"slices"
	"sync"
	"time"

	"github.com/cip8/autoname"
	utils "github.com/zauberhaus/reflect_utils"
)

var (
	mutex = sync.Mutex{}
	names = map[string]bool{}

	seed   = time.Now().UTC().UnixNano()
	random = newRandom(seed)
)

func newRandom(seed int64) *rand.Rand {
	random := rand.New(rand.New(rand.NewSource(99)))
	random.Seed(seed)
	return random
}

func RandomFor[T any](generators ...RandomGenerator) (T, error) {
	var t T
	v, err := Random(reflect.TypeOf(t), generators...)
	if err != nil {
		return *new(T), err
	}

	return v.(T), nil
}

func Random(t reflect.Type, generators ...RandomGenerator) (any, error) {
	for _, g := range generators {
		if g.Can(t) {
			value, err := g.Random()
			if err != nil {
				return nil, err
			}

			return ConvertTo(value, t)
		}
	}

	if t.Kind() == reflect.Ptr {
		value, err := Random(t.Elem(), generators...)
		if err != nil {
			return nil, err
		}

		return utils.CopyToHeap(value), nil
	}

	if t.Kind() == reflect.Interface {
		return nil, fmt.Errorf("random interface is impossible")
	}

	switch t.Kind() {
	case reflect.String:
		mutex.Lock()
		defer mutex.Unlock()

		value := autoname.Generate("-")
		for {
			if _, ok := names[value]; !ok {
				break
			}

			value = autoname.Generate("-")
		}

		names[value] = true

		return ConvertTo(value, t)

	case reflect.Bool:
		value := random.Intn(2) == 1
		return ConvertTo(value, t)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value := random.Int63()
		sgn := random.Intn(2)

		if sgn == 1 {
			value *= -1
		}

		return ConvertTo(value, t)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value := random.Uint64()

		return ConvertTo(value, t)

	case reflect.Float32, reflect.Float64:
		value := random.Float64()

		value = math.Round(value*100000) / 100000

		return ConvertTo(value, t)

	case reflect.Map:
		value := reflect.MakeMap(t)
		cnt := random.Intn(10) + 1

		for range cnt {
			k, err := Random(t.Key(), generators...)
			if err != nil {
				return nil, err
			}

			v, err := Random(t.Elem(), generators...)
			if err != nil {
				return nil, err
			}

			value.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}

		return ConvertTo(value.Interface(), t)
	case reflect.Slice:
		cnt := random.Intn(10) + 1
		value := reflect.MakeSlice(t, cnt, cnt)

		for i := range cnt {
			v, err := Random(t.Elem(), generators...)
			if err != nil {
				return nil, err
			}
			value.Index(i).Set(reflect.ValueOf(v))
		}
		return ConvertTo(value.Interface(), t)
	case reflect.Array:
		value := reflect.New(t).Elem()
		for i := 0; i < value.Len(); i++ {
			v, err := Random(t.Elem(), generators...)
			if err != nil {
				return nil, err
			}
			value.Index(i).Set(reflect.ValueOf(v))
		}

		return ConvertTo(value.Interface(), t)
	case reflect.Struct:
		value := reflect.New(t).Elem()

		for i := 0; i < value.NumField(); i++ {
			f := value.Field(i)
			ft := value.Type().Field(i)

			if ft.IsExported() && f.CanSet() {
				v, err := Random(f.Type(), generators...)
				if err != nil {
					return nil, err
				}

				f.Set(reflect.ValueOf(v))
			}
		}

		return ConvertTo(value.Interface(), t)
	}

	return nil, nil
}

func ConvertTo(value any, t reflect.Type) (any, error) {
	v := reflect.ValueOf(value)

	if v.Type() == t {
		return value, nil
	}

	if v.CanConvert(t) {
		return v.Convert(t).Interface(), nil
	}

	return nil, fmt.Errorf("can't convert %v (%v) to %v ", v, v.Type(), t)
}

func RandomOfSlice[Slice ~[]V, V any](s Slice) V {
	idx := random.Intn(len(s))
	return V(s[idx])
}

func RandomOfMap[Map ~map[K]V, K comparable, V any](m Map) K {
	keys := slices.Collect(maps.Keys(m))

	idx := random.Intn(len(keys))
	key := keys[idx]

	return K(key)
}
