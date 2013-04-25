// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rjson_test

import (
	"fmt"
	"io"
	"launchpad.net/rjson"
	"log"
	"os"
	"strings"
)

func ExampleMarshal() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := rjson.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	// Output:
	// {"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}
}

func ExampleUnmarshal() {
	var rjsonBlob = []byte(`[
		{"Name": "Platypus", "Order": "Monotremata"},
		{"Name": "Quoll",    "Order": "Dasyuromorphia"}
	]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := rjson.Unmarshal(rjsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
	// Output:
	// [{Name:Platypus Order:Monotremata} {Name:Quoll Order:Dasyuromorphia}]
}

// This example uses a Decoder to decode a stream of distinct JSON values.
func ExampleDecoder() {
	const rjsonStream = `
		{"Name": "Ed", "Text": "Knock knock."}
		{"Name": "Sam", "Text": "Who's there?"}
		{"Name": "Ed", "Text": "Go fmt."}
		{"Name": "Sam", "Text": "Go fmt who?"}
		{"Name": "Ed", "Text": "Go fmt yourself!"}
	`
	type Message struct {
		Name, Text string
	}
	dec := rjson.NewDecoder(strings.NewReader(rjsonStream))
	for {
		var m Message
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", m.Name, m.Text)
	}
	// Output:
	// Ed: Knock knock.
	// Sam: Who's there?
	// Ed: Go fmt.
	// Sam: Go fmt who?
	// Ed: Go fmt yourself!
}

func ExampleMarshalIndent() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon", "Dark amber"},
	}
	b, err := rjson.MarshalIndent(group, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	// Output:
	// {
	// 	ID: 1
	// 	Name: "Reds"
	// 	Colors: [
	// 		"Crimson"
	// 		"Red"
	// 		"Ruby"
	// 		"Maroon"
	// 		"Dark amber"
	// 	]
	// }
}
