package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"spexs/extenders"
	"spexs/filters"
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
		"args" : {},
		"filter" : {},
		"count": -1,
		"header": "",
		"format": "",
		"queue" : ""
	},
	"print" : {
		"count" : -1,
		"showheader" : true,
		"header" : "",
		"format" : ""
	}
}`

type Conf struct {
	Data     map[string]struct{ Files []string }
	Alphabet struct {
		Separator string
		Groups    map[string]struct{ Elements string }
	}
	Extension struct {
		Method string
		Args   map[string]json.RawMessage
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

func (conf *Conf) ApplyJson(js string) {
	dec := json.NewDecoder(bytes.NewBufferString(js))
	if err := dec.Decode(&conf); err != nil {
		log.Println("Error in json:", js)
		log.Fatal(err)
	}
}

func readBaseConfiguration(config string) Conf {
	var conf Conf
	dec := json.NewDecoder(bytes.NewBufferString(config))
	if err := dec.Decode(&conf); err != nil {
		log.Println("Error in base configuration")
		log.Fatal(err)
	}
	return conf
}

func ReadConfiguration(configs string) Conf {
	configFiles := strings.Split(configs, ",")
	conf := readBaseConfiguration(baseConfiguration)

	for _, configFile := range configFiles {
		if configFile == "" {
			continue
		}

		f, err := os.Open(configFile)
		if err != nil {
			log.Println("Unable to read configuration file: ", configFile)
			continue
		}

		dec := json.NewDecoder(f)

		if err = dec.Decode(&conf); err != nil {
			log.Println("Error in configuration file: ", configFile)
			log.Fatal(err)
		}
	}

	/*
		regArg, _ := regexp.Compile("^\\s*-*(.*)=(.*)$")

		for _, arg := range flag.Args() {
			if !regArg.MatchString(arg) {
				log.Fatal("Argument was not in correct form: ", arg)
			}
			tokens := regArg.FindStringSubmatch(arg)
			name := tokens[1]
			value := tokens[2]

			conf.Set(name, value)
		}
	*/

	return conf
}
