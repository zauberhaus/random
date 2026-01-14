// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package random

import (
	"reflect"
	"regexp"
	"time"

	"github.com/cip8/autoname"
)

var RandomTime = NewRandomGenerator[time.Time](func() (any, error) {
	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC)

	duration := end.Sub(start)
	randomNanos := random.Int63n(duration.Nanoseconds())
	randomDuration := time.Duration(randomNanos) * time.Nanosecond

	value := start.Add(randomDuration)
	return value, nil

})

var RandomRegexp = NewRandomGenerator[*regexp.Regexp](func() (any, error) {
	txt := autoname.Generate("-")
	return regexp.Compile(txt)
})

var RandomDuration = NewStrictRandomGenerator[time.Duration](func() (any, error) {
	duration := 24 * time.Hour
	randomNanos := random.Int63n(duration.Nanoseconds())
	value := time.Duration(randomNanos) * time.Nanosecond
	return value, nil
})

type RandomGenerator interface {
	Can(t reflect.Type) bool
	Random() (any, error)
}

func NewRandomGenerator[T any](random func() (any, error)) RandomGenerator {
	return &randomGen[T]{random: random}
}

func NewStrictRandomGenerator[T any](random func() (any, error)) RandomGenerator {
	return &randomGen[T]{random: random, strict: true}
}

type randomGen[T any] struct {
	strict bool
	random func() (any, error)
}

func (r *randomGen[T]) Can(t reflect.Type) bool {
	v := reflect.New(t)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t = reflect.TypeFor[T]()

	if r.strict {
		return v.Type() == t
	} else {
		return v.CanConvert(t)
	}
}

func (r *randomGen[T]) Random() (any, error) {
	return r.random()
}
