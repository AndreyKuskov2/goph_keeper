// Package config provides configuration management for the GophKeeper server.
// It handles loading configuration from YAML files, command-line flags, and environment variables.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// NewConfig creates and loads a new configuration instance.
func NewConfig() (*Config, error) {
	var flags Flags

	pflag.StringVarP(&flags.ConfigPath, "config", "c", "./config/config.yml", "specify this variable to set the path to the config file")

	pflag.Parse()

	for _, arg := range pflag.Args() {
		if !strings.HasPrefix(arg, "-") {
			return nil, fmt.Errorf("unknown flag: %v", arg)
		}
	}

	if err := env.Parse(&flags); err != nil {
		return nil, err
	}

	configRaw, err := os.ReadFile(flags.ConfigPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(configRaw, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
