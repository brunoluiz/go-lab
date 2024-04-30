package main

import (
	"os"

	"github.com/google/ko/pkg/build"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Docker struct {
		Registry string
	} `yaml:"docker"`
	Builds           []build.Config `yaml:"builds"`
	DefaultPlatforms []string       `yaml:"defaultPlatforms"`
}

func loadBobYAML() (*Config, error) {
	yamlFile, err := os.ReadFile(".bob.yaml")
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
