package main

import (
	"bytes"
	"encoding/json"
	"log"
)

func ApplyObject(src interface{}, dest interface{}) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(src); err != nil {
		log.Println(err)
		return
	}
	dec := json.NewDecoder(&buf)
	if err := dec.Decode(dest); err != nil {
		log.Println(err)
	}
}
