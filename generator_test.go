// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package random_test

import (
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zauberhaus/random"
)

func TestRandomTime_Can(t *testing.T) {
	t.Parallel()
	gen := random.RandomTime

	assert.True(t, gen.Can(reflect.TypeOf(time.Time{})), "should be true for time.Time")
	assert.False(t, gen.Can(reflect.TypeOf("")), "should be false for string")
}

func TestRandomTime_Random(t *testing.T) {
	t.Parallel()
	gen := random.RandomTime

	val, err := gen.Random()
	require.NoError(t, err)

	randomVal, ok := val.(time.Time)
	require.True(t, ok, "should return a time.Time, but got %T", val)

	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC)

	assert.False(t, randomVal.Before(start), "random time should not be before the start date")
	assert.False(t, randomVal.After(end), "random time should not be after the end date")
}

func TestRandomRegexp_Can(t *testing.T) {
	t.Parallel()
	gen := random.RandomRegexp

	assert.True(t, gen.Can(reflect.TypeOf(&regexp.Regexp{})), "should be true for *regexp.Regexp")
	assert.False(t, gen.Can(reflect.TypeOf(regexp.Regexp{})), "should be false for regexp.Regexp")
	assert.False(t, gen.Can(reflect.TypeOf("")), "should be false for string")
}

func TestRandomRegexp_Random(t *testing.T) {
	t.Parallel()
	gen := random.RandomRegexp

	val, err := gen.Random()
	require.NoError(t, err)

	_, ok := val.(*regexp.Regexp)
	assert.True(t, ok, "should return a *regexp.Regexp, but got %T", val)
}

func TestRandomDuration_Can(t *testing.T) {
	t.Parallel()
	gen := random.RandomDuration

	assert.True(t, gen.Can(reflect.TypeOf(time.Duration(0))), "should be true for time.Duration")
	assert.False(t, gen.Can(reflect.TypeOf(0)), "should be true for int")
	assert.False(t, gen.Can(reflect.TypeOf("")), "should be false for string")
}

func TestRandomDuration_Random(t *testing.T) {
	t.Parallel()
	gen := random.RandomDuration

	val, err := gen.Random()
	require.NoError(t, err)

	randomVal, ok := val.(time.Duration)
	require.True(t, ok, "should return a time.Duration, but got %T", val)

	maxDuration := 24 * time.Hour
	assert.GreaterOrEqual(t, randomVal, time.Duration(0), "duration should not be negative")
	assert.Less(t, randomVal, maxDuration, "duration should be less than the max duration")
}

func TestRandomEnum(t *testing.T) {
	t.Parallel()
	gen := random.NewRandomGenerator[MyEnumer](func() (any, error) {
		val := random.RandomOfSlice(MyEnumerValues())
		return val, nil
	})

	val, err := random.RandomFor[MyEnumer](gen)
	require.NoError(t, err)
	assert.Contains(t, MyEnumerStrings(), val.String())
}

func TestRandomStructEnum(t *testing.T) {
	type myEnumStruct struct {
		Enum MyEnumer
	}

	t.Parallel()
	gen := random.NewRandomGenerator[MyEnumer](func() (any, error) {
		val := random.RandomOfSlice(MyEnumerValues())
		return val, nil
	})

	val, err := random.RandomFor[myEnumStruct](gen)
	require.NoError(t, err)
	assert.Contains(t, MyEnumerStrings(), val.Enum.String())
}

func TestRandomSliceEnum(t *testing.T) {

	t.Parallel()
	gen := random.NewRandomGenerator[MyEnumer](func() (any, error) {
		val := random.RandomOfSlice(MyEnumerValues())
		return val, nil
	})

	val, err := random.RandomFor[[]MyEnumer](gen)
	require.NoError(t, err)

	assert.NotEmpty(t, val)

	for _, v := range val {
		assert.Contains(t, MyEnumerStrings(), v.String())
	}
}

func TestRandomArrayEnum(t *testing.T) {

	t.Parallel()
	gen := random.NewRandomGenerator[MyEnumer](func() (any, error) {
		val := random.RandomOfSlice(MyEnumerValues())
		return val, nil
	})

	val, err := random.RandomFor[[5]MyEnumer](gen)
	require.NoError(t, err)

	assert.NotEmpty(t, val)
	assert.Len(t, val, 5)

	for _, v := range val {
		assert.Contains(t, MyEnumerStrings(), v.String())
	}
}

func TestRandomMapEnum(t *testing.T) {
	t.Parallel()
	gen := random.NewRandomGenerator[MyEnumer](func() (any, error) {
		val := random.RandomOfSlice(MyEnumerValues())
		return val, nil
	})

	val, err := random.RandomFor[map[string]MyEnumer](gen)
	require.NoError(t, err)

	assert.NotEmpty(t, val)

	for _, v := range val {
		assert.Contains(t, MyEnumerStrings(), v.String())
	}
}

func TestStrictRandomGenerator_Can(t *testing.T) {
	t.Parallel()

	// Define a custom type
	type MyDuration time.Duration

	// Strict generator for time.Duration
	gen := random.NewStrictRandomGenerator[time.Duration](func() (any, error) {
		return time.Duration(0), nil
	})

	assert.True(t, gen.Can(reflect.TypeOf(time.Duration(0))), "should be true for exact type")
	assert.False(t, gen.Can(reflect.TypeOf(MyDuration(0))), "should be false for convertible type in strict mode")

	// Non-strict generator
	genNonStrict := random.NewRandomGenerator[time.Duration](func() (any, error) {
		return time.Duration(0), nil
	})

	assert.True(t, genNonStrict.Can(reflect.TypeOf(time.Duration(0))), "should be true for exact type")
	assert.True(t, genNonStrict.Can(reflect.TypeOf(MyDuration(0))), "should be true for convertible type in non-strict mode")
}
