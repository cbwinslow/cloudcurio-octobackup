package core

import (
	"fmt"
	"os/exec"
	"strings"
)

// Config holds the configuration for the dotfiles manager
type Config struct {
	SourceDir string
	TargetDir string
	DryRun    bool
	Verbose   bool
}

// ChezmoiWrapper wraps common Chezmoi operations
type ChezmoiWrapper struct {
	config *Config
}

// NewChezmoiWrapper creates a new instance of ChezmoiWrapper
func NewChezmoiWrapper(config *Config) *ChezmoiWrapper {
	return &ChezmoiWrapper{
		config: config,
	}
}

// RunCommand executes a chezmoi command
func (c *ChezmoiWrapper) RunCommand(args ...string) (string, error) {
	cmd := exec.Command("chezmoi", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running chezmoi command: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

// Status returns the status of managed files
func (c *ChezmoiWrapper) Status() (string, error) {
	return c.RunCommand("status")
}

// Apply applies the dotfiles to the target system
func (c *ChezmoiWrapper) Apply() error {
	_, err := c.RunCommand("apply")
	return err
}

// Add adds a file to be managed by chezmoi
func (c *ChezmoiWrapper) Add(filePath string) error {
	_, err := c.RunCommand("add", filePath)
	return err
}

// Diff shows the differences between source and target
func (c *ChezmoiWrapper) Diff() (string, error) {
	return c.RunCommand("diff")
}

// Init initializes a new chezmoi repository
func (c *ChezmoiWrapper) Init() error {
	_, err := c.RunCommand("init")
	return err
}

// Update updates the source directory from the remote repository
func (c *ChezmoiWrapper) Update() error {
	_, err := c.RunCommand("update")
	return err
}

// Edit opens an editor to edit a managed file
func (c *ChezmoiWrapper) Edit(filePath string) error {
	_, err := c.RunCommand("edit", filePath)
	return err
}

// Managed returns a list of managed files
func (c *ChezmoiWrapper) Managed() ([]string, error) {
	output, err := c.RunCommand("managed")
	if err != nil {
		return nil, err
	}
	
	// Split the output by lines and clean up
	lines := strings.Split(output, "\n")
	var files []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	
	return files, nil
}