package spexs

import (
	"bytes"
	"testing"
	"unicode/utf8"
)

func chars(s string) []Token {
	a := make([]Token, 0, 100)
	for _, c := range s {
		a = append(a, Token(c))
	}
	return a
}

func seq(data string) Sequence {
	p := Sequence{}
	b := bytes.NewBufferString(data)
	p.Pat = b.Bytes()
	p.Count = utf8.RuneCount(p.Pat)
	p.Group = 0
	return p
}

func createTestReference() *Database {
	u := NewReference(10)
	u.Alphabet = chars("ACGT")

	u.AddGroup(*NewGroup("[AC]", '1', chars("AC")))
	u.AddGroup(*NewGroup("[AG]", '2', chars("AG")))
	u.AddGroup(*NewGroup("[AT]", '3', chars("AT")))
	u.AddGroup(*NewGroup("[CG]", '4', chars("CG")))
	u.AddGroup(*NewGroup("[CT]", '5', chars("CT")))
	u.AddGroup(*NewGroup("[GT]", '6', chars("GT")))

	u.AddSequence(seq("ACGTACGG"))
	u.AddSequence(seq("CAGTCCG"))
	u.AddSequence(seq("ACGGCTA"))
	u.AddSequence(seq("GGTCAACTG"))

	return u
}

func TestReferenceNext(t *testing.T) {
	u := createTestReference()

	testStr := func(idx int, str string) {
		var x Token
		var ok bool
		idx, pos := idx, uint(0)
		for _, c := range str {
			x, pos, ok = u.Next(idx, pos)
			if !ok {
				t.Errorf("string '%s' ended too early", str)
			}
			if Token(c) != x {
				t.Errorf("wrong char: str='%s' got='%v' expected='%v'", str, x, c)
			}
		}
	}

	testStr(0, "ACGTACGG")
	testStr(2, "ACGGCTA")
	testStr(1, "CAGTCCG")
	testStr(3, "GGTCAACTG")
}