package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

const baseConfiguration = `{
	"data" : {},
	"alphabet" : {
		"separator" : "",
		"groups" : {}
	},
	"extension" : {
		"method" : "",
		"args" : {},
		"filter" : {}
	},
	"output" : {
		"order" : "",
		"sort" : "desc",
		"filter" : {},
		"queue" : ""
	},
	"print" : {
		"count" : -1,
		"showheader" : true,
		"header" : "",
		"format" : ""
	},
	"aliases" : {}
}`

type FileGroup struct {
	File     string
	Files    []string
	FileList string

	CountSeparator string
}

type Conf struct {
	Data     map[string]FileGroup
	Alphabet struct {
		Separator string
		Groups    map[string]struct{ Elements string }
	}
	Extension struct {
		Method string
		Args   json.RawMessage
		Filter map[string]json.RawMessage
	}
	Output struct {
		Order  string
		Sort   string
		Filter map[string]json.RawMessage
		Queue  string
	}
	Print struct {
		Count      int
		ShowHeader bool
		Header     string
		Format     string
	}
}

func (conf *Conf) WriteToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(file)
	if err = enc.Encode(&conf); err != nil {
		log.Fatal(err)
	}
}

func readBaseConfiguration(config string) *Conf {
	conf := &Conf{}
	dec := json.NewDecoder(bytes.NewBufferString(config))
	if err := dec.Decode(conf); err != nil {
		log.Println("Error in base configuration")
		log.Fatal(err)
	}
	return conf
}

func NewConf(configFile string) *Conf {
	conf := readBaseConfiguration(baseConfiguration)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("Unable to read configuration file: ", configFile)
		log.Fatal(err)
	}

	regArg, _ := regexp.Compile(`^\s*(.*)=(.*)$`)
	for _, arg := range flag.Args() {
		if !regArg.MatchString(arg) {
			log.Fatal("Argument was not in correct form: ", arg)
		}
		tokens := regArg.FindStringSubmatch(arg)
		fmt.Printf("%+v\n", tokens)
		replace, _ := regexp.Compile(`\$` + tokens[1] + `\$`)
		replacement := ([]byte)(tokens[2])
		data = replace.ReplaceAll(data, replacement)
	}

	replace, _ := regexp.Compile(`\$[^\$]*\$`)
	data = replace.ReplaceAll(data, []byte{})

	dec := json.NewDecoder(bytes.NewReader(data))
	if err = dec.Decode(conf); err != nil {
		log.Println("Error in configuration file: ", configFile)
		log.Fatal(err)
	}

	return conf
}
