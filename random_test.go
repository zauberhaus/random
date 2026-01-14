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

func TestRandom(t *testing.T) {
	// Define a struct for testing struct generation
	type TestStruct struct {
		Name    string
		Age     int
		Friends []string
		Tags    map[string]float64
		private string
	}

	// Define a custom type for testing convertibility
	type MyString string
	type MyInt int
	type MyUint uint
	type MyBool bool
	type MyFloat32 float32
	type MyFloat64 float64
	type MyStruct TestStruct
	type MyMap map[string]MyString
	type MySlice []MyString
	type MyArray [2]MyString
	type MyTime time.Time
	type MyDuration time.Duration
	type MyRegexp *regexp.Regexp

	testCases := []struct {
		name      string
		typ       reflect.Type
		assertion func(t *testing.T, val any, err error)
	}{
		{
			name: "string",
			typ:  reflect.TypeFor[string](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, "", val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "custom string type",
			typ:  reflect.TypeFor[MyString](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyString(""), val)
			},
		},
		{
			name: "custom int type",
			typ:  reflect.TypeFor[MyInt](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyInt(0), val)
			},
		},
		{
			name: "custom uint type",
			typ:  reflect.TypeFor[MyUint](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyUint(0), val)
			},
		},
		{
			name: "custom bool type",
			typ:  reflect.TypeFor[MyBool](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyBool(true), val)
			},
		},
		{
			name: "custom float32 type",
			typ:  reflect.TypeFor[MyFloat32](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyFloat32(0), val)
			},
		},
		{
			name: "custom float64 type",
			typ:  reflect.TypeFor[MyFloat64](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyFloat64(0), val)
			},
		},
		{
			name: "custom struct type",
			typ:  reflect.TypeFor[MyStruct](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyStruct{}, val)
			},
		},
		{
			name: "custom map type",
			typ:  reflect.TypeFor[MyMap](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyMap{}, val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "custom slice type",
			typ:  reflect.TypeFor[MySlice](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MySlice{}, val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "custom array type",
			typ:  reflect.TypeFor[MyArray](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyArray{}, val)
				assert.Len(t, val, 2)
			},
		},
		{
			name: "custom time type",
			typ:  reflect.TypeFor[MyTime](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyTime{}, val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "custom duration type",
			typ:  reflect.TypeFor[MyDuration](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyDuration(0), val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "custom regexp type",
			typ:  reflect.TypeFor[MyRegexp](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, MyRegexp(&regexp.Regexp{}), val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "int",
			typ:  reflect.TypeFor[int](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, 0, val)
			},
		},
		{
			name: "int8",
			typ:  reflect.TypeFor[int8](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, int8(0), val)
			},
		},
		{
			name: "int16",
			typ:  reflect.TypeFor[int16](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, int16(0), val)
			},
		},
		{
			name: "int32",
			typ:  reflect.TypeFor[int32](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, int32(0), val)
			},
		},
		{
			name: "int64",
			typ:  reflect.TypeFor[int64](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, int64(0), val)
			},
		},
		{
			name: "uint",
			typ:  reflect.TypeFor[uint](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, uint(0), val)
			},
		},
		{
			name: "uint8",
			typ:  reflect.TypeFor[uint8](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, uint8(0), val)
			},
		},
		{
			name: "uint16",
			typ:  reflect.TypeFor[uint16](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, uint16(0), val)
			},
		},
		{
			name: "uint32",
			typ:  reflect.TypeFor[uint32](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, uint32(0), val)
			},
		},
		{
			name: "uint64",
			typ:  reflect.TypeFor[uint64](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, uint64(0), val)
			},
		},
		{
			name: "float32",
			typ:  reflect.TypeFor[float32](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, float32(0), val)
			},
		},
		{
			name: "float64",
			typ:  reflect.TypeFor[float64](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, float64(0), val)
			},
		},
		{
			name: "bool",
			typ:  reflect.TypeFor[bool](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, false, val)
			},
		},
		{
			name: "pointer to int",
			typ:  reflect.TypeFor[*int](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, new(int), val)
			},
		},
		{
			name: "map[string]int",
			typ:  reflect.TypeFor[map[string]int](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, map[string]int{}, val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "map[int]string",
			typ:  reflect.TypeFor[map[int]string](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, map[int]string{}, val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "slice",
			typ:  reflect.TypeFor[[]int](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, []int{}, val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "slice of int pointer",
			typ:  reflect.TypeFor[[]*int](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, []*int{}, val)
				assert.NotEmpty(t, val)
			},
		},
		{
			name: "array",
			typ:  reflect.TypeFor[[5]string](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				assert.IsType(t, [5]string{}, val)
			},
		},
		{
			name: "struct",
			typ:  reflect.TypeFor[TestStruct](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				v, ok := val.(TestStruct)
				require.True(t, ok)
				assert.NotEmpty(t, v.Name)
				assert.NotEmpty(t, v.Friends)
				assert.Empty(t, v.private)
			},
		},
		{
			name: "array of struct",
			typ:  reflect.TypeFor[[]TestStruct](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				v, ok := val.([]TestStruct)
				require.True(t, ok)
				assert.NotEmpty(t, v)
				//assert.NotEmpty(t, v.Name)
				//assert.NotEmpty(t, v.Friends)
			},
		},
		{
			name: "pointer to struct",
			typ:  reflect.TypeFor[*TestStruct](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotNil(t, val)
				v, ok := val.(*TestStruct)
				require.True(t, ok)
				assert.NotEmpty(t, v.Name)
				// Age is an int, its zero value is 0, which is not "empty" in the same way.
				assert.NotEmpty(t, v.Friends)
			},
		},
		{
			name: "regexp",
			typ:  reflect.TypeFor[*regexp.Regexp](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, val)
				assert.IsType(t, &regexp.Regexp{}, val)
			},
		},
		{
			name: "time",
			typ:  reflect.TypeFor[time.Time](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, val)
				assert.IsType(t, time.Time{}, val)
			},
		},
		{
			name: "duration",
			typ:  reflect.TypeFor[time.Duration](),
			assertion: func(t *testing.T, val any, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, val)
				assert.IsType(t, 1*time.Second, val)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := random.Random(tc.typ, random.RandomTime, random.RandomRegexp, random.RandomDuration)
			tc.assertion(t, val, err)
		})
	}
}

func TestRandomFor(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		val, err := random.RandomFor[int]()
		require.NoError(t, err)
		assert.IsType(t, 0, val)
	})

	t.Run("string", func(t *testing.T) {
		val, err := random.RandomFor[string]()
		require.NoError(t, err)
		assert.IsType(t, "", val)
		assert.NotEmpty(t, val)
	})
}

func TestRandomOfSlice(t *testing.T) {
	t.Parallel()

	slice := []int{1, 2, 3, 4, 5}
	val := random.RandomOfSlice(slice)
	assert.Contains(t, slice, val)
}

func TestRandomOfMap(t *testing.T) {
	t.Parallel()

	m := map[string]int{"a": 1, "b": 2, "c": 3}
	key := random.RandomOfMap(m)
	assert.Contains(t, m, key)
}

func TestRandom_Interface(t *testing.T) {
	t.Parallel()

	type MyInterface interface {
		Foo()
	}

	_, err := random.Random(reflect.TypeFor[MyInterface]())
	assert.Error(t, err)
	assert.EqualError(t, err, "random interface is impossible")
}
