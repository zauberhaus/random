// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package main

import (
	"fmt"
	"log"

	"github.com/zauberhaus/random/pkg/stringer"
)

type SensitiveConfig struct {
	APIKey string
	Port   int
}

func main() {
	config := SensitiveConfig{
		APIKey: "secret-key-12345",
		Port:   8080,
	}

	// Default string conversion (JSON-like for structs)
	str, err := stringer.String(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Default:", str)

	// Define a hook to redact the APIKey
	redactHook := stringer.NewStringHookFor[SensitiveConfig](func(v any) (string, error) {
		c := v.(SensitiveConfig)
		return fmt.Sprintf(`{"APIKey":"REDACTED","Port":%d}`, c.Port), nil
	})

	// Use the hook
	redactedStr, err := stringer.String(config, redactHook)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Redacted:", redactedStr)
}
