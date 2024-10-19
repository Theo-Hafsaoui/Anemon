package input

import (
	core "anemon/internal/core"
	"gopkg.in/yaml.v3"
	"os"
)

type YamlSource struct{}

func (*YamlSource) GetParamsFrom(root string) (core.Params, error) {
	params := core.Params{}
	yamlFile, err := os.ReadFile(root + "/params.yml")
	if err != nil {
		return params, err
	}
	err = yaml.Unmarshal(yamlFile, &params)
	if err != nil {
		return params, err
	}
	return params, nil
}
