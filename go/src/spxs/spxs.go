package main

import (
	"fmt"
	"bytes"
	"unicode/utf8"
	. "spexs" 
	"runtime"
)

func chars(s string) []Char {
	a := make([]Char, 0, len(s))
	for _, c := range s {
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
	u.Groups[0] = *NewGroup('M', chars("AC"))
	u.Groups[1] = *NewGroup('R', chars("AG"))
	u.Groups[2] = *NewGroup('W', chars("AT"))
	u.Groups[3] = *NewGroup('S', chars("CG"))
	u.Groups[4] = *NewGroup('Y', chars("CT"))
	u.Groups[5] = *NewGroup('K', chars("GT"))
	/*
	u.Groups[6] = *NewGroup('B', chars("CGT"))
	u.Groups[7] = *NewGroup('D', chars("AGT"))
	u.Groups[8] = *NewGroup('H', chars("ACT"))
	u.Groups[9] = *NewGroup('V', chars("ACG"))
	*/

	u.Pats = make([]ReferencePattern, 6)
	u.Pats[0] = pattern("ACGTACGG")
	u.Pats[1] = pattern("CAGTCCG")
	u.Pats[2] = pattern("ACGGCTA")
	u.Pats[3] = pattern("AATCTACTG")
	u.Pats[4] = pattern("GTCACAACTG")
	u.Pats[5] = pattern("GCCCCACTG")

	return u
}

func main() {
	runtime.GOMAXPROCS(16)

	fmt.Println("SPEXS")
	
	ref := *createUnicodeReference()
	in := NewFifoPool()
	var t Pattern = *NewFullNodeFromRef(ref)
	in.Put(t)
	out := NewFifoPool()
	extender := GroupStarExtender
	acceptable := func(p Pattern) bool { 
		n := p.(TrieNode)
		return n.Pos.Length() > 8
	}

	RunParallel(ref,in,out,extender,acceptable)

	p, ok := out.Take()
	for ok {
		n := p.(TrieNode)
		fmt.Printf("%s : %v\n", n.String(), n.Pos.Length())
		p, ok = out.Take()
	}

}