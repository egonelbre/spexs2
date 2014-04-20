// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rjson

import "bytes"

// Compact appends to dst the rjson-encoded src with
// insignificant space characters elided. Note that
// calling Compact on the result of Indent can be
// smaller still, as Compact does not current omit
// quotes from strings that can be represented as identifiers.
func Compact(dst *bytes.Buffer, src []byte) error {
	// BUG it does not currently shorten strings
	// that can be identifiers.
	return compact(dst, src, false)
}

func compact(dst *bytes.Buffer, src []byte, escape bool) error {
	origLen := dst.Len()
	var scan scanner
	scan.reset()
	start := 0
	wasComma := false
scan:
	for i, c := range src {
		if escape && (c == '<' || c == '>' || c == '&') {
			if start < i {
				dst.Write(src[start:i])
			}
			dst.WriteString(`\u00`)
			dst.WriteByte(hex[c>>4])
			dst.WriteByte(hex[c&0xF])
			start = i + 1
		}
		switch v := scan.step(&scan, int(c)); v {
		case scanError:
			break scan
		case scanEnd, scanSkipSpace, scanComma:
			if start < i {
				dst.Write(src[start:i])
			}
			start = i + 1
			// defer writing a comma until the next object starts,
			// so we can avoid producing a redundant comma
			// before closing '}' or ']'.
			if v == scanComma {
				wasComma = true
			}
		case scanBeginLiteral, scanBeginObject, scanBeginArray:
			if wasComma {
				dst.WriteByte(',')
				wasComma = false
			}
		case scanEndObject, scanEndArray:
			wasComma = false
		}
	}
	if scan.eof() == scanError {
		dst.Truncate(origLen)
		return scan.err
	}
	if start < len(src) {
		dst.Write(src[start:])
	}
	return nil
}

func newline(dst *bytes.Buffer, prefix, indent string, depth int) {
	dst.WriteByte('\n')
	dst.WriteString(prefix)
	for i := 0; i < depth; i++ {
		dst.WriteString(indent)
	}
}

// Indent appends to dst an indented form of the rjson-encoded src.
// Each element in a rjson object or array begins on a new,
// indented line beginning with prefix followed by one or more
// copies of indent according to the indentation nesting.
// The data appended to dst has no trailing newline, to make it easier
// to embed inside other formatted rjson data.
// Strings are formatted as identifiers whenever possible;
// commas are elided.
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	origLen := dst.Len()
	var scan scanner
	scan.reset()
	needIndent := false
	wasComma := false
	startString := -1 // index of start of string that's a possible identifier
	depth := 0
	for i, c := range src {
		scan.bytes++
		v := scan.step(&scan, int(c))
		if v == scanSkipSpace {
			continue
		}
		if v == scanError {
			break
		}
		if needIndent && v != scanEndObject && v != scanEndArray {
			needIndent = false
			depth++
			newline(dst, prefix, indent, depth)
		}

		if v == scanContinue {
			if startString == -1 {
				// Emit semantically uninteresting bytes
				// (in particular, punctuation in strings) unmodified.
				dst.WriteByte(c)
				continue
			}
			if c == '"' {
				// End of key string - we've found a potential identifier
				ident := src[startString+1 : i]
				if len(ident) > 0 {
					// valid identifier - write without quotes.
					dst.Write(ident)
				} else {
					// Invalid identifier - write with original quotes.
					dst.Write(src[startString : i+1])
				}
				startString = -1
				continue
			}
			if !validIdentifierChar(int(c), i-startString-1) {
				// We're inside a string that has turned out not
				// to be a valid identifier, so flush the string and
				// ignore the rest of it.
				dst.Write(src[startString : i+1])
				startString = -1
			}
			continue
		}
		if startString != -1 {
			panic("string started but did not end")
		}

		// Add spacing around real punctuation.
		switch c {
		case '{', '[':
			if wasComma {
				newline(dst, prefix, indent, depth)
				wasComma = false
			}
			// delay indent so that empty object and array are formatted as {} and [].
			needIndent = true
			dst.WriteByte(c)

		case ',', '\n':
			// delay comma so we don't print an unnecessary newline before '}' or ']'.
			wasComma = true

		case ':':
			dst.WriteByte(c)
			dst.WriteByte(' ')

		case '}', ']':
			startString = -1
			if needIndent {
				// suppress indent in empty object/array
				needIndent = false
			} else {
				depth--
				newline(dst, prefix, indent, depth)
			}
			dst.WriteByte(c)
			wasComma = false

		default:
			if wasComma {
				newline(dst, prefix, indent, depth)
				wasComma = false
			}
			n := len(scan.parseState)
			if n > 0 && scan.parseState[n-1] == parseObjectKey && c == '"' {
				// delay writing a quoted object key until we decide whether it
				// can be a valid identifier.
				startString = i
			} else {
				dst.WriteByte(c)
			}
		}
	}
	if scan.eof() == scanError {
		dst.Truncate(origLen)
		return scan.err
	}
	return nil
}

// validIdentifierRune returns whether the given rune
// is valid at the given index of an identifier.
func validIdentifierChar(c int, i int) bool {
	if i == 0 {
		return isIdentifierStart(c)
	}
	return isIdentifierInside(c)
}
