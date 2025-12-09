package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// CollectionConfig represents the configuration for a collection.
type CollectionConfig struct {
	Name         string            `yaml:"name"`
	Help         string            `yaml:"help"`
	Readme       string            `yaml:"readme"`
	Runnables    []string          `yaml:"runnables"`
	Environments map[string]string `yaml:"environments"`
}

// RunnableConfig represents the configuration for a runnable.
type RunnableConfig struct {
	Name         string            `yaml:"name"`
	Help         string            `yaml:"help"`
	Readme       string            `yaml:"readme"`
	Run          string            `yaml:"run"`
	Before       string            `yaml:"before"`
	After        string            `yaml:"after"`
	Environments map[string]string `yaml:"environments"`
}

// LoadCollectionConfig loads the collection configuration from the given path.
func LoadCollectionConfig(path string) (*CollectionConfig, error) {
	data, err := os.ReadFile(filepath.Join(path, "collection.yml"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var cfg CollectionConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse collection.yml: %w", err)
	}
	return &cfg, nil
}

// LoadRunnableConfig loads the runnable configuration from the given path.
func LoadRunnableConfig(path string) (*RunnableConfig, error) {
	data, err := os.ReadFile(filepath.Join(path, "runnable.yml"))

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var cfg RunnableConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse runnable.yml: %w", err)
	}
	return &cfg, nil
}
