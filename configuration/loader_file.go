// Package configuration handles the static config.
package configuration

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// FileLoader will load the configuration from a file.
type FileLoader struct {
	Concurrents *int64         `yaml:"concurrents"`
	Filepath    *string        `yaml:"filepath"`
	Interval    *time.Duration `yaml:"interval"`
	Size        *string        `yaml:"payload-size"` //nolint:tagliatelle
	Tor         *string        `yaml:"tor"`
	Method      *string        `yaml:"method"`
	Headers     *[]string      `yaml:"headers"`
	URL         *string        `yaml:"url"`
	Protocol    *string        `yaml:"protocol"`
	Throttle    *time.Duration `yaml:"throttle"`
	Insecure    *bool          `yaml:"insecure"`
}

// NewFileLoader loads the configuration from the file.
func NewFileLoader(path string) (*FileLoader, error) {
	content, err := os.ReadFile(path) //nolint:gosec // it's okay
	if err != nil {
		return nil, fmt.Errorf("cannot read configuration file: %w", err)
	}

	var config FileLoader

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse unmarshal into the file configuration struct: %w", err)
	}

	return &config, nil
}
