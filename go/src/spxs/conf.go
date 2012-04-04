package main

import (
	"flag"
	"log"
	"os"
	"encoding/json"
	"regexp"
	"bytes"
	"strings"
	"strconv"
)

const baseConfiguration = `{
	"Data" : {
		"Input" : "",
		"Validation" : ""
	},
	"Alphabet" : {
		"Characters" : "",
		"Groups" : {}
	},
	"Extension" : {
		"Method" : {},
		"Filter" : {}
	},
	"Output" : {
		"Order" : {},
		"Filter" : {},
		"Show": -1
	},
	"Aliases" : {
		"ref" : {"Desc":"input file", "Modify":["Data.Input"]},
		"val" : {"Desc":"validation file", "Modify":["Data.Validation"]}
	}
}`

type Conf struct {
	Data struct {
		Input string
		Validation string
	}
	Alphabet struct {
		Characters string
		Groups map[string]struct {
			Group string
		}
	}
	Extension struct {
		Method map[string]interface{}
		Filter map[string]interface{}
	}
	Output struct {
		Order map[string]interface{}
		Filter map[string]interface{}
		Show int
	}
	Aliases map[string]struct {
		Desc string
		Modify []string
	}
}

func (conf *Conf) ApplyJson(js string){
	dec := json.NewDecoder(bytes.NewBufferString(js))
	if err := dec.Decode(&conf); err != nil {
		log.Println("Error in json:", js)
		log.Fatal(err)
	}
}

func (conf *Conf) SetFQN(name string, value string){
	// extension.filter.pvalue.min
	// convert to {"extension":{"filter":{"pvalue":{"min":value}}}}
	names := strings.Split(name, ".")
	js := value

	if _, err := strconv.ParseFloat(value, 64); err != nil{
		js = `"` + value + `"`;
	}
	
	for i := len(names)-1 ; i >= 0; i -= 1 {
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

func readConfiguration(configFiles []string) Conf {
	conf := readBaseConfiguration(baseConfiguration)

	for _, configFile := range configFiles {
		if configFile == "" {
			continue
		}

		f, err := os.Open(configFile)
		if err != nil {
			log.Println("Unable to read configuration file: ", configFile)
			log.Fatal(err)
		}

		dec := json.NewDecoder(f)

		if err = dec.Decode(&conf); err != nil {
			log.Println("Error in configuration file: ", configFile)
			log.Fatal(err)
		}
	}

	regArg, _ := regexp.Compile("^\\s*-*(.*)=(.*)$")

	for _, arg := range flag.Args() {
		log.Println("processing: ", arg)
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