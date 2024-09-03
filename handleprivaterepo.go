package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"golang.org/x/mod/modfile"
)

func cleanPrivateRepo(_ ProjectSpec) goPackageFunc {
	log.Print("cleaning private repos...")
	return func(pkgDirName, actDirName string) error {
		prefix := pkgPrefix(pkgDirName, actDirName)

		// remove private repo dir
		log.Print(prefix("removing private repo dir... "))
		privateRepoDir := path.Join(actDirName, PrivateRepoDir)
		if _, err := os.Stat(privateRepoDir); os.IsNotExist(err) {
			log.Print(prefix("skip: private repo dir does not exist"))
			return nil
		}
		if err := os.RemoveAll(privateRepoDir); err != nil {
			return fmt.Errorf(prefix("error: failed removing private repo dir: %w"), prefix, err)
		}

		// restore go.mod to original state
		log.Print(prefix("restoring go.mod to original state... "))
		goMod := path.Join(actDirName, "go.mod")

		// find all go.mod.*.bak files
		var goModBak string
		files, err := os.ReadDir(actDirName)
		if err != nil {
			return fmt.Errorf(prefix("error: failed reading dir: %w"), err)
		}
		for _, f := range files {
			if strings.HasPrefix(f.Name(), "go.mod.") && strings.HasSuffix(f.Name(), ".bak") {
				goModBak = path.Join(actDirName, f.Name())
				break
			}
		}

		// test checksum
		if goModBak == "" {
			return fmt.Errorf(prefix("error: failed finding go.mod backup"))
		}

		hash := strings.TrimSuffix(strings.TrimPrefix(path.Base(goModBak), "go.mod."), ".bak")
		modFile, err := os.ReadFile(goMod)
		if err != nil {
			return fmt.Errorf(prefix("error: failed reading go.mod: %w"), err)
		}
		h := sha1.New()
		h.Write(modFile)
		sum := h.Sum(nil)
		base64sum := base64.StdEncoding.EncodeToString(sum)
		if hash != base64sum {
			log.Printf(prefix("checksum: %s"), base64sum)
			log.Printf(prefix("hash: %s"), hash)
			return fmt.Errorf(prefix("error: go.mod checksum mismatch backup: %s hash: %s"), goModBak)
		}

		if f, err := os.Stat(goModBak); err == nil && !f.IsDir() {
			if err := os.Remove(goMod); err != nil {
				return fmt.Errorf(prefix("error: failed removing go.mod: %w"), err)
			}
			if err := os.Rename(goModBak, goMod); err != nil {
				return fmt.Errorf(prefix("error: failed restoring go.mod: %w"), err)
			}
		}

		return nil
	}
}

func clonePrivateRepo(projectCfg ProjectSpec) goPackageFunc {
	log.Print("checking for private repos...")

	// check if GOPRIVATE is set in projectCfg
	var privateRepos []string
	for k, v := range projectCfg.Environment {
		if k == "GOPRIVATE" {
			privateRepos = strings.Split(v, ",")
			break
		}
	}

	// if no private repos defined, skip
	if len(privateRepos) == 0 {
		log.Print("skipping: no private repos defined in GOPRIVATE environment")
		return nil
	}
	log.Print("cloning private repos for each go project...")

	return func(pkgDirName, actDirName string) error {
		prefix := pkgPrefix(pkgDirName, actDirName)

		goMod := path.Join(actDirName, "go.mod")
		modFile, err := os.ReadFile(goMod)
		if err != nil {
			return fmt.Errorf(prefix("error: failed reading go.mod: %w"), err)
		}
		mod, err := modfile.Parse(goMod, modFile, nil)
		if err != nil {
			return fmt.Errorf(prefix("error: failed parsing go.mod: %w"), err)
		}

		// make private repo dir
		log.Print(prefix("creating private repo dir... "))
		privateRepoDir := path.Join(actDirName, PrivateRepoDir)
		if _, err := os.Stat(privateRepoDir); os.IsNotExist(err) {
			os.Mkdir(privateRepoDir, 0755)
		} else {
			return fmt.Errorf(prefix("error: private repo dir already exists"))
		}

		// for each require check if it is a private repo
		goModChange := false
		for _, r := range mod.Require {
			for _, p := range privateRepos {
				if match, err := path.Match(p, r.Mod.Path); err != nil {
					return fmt.Errorf(prefix("error: failed matching private repo: %w"), err)
				} else if match {
					clonePath := path.Join(privateRepoDir, r.Mod.Path)
					newModPath := path.Join(".", PrivateRepoDir, r.Mod.Path)

					log.Printf(prefix("found private repo: %s"), r.Mod.Path)

					// clone private repo
					log.Print(prefix("cloning private repo... "))

					// recursively create private repo dir
					dir := privateRepoDir
					for _, d := range strings.Split(path.Dir(r.Mod.Path), "/") {
						dir = path.Join(dir, d)
						log.Printf(prefix("creating dir: %s"), dir)
						if _, err := os.Stat(privateRepoDir); os.IsNotExist(err) {
							os.Mkdir(privateRepoDir, 0755)
						}
					}

					cmd := exec.Command("git", "clone", "https://"+r.Mod.Path, clonePath)
					o, err := cmd.CombinedOutput()
					if err != nil {
						return fmt.Errorf(prefix("error: failed cloning private repo: %s"), o)
					}

					// modify go.mod to use private repo
					log.Print(prefix("modifying go.mod to use private repo... "))
					if err := mod.AddReplace(r.Mod.Path, r.Mod.Version, newModPath, ""); err != nil {
						return fmt.Errorf(prefix("error: failed adding replace to go.mod: %w"), err)
					}

					goModChange = true
				}
			}
		}

		if !goModChange {
			log.Print(prefix("no private repos found"))
			return nil
		}

		// write go.mod
		b, err := mod.Format()
		if err != nil {
			return fmt.Errorf(prefix("error: failed formatting go.mod: %w"), err)
		}

		// sha1 hash of bytes
		h := sha1.New()
		h.Write(b)
		sum := h.Sum(nil)
		base64sum := base64.StdEncoding.EncodeToString(sum)

		// backup go.mod
		if err := os.Rename(goMod, goMod+"."+base64sum+".bak"); err != nil {
			return fmt.Errorf(prefix("error: failed renaming go.mod: %w"), err)
		}

		// write go.mod
		if err := os.WriteFile(goMod, b, 0644); err != nil {
			return fmt.Errorf(prefix("error: failed writing go.mod: %w"), err)
		}

		return nil
	}
}
