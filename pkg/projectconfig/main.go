package projectconfig

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrNotFound = fmt.Errorf("could not find config file")
var ErrParse = fmt.Errorf("error parsing config file")

func Parse(cfgPath string) (*ProjectSpec, error) {
	cfg := &ProjectSpec{}
	b, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, errors.Join(ErrNotFound, err)
	}

	// parse yaml
	err = yaml.Unmarshal(b, cfg)
	if err != nil {
		return nil, errors.Join(ErrParse, err)
	}
	return cfg, nil
}
