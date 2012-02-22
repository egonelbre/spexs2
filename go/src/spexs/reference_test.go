package spexs

import (
	"testing"
	"unicode/utf8"
	"bytes"
)

func chars(s string) []Char {
	a := make([]Char, 0, 100)
	for c := range s {
		a = append(a, Char(c))
	}
	return a
}

func pattern(data string) ReferencePattern {
	p := ReferencePattern{}
	b := bytes.NewBufferString(data)
	p.Pat = b.Bytes()
	p.Count = utf8.RuneCount(p.Pat)
	return p
}

func createUnicodeReference() *UnicodeReference {
	u := &UnicodeReference{}
	u.Alphabet = chars("ACGT")

	u.Groups = make([]Group, 6)
	u.Groups[0] = *NewGroup('1', chars("AC"))
	u.Groups[1] = *NewGroup('2', chars("AG"))
	u.Groups[2] = *NewGroup('3', chars("AT"))
	u.Groups[3] = *NewGroup('4', chars("CG"))
	u.Groups[4] = *NewGroup('5', chars("CT"))
	u.Groups[5] = *NewGroup('6', chars("GT"))

	u.Pats = make([]ReferencePattern, 4)
	u.Pats[0] = pattern("ACGTACGG")
	u.Pats[1] = pattern("CAGTCCG")
	u.Pats[2] = pattern("ACGGCTA")
	u.Pats[3] = pattern("GGTCAACTG")

	return u
}

func TestUnicodeReferenceNext(t *testing.T){
	u := createUnicodeReference()

	testStr := func (idx int, str string) {
		p := PosEncode(idx, 0)
		var x Char
		var ok bool
		for _, c := range str {
			x, p, ok = u.Next(p)
			if !ok {
				t.Errorf("string '%s' ended too early", str)
			}
			if Char(c) != x {
				t.Errorf("wrong char: str='%s' got='%v' expected='%v'", str, x, c)
			}
		}
	}
	
	testStr(0, "ACGTACGG")
	testStr(2, "ACGGCTA")
	testStr(1, "CAGTCCG")
	testStr(3, "GGTCAACTG")
}
