package main

import (
	. "spexs"
	"bytes"
	"unicode/utf8"
	"os"
	"io"
	"bufio"
	"strings"
	"regexp"
	"fmt"
	"log"
)

func chars(s string) []Char {
	a := make([]Char, 0, len(s))
	for _, c := range s {
		a = append(a, Char(c))
	}
	return a
}

func pattern(data string) *ReferencePattern {
	p := ReferencePattern{}
	b := bytes.NewBufferString(data)
	p.Pat = b.Bytes()
	p.Count = utf8.RuneCount(p.Pat)
	return &p
}

func NewReferenceFromFile(refName string, charName string) (ref *UnicodeReference) {
	var (
		file *os.File
		reader *bufio.Reader
		line string
		err error
	)

	if file, err = os.Open(refName); err != nil {
		log.Fatal(err)
	}
	
	ref = NewUnicodeReference(1024)

	reader = bufio.NewReader(file)
	for err == nil {
		if line, err = reader.ReadString('\n'); err != nil && err != io.EOF {
			log.Fatal(err)
		}
		
		line = strings.TrimSpace(line)
		
		if len(line) == 0 {
			continue
		}

		p := pattern(line)
		if p != nil {
			ref.Pats = append(ref.Pats, *p)
		}
	}

	file.Close()
	err = nil

	if file, err = os.Open(charName); err != nil {
		log.Fatal(err)
	}

	regComment, _ := regexp.Compile("^\\s*(#.*)?$")
	regAlphabet, _ := regexp.Compile("^\\s*(\\S+)\\s*(#.*)?$")
	regGroup, _ := regexp.Compile("^\\s*(\\S+)\\s*,\\s*(\\S)\\s*,\\s*(\\S+)\\s*(#.*)?$")

	first := true
	lineNo := 0
	reader = bufio.NewReader(file)

	for err == nil {
		if line, err = reader.ReadString('\n'); err != nil && err != io.EOF {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		lineNo += 1;

		if regComment.MatchString(line) {
			continue
		}

		if first && regAlphabet.MatchString(line) {
			tokens := regAlphabet.FindStringSubmatch(line)
			ref.Alphabet = chars(tokens[1])
			first = false
			continue
		}
		
		if !first && regGroup.MatchString(line) {
			tokens := regGroup.FindStringSubmatch(line)
			id := chars(tokens[2])[0]
			g := *NewGroup(tokens[3], id, chars(tokens[1]))
			ref.Groups[id] = g
			continue
		}

		fmt.Printf("Invalid charset entry on line %v : %v\n", lineNo, line)
	}

	ref.Groupings = make([]int, 2)
	ref.Groupings[0] = len(ref.Pats)
	ref.Groupings[1] = 0

	return
}