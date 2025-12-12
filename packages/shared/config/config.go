package config

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DevicesPort int  `yaml:"devicesPort"`
	GatewayPort int  `yaml:"gatewayPort"`
	Debug       bool `yaml:"debug"`
}

func (c *Config) Load(logger *zap.Logger, path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(file, c); err != nil {
		return fmt.Errorf("failed to unmarshal config file %s: %w", path, err)
	}

	logger.Info("config loaded successfully", zap.String("path", path))
	return nil
}
