// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rjson

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
)

// Tests of simple examples.

type example struct {
	orig    string
	compact string
	indent  string
}

var examples = []example{
	{`1`, `1`, `1`},
	{`{}`, `{}`, `{}`},
	{`[]`, `[]`, `[]`},
	{`{"":2}`, `{"":2}`, "{\n\t\"\": 2\n}"},
	{`[3]`, `[3]`, "[\n\t3\n]"},
	{`[1,2,3]`, `[1,2,3]`, "[\n\t1\n\t2\n\t3\n]"},
	{`{x:1}`, `{x:1}`, "{\n\tx: 1\n}"},
	{`{"null":1}`, `{"null":1}`, "{\n\tnull: 1\n}"},
	{`{"true":1}`, `{"true":1}`, "{\n\ttrue: 1\n}"},
	{`{"false":1}`, `{"false":1}`, "{\n\tfalse: 1\n}"},
	{ex1, ex1, ex1i},
	// rjson-specific examples:
	{"{truenull\n :\n \"falsetrue\",\n}", `{truenull:"falsetrue"}`, "{\n\ttruenull: \"falsetrue\"\n}"},
	{"{\n\ttruenull: \"falsetrue\"\n}", `{truenull:"falsetrue"}`, "{\n\ttruenull: \"falsetrue\"\n}"},
}

var ex1 = `[true,false,null,"x",1,1.5,0,-5e+2]`

var ex1i = `[
	true
	false
	null
	"x"
	1
	1.5
	0
	-5e+2
]`

func TestCompactIndent(t *testing.T) {
	var buf bytes.Buffer
	for _, tt := range examples {
		buf.Reset()
		if err := Compact(&buf, []byte(tt.orig)); err != nil {
			t.Errorf("Compact(%#q): %v", tt.orig, err)
		}
		if s := buf.String(); s != tt.compact {
			t.Errorf("Compact(%#q) = %#q, want %#q", tt.orig, s, tt.compact)
		}

		buf.Reset()
		if err := Indent(&buf, []byte(tt.orig), "", "\t"); err != nil {
			t.Errorf("Indent(%#q): %v", tt.orig, err)
		}
		if s := buf.String(); s != tt.indent {
			t.Errorf("Indent(%#q) = %#q, want %#q", tt.orig, s, tt.indent)
		}
	}
}

// Tests of a large random structure.

func TestCompactBig(t *testing.T) {
	initBig()
	var buf bytes.Buffer
	if err := Compact(&buf, jsonBig); err != nil {
		t.Fatalf("Compact: %v", err)
	}
	b := buf.Bytes()
	if bytes.Compare(b, jsonBig) != 0 {
		t.Error("Compact(jsonBig) != jsonBig")
		t.Error(diff(b, jsonBig))
		return
	}
}

func checkUnmarshalledEquality(data []byte, v interface{}) error {
	// Should marshal back to original.
	var uv interface{}
	if err := Unmarshal(data, &uv); err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}
	if !reflect.DeepEqual(uv, v) {
		return fmt.Errorf("not equal to original when unmarshalled")
	}
	return nil
}

func TestIndentBig(t *testing.T) {
	initBig()

	data, err := findProblem(jsonBigData, func(v interface{}) error {
		data, err := Marshal(v)
		if err != nil {
			return fmt.Errorf("Marshal1: %v", err)
		}
		var buf1 bytes.Buffer
		if err := Compact(&buf1, data); err != nil {
			return fmt.Errorf("Compact1: %v", err)
		}
		b_compact := buf1.Bytes()
		if err := checkUnmarshalledEquality(b_compact, v); err != nil {
			return fmt.Errorf("Compact(data): %v", err)
		}

		// Compact should be idempotent
		var buf bytes.Buffer
		if err := Compact(&buf, b_compact); err != nil {
			return fmt.Errorf("Compact2: %v", err)
		}
		if !bytes.Equal(buf.Bytes(), b_compact) {
			return fmt.Errorf("Compact(Compact(data)) != Compact(data): %v", diff(buf.Bytes(), b_compact))
		}
		buf.Reset()

		var buf2 bytes.Buffer
		if err := Indent(&buf2, data, "", "\t"); err != nil {
			return fmt.Errorf("Indent1: %v", err)
		}
		b_indent := buf2.Bytes()
		if len(b_indent) <= len(data) && v == jsonBigData {
			// jsonBig is compact (no unnecessary spaces);
			// indenting should make it bigger, at the top
			// level at least.
			return fmt.Errorf("Indent(data) did not get bigger")
		}
		if err := checkUnmarshalledEquality(b_indent, v); err != nil {
			return fmt.Errorf("Indent(data): %v", err)
		}
		// Indent should be idempotent
		buf.Reset()
		if err := Indent(&buf, b_indent, "", "\t"); err != nil {
			return fmt.Errorf("Indent2: %v", err)
		}
		if !bytes.Equal(buf.Bytes(), b_indent) {
			return fmt.Errorf("Indent(Indent(data)) != Indent(data): %v", diff(buf.Bytes(), b_indent))
		}

		buf.Reset()
		if err := Indent(&buf, b_compact, "", "\t"); err != nil {
			return fmt.Errorf("Indent3: %v", err)
		}
		if !bytes.Equal(buf.Bytes(), b_indent) {
			return fmt.Errorf("Indent(Compact(data)) != Indent(data): %v", diff(buf.Bytes(), b_indent))
		}
		// N.B. compact(indent(jsonBig)) != jsonBig because indent
		// removes unnecessary quotes, but compact does not, currently.
		return nil
	})
	if err != nil {
		t.Errorf("problem in %#v: %v", data, err)
	}
}

type indentErrorTest struct {
	in  string
	err error
}

var indentErrorTests = []indentErrorTest{
	{`{"X": "foo", "Y"}`, &SyntaxError{"invalid character '}' after object key", 17}},
	{`{"X": "foo" "Y": "bar"}`, &SyntaxError{"invalid character '\"' after object key:value pair", 13}},
}

func TestIndentErrors(t *testing.T) {
	for i, tt := range indentErrorTests {
		slice := make([]uint8, 0)
		buf := bytes.NewBuffer(slice)
		if err := Indent(buf, []uint8(tt.in), "", ""); err != nil {
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("#%d: Indent: %#v", i, err)
				continue
			}
		}
	}
}

func TestNextValueBig(t *testing.T) {
	initBig()
	var scan scanner
	item, rest, err := nextValue(jsonBig, &scan)
	if err != nil {
		t.Fatalf("nextValue: %s", err)
	}
	if len(item) != len(jsonBig) || &item[0] != &jsonBig[0] {
		t.Errorf("invalid item: %d %d", len(item), len(jsonBig))
	}
	if len(rest) != 0 {
		t.Errorf("invalid rest: %d", len(rest))
	}

	item, rest, err = nextValue(append(jsonBig, "HELLO WORLD"...), &scan)
	if err != nil {
		t.Fatalf("nextValue extra: %s", err)
	}
	if len(item) != len(jsonBig) {
		t.Errorf("invalid item: %d %d", len(item), len(jsonBig))
	}
	if string(rest) != "HELLO WORLD" {
		t.Errorf("invalid rest: %d", len(rest))
	}
}

var benchScan scanner

func BenchmarkSkipValue(b *testing.B) {
	initBig()
	for i := 0; i < b.N; i++ {
		nextValue(jsonBig, &benchScan)
	}
	b.SetBytes(int64(len(jsonBig)))
}

// diff returns a description of the difference
// between and b
func diff(a, b []byte) string {
	for i := 0; ; i++ {
		if i >= len(a) || i >= len(b) || a[i] != b[i] {
			j := i - 10
			if j < 0 {
				j = 0
			}
			return fmt.Sprintf("diverge at %d/%d: «%s» vs «%s»", i, len(a), trim(a[j:]), trim(b[j:]))
		}
	}
	return "no difference"
}

func trim(b []byte) []byte {
	if len(b) > 20 {
		return b[0:20]
	}
	return b
}

// Generate a random JSON object.

var jsonBig []byte
var jsonBigData interface{}

const (
	big   = 10000
	small = 100
)

func initBig() {
	n := big
	if testing.Short() {
		n = small
	}
	if len(jsonBig) != n {
		v := genValue(n)
		b, err := Marshal(v)
		if err != nil {
			panic(err)
		}
		jsonBigData = v
		jsonBig = b
	}
}

func genValue(n int) interface{} {
	if n > 1 {
		switch rand.Intn(2) {
		case 0:
			return genArray(n)
		case 1:
			return genMap(n)
		}
	}
	switch rand.Intn(3) {
	case 0:
		return rand.Intn(2) == 0
	case 1:
		return rand.NormFloat64()
	case 2:
		return genString(30)
	}
	panic("unreachable")
}

func genString(stddev float64) string {
	n := int(math.Abs(rand.NormFloat64()*stddev + stddev/2))
	c := make([]rune, n)
	for i := range c {
		f := math.Abs(rand.NormFloat64()*64 + 32)
		if f > 0x10ffff {
			f = 0x10ffff
		}
		c[i] = rune(f)
	}
	return string(c)
}

func genArray(n int) []interface{} {
	f := int(math.Abs(rand.NormFloat64()) * math.Min(10, float64(n/2)))
	if f > n {
		f = n
	}
	if n > 0 && f == 0 {
		f = 1
	}
	x := make([]interface{}, f)
	for i := range x {
		x[i] = genValue(((i+1)*n)/f - (i*n)/f)
	}
	return x
}

func genMap(n int) map[string]interface{} {
	f := int(math.Abs(rand.NormFloat64()) * math.Min(10, float64(n/2)))
	if f > n {
		f = n
	}
	if n > 0 && f == 0 {
		f = 1
	}
	x := make(map[string]interface{})
	for i := 0; i < f; i++ {
		x[genString(10)] = genValue(((i+1)*n)/f - (i*n)/f)
	}
	return x
}
