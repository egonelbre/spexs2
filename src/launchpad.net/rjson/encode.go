// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rjson implements a backwardly compatible
// version of JSON with the aim of making it easier for humans to read
// and write.
//
// The data model is exactly that of JSON's and this package implements
// all the marshalling and unmarshalling operations supported by the
// standard library's json package.
//
// The three principal differences are:
//	- Quotes may be omitted for some object key values (see below).
//	- Commas are automatically inserted when a newline
//	is encountered after a value (but not an object key).
//	- Commas are optional at the end of an array or object.
//
// The quotes around an object key may be omitted if the
// key matches the following regular expression:
//	[a-zA-Z][a-zA-Z0-9\-_]*
// This rule will be relaxed in the future to allow unicode characters,
// probably following the same rules as Go's identifiers, and probably a
// few more non-alphabetical characters.
package rjson

import (
	"bytes"
	"encoding/base64"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// Marshal returns the JSON encoding of v. Note that despite
// the fact that this is the rjson package, this function
// does returns JSON. MarshalIndent and Indent can be used
// to produce rjson-specific output.
//
// Marshal traverses the value v recursively.
// If an encountered value implements the Marshaler interface
// and is not a nil pointer, Marshal calls its MarshalJSON method
// to produce JSON.  The nil pointer exception is not strictly necessary
// but mimics a similar, necessary exception in the behavior of
// UnmarshalJSON.
//
// Otherwise, Marshal uses the following type-dependent default encodings:
//
// Boolean values encode as JSON booleans.
//
// Floating point, integer, and Number values encode as JSON numbers.
//
// String values encode as JSON strings, with each invalid UTF-8 sequence
// replaced by the encoding of the Unicode replacement character U+FFFD.
// The angle brackets "<" and ">" are escaped to "\u003c" and "\u003e"
// to keep some browsers from misinterpreting JSON output as HTML.
//
// Array and slice values encode as JSON arrays, except that
// []byte encodes as a base64-encoded string, and a nil slice
// encodes as the null JSON object.
//
// Struct values encode as JSON objects. Each exported struct field
// becomes a member of the object unless
//   - the field's tag is "-", or
//   - the field is empty and its tag specifies the "omitempty" option.
// The empty values are false, 0, any
// nil pointer or interface value, and any array, slice, map, or string of
// length zero. The object's default key string is the struct field name
// but can be specified in the struct field's tag value. The "json" key in
// the struct field's tag value is the key name, followed by an optional comma
// and options. Examples:
//
//   // Field is ignored by this package.
//   Field int `json:"-"`
//
//   // Field appears in JSON as key "myName".
//   Field int `json:"myName"`
//
//   // Field appears in JSON as key "myName" and
//   // the field is omitted from the object if its value is empty,
//   // as defined above.
//   Field int `json:"myName,omitempty"`
//
//   // Field appears in JSON as key "Field" (the default), but
//   // the field is skipped if empty.
//   // Note the leading comma.
//   Field int `json:",omitempty"`
//
// The "string" option signals that a field is stored as JSON inside a
// JSON-encoded string.  This extra level of encoding is sometimes
// used when communicating with JavaScript programs:
//
//    Int64String int64 `json:",string"`
//
// The key name will be used if it's a non-empty string consisting of
// only Unicode letters, digits, dollar signs, percent signs, hyphens,
// underscores and slashes.
//
// Anonymous struct fields are usually marshaled as if their inner exported fields
// were fields in the outer struct, subject to the usual Go visibility rules.
// An anonymous struct field with a name given in its JSON tag is treated as
// having that name instead of as anonymous.
//
// Handling of anonymous struct fields is new in Go 1.1.
// Prior to Go 1.1, anonymous struct fields were ignored. To force ignoring of
// an anonymous struct field in both current and earlier versions, give the field
// a JSON tag of "-".
//
// Map values encode as JSON objects.
// The map's key type must be string; the object keys are used directly
// as map keys.
//
// Pointer values encode as the value pointed to.
// A nil pointer encodes as the null JSON object.
//
// Interface values encode as the value contained in the interface.
// A nil interface value encodes as the null JSON object.
//
// Channel, complex, and function values cannot be encoded in JSON.
// Attempting to encode such a value causes Marshal to return
// an UnsupportedTypeError.
//
// JSON cannot represent cyclic data structures and Marshal does not
// handle them.  Passing cyclic structures to Marshal will result in
// an infinite recursion.
//
func Marshal(v interface{}) ([]byte, error) {
	e := &encodeState{}
	err := e.marshal(v)
	if err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

// MarshalIndent is like Marshal but applies Indent to format the output.
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	b, err := Marshal(v)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// HTMLEscape appends to dst the JSON-encoded src with <, >, and &
// characters inside string literals changed to \u003c, \u003e, \u0026
// so that the JSON will be safe to embed inside HTML <script> tags.
// For historical reasons, web browsers don't honor standard HTML
// escaping within <script> tags, so an alternative JSON encoding must
// be used.
func HTMLEscape(dst *bytes.Buffer, src []byte) {
	// < > & can only appear in string literals,
	// so just scan the string one byte at a time.
	start := 0
	for i, c := range src {
		if c == '<' || c == '>' || c == '&' {
			if start < i {
				dst.Write(src[start:i])
			}
			dst.WriteString(`\u00`)
			dst.WriteByte(hex[c>>4])
			dst.WriteByte(hex[c&0xF])
			start = i + 1
		}
	}
	if start < len(src) {
		dst.Write(src[start:])
	}
}

// Marshaler is the interface implemented by objects that
// can marshal themselves into valid JSON (and hence valid
// rjson)
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "json: unsupported type: " + e.Type.String()
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return "json: unsupported value: " + e.Str
}

type InvalidUTF8Error struct {
	S string
}

func (e *InvalidUTF8Error) Error() string {
	return "json: invalid UTF-8 in string: " + strconv.Quote(e.S)
}

type MarshalerError struct {
	Type reflect.Type
	Err  error
}

func (e *MarshalerError) Error() string {
	return "json: error calling MarshalJSON for type " + e.Type.String() + ": " + e.Err.Error()
}

var hex = "0123456789abcdef"

// An encodeState encodes JSON into a bytes.Buffer.
type encodeState struct {
	bytes.Buffer // accumulated output
	scratch      [64]byte
}

func (e *encodeState) marshal(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()
	e.reflectValue(reflect.ValueOf(v))
	return nil
}

func (e *encodeState) error(err error) {
	panic(err)
}

var byteSliceType = reflect.TypeOf([]byte(nil))

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func (e *encodeState) reflectValue(v reflect.Value) {
	e.reflectValueQuoted(v, false)
}

// reflectValueQuoted writes the value in v to the output.
// If quoted is true, the serialization is wrapped in a JSON string.
func (e *encodeState) reflectValueQuoted(v reflect.Value, quoted bool) {
	if !v.IsValid() {
		e.WriteString("null")
		return
	}

	m, ok := v.Interface().(Marshaler)
	if !ok {
		// T doesn't match the interface. Check against *T too.
		if v.Kind() != reflect.Ptr && v.CanAddr() {
			m, ok = v.Addr().Interface().(Marshaler)
			if ok {
				v = v.Addr()
			}
		}
	}
	if ok && (v.Kind() != reflect.Ptr || !v.IsNil()) {
		b, err := m.MarshalJSON()
		if err == nil {
			// copy JSON into buffer, checking validity.
			err = compact(&e.Buffer, b, true)
		}
		if err != nil {
			e.error(&MarshalerError{v.Type(), err})
		}
		return
	}

	writeString := (*encodeState).WriteString
	if quoted {
		writeString = (*encodeState).string
	}

	switch v.Kind() {
	case reflect.Bool:
		x := v.Bool()
		if x {
			writeString(e, "true")
		} else {
			writeString(e, "false")
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b := strconv.AppendInt(e.scratch[:0], v.Int(), 10)
		if quoted {
			writeString(e, string(b))
		} else {
			e.Write(b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		b := strconv.AppendUint(e.scratch[:0], v.Uint(), 10)
		if quoted {
			writeString(e, string(b))
		} else {
			e.Write(b)
		}
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if math.IsInf(f, 0) || math.IsNaN(f) {
			e.error(&UnsupportedValueError{v, strconv.FormatFloat(f, 'g', -1, v.Type().Bits())})
		}
		b := strconv.AppendFloat(e.scratch[:0], f, 'g', -1, v.Type().Bits())
		if quoted {
			writeString(e, string(b))
		} else {
			e.Write(b)
		}
	case reflect.String:
		if v.Type() == numberType {
			numStr := v.String()
			if numStr == "" {
				numStr = "0" // Number's zero-val
			}
			e.WriteString(numStr)
			break
		}
		if quoted {
			sb, err := Marshal(v.String())
			if err != nil {
				e.error(err)
			}
			e.string(string(sb))
		} else {
			e.string(v.String())
		}

	case reflect.Struct:
		e.WriteByte('{')
		first := true
		for _, f := range cachedTypeFields(v.Type()) {
			fv := fieldByIndex(v, f.index)
			if !fv.IsValid() || f.omitEmpty && isEmptyValue(fv) {
				continue
			}
			if first {
				first = false
			} else {
				e.WriteByte(',')
			}
			e.string(f.name)
			e.WriteByte(':')
			e.reflectValueQuoted(fv, f.quoted)
		}
		e.WriteByte('}')

	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			e.error(&UnsupportedTypeError{v.Type()})
		}
		if v.IsNil() {
			e.WriteString("null")
			break
		}
		e.WriteByte('{')
		var sv stringValues = v.MapKeys()
		sort.Sort(sv)
		for i, k := range sv {
			if i > 0 {
				e.WriteByte(',')
			}
			e.string(k.String())
			e.WriteByte(':')
			e.reflectValue(v.MapIndex(k))
		}
		e.WriteByte('}')

	case reflect.Slice:
		if v.IsNil() {
			e.WriteString("null")
			break
		}
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// Byte slices get special treatment; arrays don't.
			s := v.Bytes()
			e.WriteByte('"')
			if len(s) < 1024 {
				// for small buffers, using Encode directly is much faster.
				dst := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
				base64.StdEncoding.Encode(dst, s)
				e.Write(dst)
			} else {
				// for large buffers, avoid unnecessary extra temporary
				// buffer space.
				enc := base64.NewEncoder(base64.StdEncoding, e)
				enc.Write(s)
				enc.Close()
			}
			e.WriteByte('"')
			break
		}
		// Slices can be marshalled as nil, but otherwise are handled
		// as arrays.
		fallthrough
	case reflect.Array:
		e.WriteByte('[')
		n := v.Len()
		for i := 0; i < n; i++ {
			if i > 0 {
				e.WriteByte(',')
			}
			e.reflectValue(v.Index(i))
		}
		e.WriteByte(']')

	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			e.WriteString("null")
			return
		}
		e.reflectValue(v.Elem())

	default:
		e.error(&UnsupportedTypeError{v.Type()})
	}
	return
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}

func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}

// stringValues is a slice of reflect.Value holding *reflect.StringValue.
// It implements the methods to sort by string.
type stringValues []reflect.Value

func (sv stringValues) Len() int           { return len(sv) }
func (sv stringValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv stringValues) Less(i, j int) bool { return sv.get(i) < sv.get(j) }
func (sv stringValues) get(i int) string   { return sv[i].String() }

func (e *encodeState) string(s string) (int, error) {
	len0 := e.Len()
	e.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' {
				i++
				continue
			}
			if start < i {
				e.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				e.WriteByte('\\')
				e.WriteByte(b)
			case '\n':
				e.WriteByte('\\')
				e.WriteByte('n')
			case '\r':
				e.WriteByte('\\')
				e.WriteByte('r')
			default:
				// This encodes bytes < 0x20 except for \n and \r,
				// as well as < and >. The latter are escaped because they
				// can lead to security holes when user-controlled strings
				// are rendered into JSON and served to some browsers.
				e.WriteString(`\u00`)
				e.WriteByte(hex[b>>4])
				e.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			e.error(&InvalidUTF8Error{s})
		}
		i += size
	}
	if start < len(s) {
		e.WriteString(s[start:])
	}
	e.WriteByte('"')
	return e.Len() - len0, nil
}

// A field represents a single field found in a struct.
type field struct {
	name      string
	tag       bool
	index     []int
	typ       reflect.Type
	omitEmpty bool
	quoted    bool
}

// byName sorts field by name, breaking ties with depth,
// then breaking ties with "name came from json tag", then
// breaking ties with index sequence.
type byName []field

func (x byName) Len() int { return len(x) }

func (x byName) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byName) Less(i, j int) bool {
	if x[i].name != x[j].name {
		return x[i].name < x[j].name
	}
	if len(x[i].index) != len(x[j].index) {
		return len(x[i].index) < len(x[j].index)
	}
	if x[i].tag != x[j].tag {
		return x[i].tag
	}
	return byIndex(x).Less(i, j)
}

// byIndex sorts field by index sequence.
type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}

// typeFields returns a list of fields that JSON should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func typeFields(t reflect.Type) []field {
	// Anonymous fields to explore at the current level and the next.
	current := []field{}
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.PkgPath != "" { // unexported
					continue
				}
				tag := sf.Tag.Get("json")
				if tag == "-" {
					continue
				}
				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i
				// Record found field and index sequence.
				if name != "" || !sf.Anonymous {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					fields = append(fields, field{
						name:      name,
						tag:       tagged,
						index:     index,
						typ:       sf.Type,
						omitEmpty: opts.Contains("omitempty"),
						quoted:    opts.Contains("string"),
					})
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				ft := sf.Type
				if ft.Name() == "" {
					// Must be pointer.
					ft = ft.Elem()
				}
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	sort.Sort(byName(fields))

	// Remove fields with annihilating name collisions
	// and also fields shadowed by fields with explicit JSON tags.
	name := ""
	out := fields[:0]
	for _, f := range fields {
		if f.name != name {
			name = f.name
			out = append(out, f)
			continue
		}
		if n := len(out); n > 0 && out[n-1].name == name && (!out[n-1].tag || f.tag) {
			out = out[:n-1]
		}
	}
	fields = out

	sort.Sort(byIndex(fields))

	return fields
}

var fieldCache struct {
	sync.RWMutex
	m map[reflect.Type][]field
}

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field {
	fieldCache.RLock()
	f := fieldCache.m[t]
	fieldCache.RUnlock()
	if f != nil {
		return f
	}

	// Compute fields without lock.
	// Might duplicate effort but won't hold other computations back.
	f = typeFields(t)
	if f == nil {
		f = []field{}
	}

	fieldCache.Lock()
	if fieldCache.m == nil {
		fieldCache.m = map[reflect.Type][]field{}
	}
	fieldCache.m[t] = f
	fieldCache.Unlock()
	return f
}