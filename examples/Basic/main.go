// Copyright 2026 Zauberhaus
// Licensed to Zauberhaus under one or more agreements.
// Zauberhaus licenses this file to you under the Apache 2.0 License.
// See the LICENSE file in the project root for more information.

package main

import (
	"fmt"
	"log"

	"github.com/zauberhaus/random"
	"github.com/zauberhaus/random/pkg/stringer"
)

type Address struct {
	Street string
	City   string
	Zip    int
}

type User struct {
	ID      int
	Name    string
	Active  bool
	Address Address
	Tags    []string
}

func main() {
	// Generate a random User struct
	// The library recursively populates fields with random values
	user, err := random.RandomFor[User]()
	if err != nil {
		log.Fatalf("Failed to generate user: %v", err)
	}

	txt, err := stringer.String(user, stringer.ToYaml[User]())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Generated User: %s\n", txt)
}
