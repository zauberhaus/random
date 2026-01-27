// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/zauberhaus/random"
)

type HexCode string

// Define a custom generator for HexCode
var RandomHex = random.NewRandomGenerator[HexCode](func() (any, error) {
	// In a real scenario, you'd generate a random hex string here.
	// For this example, we return a fixed format with random-like content.
	return HexCode("#FF5733"), nil
})

type Palette struct {
	Primary   HexCode
	Secondary HexCode
	CreatedAt time.Time
}

func main() {
	// Generate a Palette using RandomHex for HexCode fields
	// and random.RandomTime for time.Time fields.
	palette, err := random.RandomFor[Palette](RandomHex)
	if err != nil {
		log.Fatalf("Failed to generate palette: %v", err)
	}

	fmt.Printf("Palette Created At: %s\n", palette.CreatedAt)
	fmt.Printf("Primary Color: %s\n", palette.Primary)
	fmt.Printf("Secondary Color: %s\n", palette.Secondary)
}
