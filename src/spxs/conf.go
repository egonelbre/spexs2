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
	"data" : {
		"input" : "",
		"reference" : ""
	},
	"alphabet" : {
		"characters" : "",
		"groups" : {},
		"separator" : ""
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
	"aliases" : {
	}
}`

type Conf struct {
	Data struct {
		Input     string
		Reference string
	}
	Alphabet struct {
		Characters string
		Groups     map[string]struct{ Group string }
		Separator  string
	}
	Extension struct {
		Method string
		Args   map[string]extenders.Conf
		Filter map[string]filters.Conf
	}
	Output struct {
		Order  string
		Sort   string
		Args   map[string]map[string]interface{}
		Filter map[string]filters.Conf
		Count  int
		Header string
		Format string
		Queue  string
	}
	Aliases map[string]struct {
		Desc   string
		Modify []string
	}
}

func (conf *Conf) ApplyJson(js string) {
	dec := json.NewDecoder(bytes.NewBufferString(js))
	if err := dec.Decode(&conf); err != nil {
		log.Println("Error in json:", js)
		log.Fatal(err)
	}
}

func (conf *Conf) SetFQN(name string, value string) {
	// extension.filter.pvalue.min
	// convert to {"extension":{"filter":{"pvalue":{"min":value}}}}
	// then apply as an json
	names := strings.Split(name, ".")
	js := value

	_, err := strconv.ParseFloat(value, 64)
	isNumeric := err == nil
	isJson := len(value) > 1 && value[0] == '{'

	if !isNumeric && !isJson {
		// probably a string
		js = `"` + value + `"`
	}

	for i := len(names) - 1; i >= 0; i -= 1 {
		js = `{"` + names[i] + `":` + js + `}`
	}

	conf.ApplyJson(js)
}

func (conf *Conf) Set(ref string, value string) {
	names := make([]string, 1)
	names[0] = ref

	if _, valid := conf.Aliases[ref]; valid {
		names = conf.Aliases[ref].Modify
	}

	for _, name := range names {
		conf.SetFQN(name, value)
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

	return conf
}
