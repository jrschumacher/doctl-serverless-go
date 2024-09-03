package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
)

const ProjectCfgFile = "project"
const DOConfigDir = ".do"
const DeployTemplateCfgFile = "deploy.template"

var PrivateRepoDir string

var CfgExt = []string{".yaml", ".yml"}

type goPackageFunc func(pkgDirName, actDirName string) error

const usageMessage = `
usage: doctl-serverless-go <command> [<monorepo-path>]

	deploy	Deploy the go monorepo to digitalocean serverless
	clean	Clean the go monorepo after deployment

`

var ErrInvalidUsage = errors.New("invalid usage")

func main() {
	var err error

	flag.StringVar(&PrivateRepoDir, "private", PrivateRepoDir, "private repo directory")
	flag.Parse()

	subcommand := flag.Arg(0)
	if subcommand == "" || (subcommand != "deploy" && subcommand != "clean") {
		exitInvalidUsage()
	}

	monorepoPath := flag.Arg(1)
	if monorepoPath == "" {
		monorepoPath = "."
	}
	log.Printf("scanning monorepo... path:\"%s\"", monorepoPath)
	if _, err := os.Stat(monorepoPath); err != nil {
		log.Fatal(errors.Join(ErrNotMonorepo, err))
	}

	// read project yaml
	log.Print("parsing project.yaml... ")
	var projectCfg ProjectSpec
	if err := parseCfg(monorepoPath, ProjectCfgFile, &projectCfg); err != nil {
		log.Fatal(errors.Join(ErrParseProjectCfg, errors.Join(ErrNotMonorepo, err)))
	}

	// stat packages
	pkgsDirName := path.Join(monorepoPath, "packages")
	pkgsDirStat, err := os.Stat(pkgsDirName)
	if err != nil || !pkgsDirStat.IsDir() {
		log.Fatal(errors.Join(ErrNotMonorepo, err))
	}

	errs := make(errMap)
	switch subcommand {
	case "deploy":
		errs = forEveryPackage(pkgsDirName, projectCfg,
			clonePrivateRepo(projectCfg),
		)
	case "clean":
		errs = forEveryPackage(pkgsDirName, projectCfg,
			cleanPrivateRepo(projectCfg),
		)
	}

	if errs.HasErrors() {
		errs.Print()
		os.Exit(1)
	}

	log.Print("completed")
}

func exitInvalidUsage() {
	fmt.Println("fatal: missing command")
	fmt.Print(usageMessage)
	os.Exit(1)
}

func forEveryPackage(pkgsDirName string, projectCfg ProjectSpec, fns ...goPackageFunc) errMap {
	pkgErrs := make(map[string]error)
	for _, pkg := range projectCfg.Packages {
		scope := pkg.Name
		pkgDirName := path.Join(pkgsDirName, pkg.Name)

		// check if exists and is a directory
		pkgDirStat, err := os.Stat(pkgDirName)
		if err != nil || !pkgDirStat.IsDir() {
			pkgErrs[scope] = errors.Join(ErrNoPackagesFound, err)
			continue
		}

		// for each action
		for _, act := range pkg.Actions {
			scope := path.Join(pkg.Name, act.Name)
			actDirName := path.Join(pkgDirName, act.Name)

			// check if exists and is a directory
			actDirStat, err := os.Stat(actDirName)
			if err != nil || !actDirStat.IsDir() {
				pkgErrs[scope] = errors.Join(ErrNoPackagesFound, err)
				continue
			}

			// is go package? check for go.mod
			modPath := path.Join(actDirName, "go.mod")
			modStat, err := os.Stat(modPath)
			if err != nil {
				if !os.IsNotExist(err) {
					continue
				}
				pkgErrs[scope] = errors.Join(ErrCheckingForMod, err)
			}
			if modStat.IsDir() {
				continue
			}

			// run functions
			for _, fn := range fns {
				// reflect fn name
				fnName := reflect.TypeOf(fn).Name()
				scope = scope + " (" + fnName + ")"
				if err := fn(pkgDirName, actDirName); err != nil {
					pkgErrs[scope] = err
				}
			}
		}
	}
	return pkgErrs
}
