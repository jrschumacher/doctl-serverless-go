package main

import (
	"errors"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

func pkgPrefix(s ...string) func(string) string {
	for i, v := range s {
		s[i] = path.Base(v)
	}

	return func(m string) string {
		return "[" + strings.Join(s, "/") + "] " + m
	}
}

func parseCfg(monorepoPath, cfgPath string, cfg interface{}) error {
	for _, ext := range CfgExt {
		b, err := os.ReadFile(path.Join(monorepoPath, cfgPath+ext))
		if err == nil {
			// check if cfg is nil and return
			if cfg == nil {
				return nil
			}
			// parse yaml
			err = yaml.Unmarshal(b, cfg)
			if err != nil {
				return errors.Join(ErrCfgParse, err)
			}
			return nil
		}
	}
	return ErrCfgNotFound
}
