package utils

import (
	"bytes"
	"log"

	"launchpad.net/rjson"
)

func ApplyObject(src interface{}, dest interface{}) {
	var buf bytes.Buffer
	enc := rjson.NewEncoder(&buf)
	if err := enc.Encode(src); err != nil {
		log.Println(err)
		return
	}
	dec := rjson.NewDecoder(&buf)
	if err := dec.Decode(dest); err != nil {
		log.Println(err)
	}
}
