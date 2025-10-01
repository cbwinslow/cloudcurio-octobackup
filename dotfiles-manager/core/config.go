package core

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// DotfilesConfig represents the configuration for the dotfiles manager
type DotfilesConfig struct {
	SourceDir    string            `yaml:"source_dir"`
	TargetDir    string            `yaml:"target_dir"`
	GitRepo      string            `yaml:"git_repo"`
	AutoApply    bool              `yaml:"auto_apply"`
	IgnoredFiles map[string]bool   `yaml:"ignored_files"`
	Commands     map[string]string `yaml:"commands"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(configPath string) (*DotfilesConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config DotfilesConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %v", err)
	}

	// Set defaults if not provided
	if config.SourceDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %v", err)
		}
		config.SourceDir = filepath.Join(homeDir, ".local", "share", "chezmoi")
	}

	if config.TargetDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %v", err)
		}
		config.TargetDir = homeDir
	}

	return &config, nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *DotfilesConfig, configPath string) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

// DefaultConfig returns a default configuration
func DefaultConfig() *DotfilesConfig {
	homeDir, _ := os.UserHomeDir()
	return &DotfilesConfig{
		SourceDir: filepath.Join(homeDir, ".local", "share", "chezmoi"),
		TargetDir: homeDir,
		AutoApply: false,
		IgnoredFiles: make(map[string]bool),
		Commands:     make(map[string]string),
	}
}