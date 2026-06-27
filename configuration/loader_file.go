package configuration

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type fileLoader struct {
	Concurrents *int64         `yaml:"concurrents"`
	Filepath    *string        `yaml:"filepath"`
	Interval    *time.Duration `yaml:"interval"`
	Size        *string        `yaml:"payload-size"`
	Tor         *string        `yaml:"tor"`
	Method      *string        `yaml:"method"`
	Headers     *[]string      `yaml:"headers"`
	URL         *string        `yaml:"url"`
}

func NewFileLoader(path string) (*fileLoader, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read configuration file: %v", err)
	}

	var config fileLoader

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse unmarshal into the file configuration struct: %v", err)
	}

	return &config, nil
}
