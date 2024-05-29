package config

import (
	"github.com/alkosuv/in-mem-kv/internal/logger"
	"github.com/alkosuv/in-mem-kv/internal/network"
	"gopkg.in/yaml.v3"
	"os"
)

// const DefaultConfigPath string = "/etc/in-mem-kv/error.log"
const DefaultConfigPath string = "configs/config.dev.yaml"

type Config struct {
	Network network.TCPServerConfig `yaml:"network"`
	Logging logger.Config           `yaml:"logging"`
	Engine  struct {
		Type string `yaml:"type"`
	} `yaml:"engine"`
}

func New(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigPath
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
