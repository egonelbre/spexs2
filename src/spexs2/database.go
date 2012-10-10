package main

import (
	"bufio"
	"io"
	"log"
	"os"
	. "spexs"
	"strings"
)

func removeEmpty(names []string) []string {
	result := make([]string, len(names))
	i := 0
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed != "" {
			result[i] = trimmed
			i += 1
		}
	}
	return result[0:i]
}

func CreateDatabase(conf Conf) *Database {
	db := NewDatabase(1024)
	db.Separator = conf.Alphabet.Separator

	for id, grp := range conf.Alphabet.Groups {
		group := Group{}

		group.Name = id
		group.FullName = "[" + grp.Elements + "]"

		tokenNames := strings.Split(grp.Elements, db.Separator)
		tokenNames = removeEmpty(tokenNames)
		group.Elems = db.ToTokens(tokenNames)

		db.AddGroup(group)
	}

	if conf.Data.Input == "" {
		log.Fatal("Data file not defined")
	}

	addSeqsFromFile(db, conf.Data.Input, 0)
	if conf.Data.Reference != "" {
		addSeqsFromFile(db, conf.Data.Reference, 1)
	}

	return db
}

func addSeqsFromFile(db *Database, filename string, section int) {
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
		tokenNames := strings.Split(line, db.Separator)
		tokenNames = removeEmpty(tokenNames)
		tokens := db.ToTokens(tokenNames)

		if len(tokens) <= 0 {
			continue
		}

		seq := Sequence{
			Tokens:  tokens,
			Len:     len(tokens),
			Section: section,
			Count:   1,
		}
		db.AddSequence(seq)
	}
}
