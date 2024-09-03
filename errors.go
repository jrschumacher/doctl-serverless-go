package main

import (
	"fmt"
	"log"
)

var ErrCfgNotFound = fmt.Errorf("could not find config file")
var ErrCfgParse = fmt.Errorf("error parsing config file")
var ErrParseProjectCfg = fmt.Errorf("error parsing project.yaml")
var ErrParseDeployTemplateCfg = fmt.Errorf("error parsing deploy.template.yaml")

var ErrNotMonorepo = fmt.Errorf("not a digitalocean serverless monorepo")
var ErrNoPackagesFound = fmt.Errorf("no packages found")
var ErrCheckingForMod = fmt.Errorf("error checking for go.mod")

type errMap map[string]error

func (e errMap) HasErrors() bool {
	return len(e) > 0
}

func (e errMap) Print() {
	if !e.HasErrors() {
		return
	}

	for _, v := range e {
		log.Printf("%s", v)
	}
}
