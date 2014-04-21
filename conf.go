package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/egonelbre/spexs2/vendor/rjson"
)

const baseConfiguration = `{
	"Dataset" : {},
	"Reader" : {
		"Method" : "Delimited",
		"Separator" : "",
		"CountSeparator" : "",
		"Skip" : ""
	},
	"Extension" : {
		"Method" : "",
		"Groups" : {},
		"Extendable" : {},
		"Outputtable" : {}
	},
	"Output" : {
		"Method" : "Sorted",
		"Count" : -1
	},
	"Print" : {
		"ShowHeader" : true
	}
}`

type FileGroup struct {
	File     string
	Files    []string
	FileList string
}

type Conf struct {
	Dataset map[string]FileGroup
	Reader  struct {
		Method         string
		Separator      string
		CountSeparator string
		Skip           string
	}
	Extension struct {
		Method string
		Groups map[string]struct{ Elements string }

		Extendable  map[string]rjson.RawMessage
		Outputtable map[string]rjson.RawMessage
	}
	Output struct {
		Method string
		SortBy []string
		Count  int
	}
	Printer struct {
		Method     string
		ShowHeader bool
		Reverse    bool
		Header     string
		Format     string
	}
}

func (conf *Conf) WriteToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	enc := rjson.NewEncoder(file)
	if err = enc.Encode(&conf); err != nil {
		log.Fatal(err)
	}
}

func readBaseConfiguration(config string) *Conf {
	conf := &Conf{}
	dec := rjson.NewDecoder(bytes.NewBufferString(config))
	if err := dec.Decode(conf); err != nil {
		log.Println("Error in base configuration")
		log.Fatal(err)
	}
	return conf
}

func NewConf(configFile string) *Conf {
	conf := readBaseConfiguration(baseConfiguration)

	if configFile == "" {
		return conf
	}

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

		replace, _ := regexp.Compile(`\$` + tokens[1] + `(=[^$]*)?\$`)
		replacement := ([]byte)(tokens[2])
		data = replace.ReplaceAll(data, replacement)
	}

	regDefaults, _ := regexp.Compile(`\$[^\$]+\$`)
	regDefault, _ := regexp.Compile(`\$[^=]+=(.*)\$`)
	data = regDefaults.ReplaceAllFunc(data, func(repl []byte) []byte {
		defaults := regDefault.FindSubmatch(repl)
		if defaults != nil {
			return defaults[1]
		}
		return nil
	})

	dec := rjson.NewDecoder(bytes.NewReader(data))
	if err = dec.Decode(conf); err != nil {
		log.Println("Error in configuration file: ", configFile)
		log.Fatal(err)
	}

	return conf
}
