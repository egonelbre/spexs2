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

func NewReferenceFromFile(refName string, charName string) (ref *UnicodeReference, err error) {
	var (
		file *os.File
		reader *bufio.Reader
		line string
	)

	if file, err = os.Open(refName); err != nil {
		return nil, err
	}
	
	ref = &UnicodeReference{}
	ref.Pats = make([]ReferencePattern, 0, 1024)

	reader = bufio.NewReader(file)
	for err == nil {
		if line, err = reader.ReadString('\n'); err != nil && err != io.EOF {
			return nil, fmt.Errorf("Unable to read reference file '%v'", refName)
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

	ref.Groups = make([]Group, 0, 255)
	
	if file, err = os.Open(charName); err != nil {
		return nil, fmt.Errorf("Unable to read charset file '%v'", charName)
	}

	regComment, err1 := regexp.Compile("^\\s*(#.*)?$")
	regAlphabet, err2 := regexp.Compile("^\\s*(\\S+)\\s*(#.*)?$")
	regGroup, err3 := regexp.Compile("^\\s*(\\S+)\\s*,\\s*(\\S)\\s*,\\s*(\\S+)\\s*(#.*)?$")
	
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Printf("Invalid regex: %v\n", err1)
		fmt.Printf("Invalid regex: %v\n", err2)
		fmt.Printf("Invalid regex: %v\n", err3)
		return nil, fmt.Errorf("found invalid regex\n")
	}

	first := true
	lineNo := 0
	reader = bufio.NewReader(file)

	for err == nil {
		if line, err = reader.ReadString('\n'); err != nil && err != io.EOF {
			return nil, err
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
			g := *NewGroup(tokens[3], chars(tokens[2])[0], chars(tokens[1]))
			ref.Groups = append(ref.Groups, g)
			continue
		}

		fmt.Printf("Invalid charset entry on line %v : %v\n", lineNo, line)
	}
	err = nil

	return
}