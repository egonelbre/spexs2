package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	. "spexs"

	"spexs/features"
)

func parseCall(call string) (name string, groups []string, info bool, positive bool) {
	regNameArgs, _ := regexp.Compile(`([-+]?)([a-zA-Z0-9]+)(\??)\((.*)\)`)
	regArgs, _ := regexp.Compile("~?[@a-zA-Z0-9]+")
	tokens := regNameArgs.FindStringSubmatch(call)
	if tokens == nil {
		log.Fatalf("Invalid name: %v", call)
	}
	positive = tokens[1] != "-"
	name = tokens[2]
	info = (tokens[3] == "?")
	groups = regArgs.FindAllString(tokens[4], -1)
	return
}

func (s *AppSetup) groupToIds(group string) []int {
	if group == "@" {
		ids := make([]int, len(s.Dataset.Files))
		for i, _ := range ids {
			ids[i] = i
		}
		return ids
	}

	return s.Dataset.Groups[group]
}

func (s *AppSetup) parseFeature(call string) (name string, args []interface{}, info bool, positive bool) {
	name, groups, info, positive := parseCall(call)
	args = make([]interface{}, len(groups))
	for i, group := range groups {
		args[i] = s.groupToIds(group)
	}
	return
}

func (s *AppSetup) makeFeature(call string) Feature {
	feature, _ := s.makeFeatureEx(call)
	return feature
}

func (s *AppSetup) makeFeatureEx(call string) (Feature, bool) {
	name, args, info, positive := s.parseFeature(call)

	normalized := fmt.Sprintf("%+v%+v", name, args)
	if feature, ok := s.Features[normalized]; ok {
		return feature, info
	}

	create, ok := features.Get(name)
	if !ok {
		log.Fatal("No feature named ", name)
	}

	feature, err := features.CallCreateWithArgs(create, args)
	if err != nil {
		log.Fatal(err)
	}

	s.Features[normalized] = feature

	if !positive {
		return func(q *Query) (float64, string) {
			v, info := q.Memoized(feature)
			return -v, info
		}, info
	}

	return feature, info
}

func isDisabled(data []byte) bool {
	var enabled struct{ Enabled *string }
	err := json.Unmarshal(data, &enabled)
	if err != nil {
		log.Fatal(err)
	}

	if (enabled.Enabled != nil) && (*enabled.Enabled == "false") {
		return true
	}
	return false
}
