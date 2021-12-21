package util

import (
	"fmt"
	"io/ioutil"

	"github.com/chwetion/ezsave-for-docker-image/model"
	"gopkg.in/yaml.v3"
)

func LoadConfiguration(path string) (*model.Configuration, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s: %s\n", path, err)
	}
	cfg := &model.Configuration{
		DefaultTarget: model.Image{
			Registry: "docker.io",
			Project:  "library",
		},
		DefaultFrom: model.Image{
			Registry: "docker.io",
			Project:  "library",
		},
	}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse import cfg: %s: %s\n", path, err)
	}
	return cfg, nil
}
