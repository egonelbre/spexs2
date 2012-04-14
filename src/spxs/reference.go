package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	. "spexs"
	"strings"
	"unicode/utf8"
)

func CreateReference(conf Conf) *UnicodeReference {
	ref := NewUnicodeReference(1024)

	if conf.Alphabet.Characters == "" {
		log.Fatal("No alphabet defined!")
	}
	ref.Alphabet = chars(conf.Alphabet.Characters)

	for id, grp := range conf.Alphabet.Groups {
		group := Group{}

		if len(id) != 1 {
			log.Fatal("Group identifier must be of length 1.")
		}

		group.Id = Char(id[0])
		group.Long = "[" + grp.Group + "]"
		group.Chars = chars(grp.Group)

		ref.AddGroup(group)
	}

	if conf.Data.Input == "" {
		log.Fatal("Data file not defined")
	}

	addPatternsFromFile(ref, conf.Data.Input, 0)

	if conf.Data.Validation != "" {
		addPatternsFromFile(ref, conf.Data.Validation, 1)
	}

	return ref
}

func chars(s string) []Char {
	a := make([]Char, 0, len(s))
	for _, c := range s {
		a = append(a, Char(c))
	}
	return a
}

func pattern(data string, group int) ReferencePattern {
	p := ReferencePattern{}
	b := bytes.NewBufferString(data)
	p.Pat = b.Bytes()
	p.Count = utf8.RuneCount(p.Pat)
	p.Group = group
	return p
}

func addPatternsFromFile(ref *UnicodeReference, filename string, group int) {
	var (
		file   *os.File
		reader *bufio.Reader
		line   string
		err    error
	)

	if file, err = os.Open(filename); err != nil {
		log.Println("Did not find data file:", filename)
		log.Fatal(err)
	}

	reader = bufio.NewReader(file)
	for err == nil {
		if line, err = reader.ReadString('\n'); err != nil && err != io.EOF {
			log.Fatal(err)
		}

		line = strings.TrimSpace(line)

		if len(line) == 0 {
			continue
		}

		p := pattern(line, group)
		ref.AddPattern(p)
	}
}
