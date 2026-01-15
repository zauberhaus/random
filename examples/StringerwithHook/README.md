# Stringer with Hook Example

This example demonstrates how to use the `stringer` package with custom hooks for string conversion.

## Features

- Default string conversion for structs.
- Creating a `StringHook` to customize output (e.g., redacting sensitive data).
- Applying the hook using `stringer.String`.

## Running

```bash
go run main.go
```