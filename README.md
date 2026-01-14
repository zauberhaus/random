# random

`random` is a Go library designed to facilitate the generation of random data for testing purposes. It uses reflection to recursively populate complex data structures such as structs, maps, slices, and arrays with random values.

## Installation

```bash
go get github.com/zauberhaus/random
```

## Usage

### Generating Values

The primary entry point is `RandomFor[T]`, which generates a random value of type `T`.

```go
package main

import (
	"fmt"
	"github.com/zauberhaus/random"
)

func main() {
	// Generate a random integer
	i, _ := random.RandomFor[int]()
	fmt.Println("Random Int:", i)

	// Generate a random string
	s, _ := random.RandomFor[string]()
	fmt.Println("Random String:", s)
}
```

### Working with Structs

`random` automatically populates exported fields of structs, including nested structures.

```go
type User struct {
	ID       int
	Name     string
	IsActive bool
	Tags     map[string]int
}

func main() {
	user, err := random.RandomFor[User]()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Random User: %+v\n", user)
}
```

### Custom Generators

For types that require specific generation logic (like `time.Time` or types with validation rules), you can provide custom `RandomGenerator` implementations.

The library includes built-in generators for common types:
- `random.RandomTime`: Generates `time.Time`
- `random.RandomDuration`: Generates `time.Duration`
- `random.RandomRegexp`: Generates `*regexp.Regexp`

To use them (or your own), pass them to `RandomFor`:

```go
import (
	"time"
	"github.com/zauberhaus/random"
)

func main() {
	// Generate a random time using the specific generator
	t, _ := random.RandomFor[time.Time](random.RandomTime)
}
```

### Helper Functions

The library provides helpers for selecting random items from collections.

- `RandomOfSlice(slice)`: Returns a random element from the provided slice.
- `RandomOfMap(map)`: Returns a random key from the provided map.

### String Conversion

The `String` function provides a robust way to convert arbitrary values to strings. It handles various interfaces and complex types automatically.

Supported conversions:
- `encoding.TextMarshaler`
- `yaml.Marshaler`
- `fmt.Stringer`
- Slices and Arrays (comma-separated values)
- Maps (key=value pairs, sorted by key)
- Structs (serialized to JSON)
- Basic types

#### Custom String Conversion

You can customize string conversion for specific types using `StringHook`.

```go
type MyType struct {
	Value string
}

hook := random.NewStringHookFor[MyType](func(v any) (string, error) {
	return "custom:" + v.(MyType).Value, nil
})

val := MyType{Value: "test"}
str, err := random.String(val, hook)
// str: "custom:test"
```