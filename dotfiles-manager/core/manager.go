package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// DotfilesManager manages the overall dotfiles management process
type DotfilesManager struct {
	config       *DotfilesConfig
	chezmoi      *ChezmoiWrapper
	logger       *log.Logger
}

// NewDotfilesManager creates a new instance of DotfilesManager
func NewDotfilesManager(config *DotfilesConfig, logger *log.Logger) *DotfilesManager {
	chezmoiConfig := &Config{
		SourceDir: config.SourceDir,
		TargetDir: config.TargetDir,
	}
	
	manager := &DotfilesManager{
		config:  config,
		chezmoi: NewChezmoiWrapper(chezmoiConfig),
		logger:  logger,
	}
	
	return manager
}

// Initialize sets up the dotfiles manager
func (dm *DotfilesManager) Initialize() error {
	dm.logger.Println("Initializing dotfiles manager...")
	
	// Check if chezmoi is installed
	if !dm.isChezmoiInstalled() {
		return fmt.Errorf("chezmoi is not installed. Please install chezmoi first")
	}
	
	// Ensure source directory exists
	if err := os.MkdirAll(dm.config.SourceDir, 0755); err != nil {
		return fmt.Errorf("error creating source directory: %v", err)
	}
	
	dm.logger.Printf("Source directory: %s\n", dm.config.SourceDir)
	dm.logger.Printf("Target directory: %s\n", dm.config.TargetDir)
	
	return nil
}

// isChezmoiInstalled checks if chezmoi is installed
func (dm *DotfilesManager) isChezmoiInstalled() bool {
	_, err := dm.chezmoi.RunCommand("--version")
	return err == nil
}

// AddFile adds a file to be managed
func (dm *DotfilesManager) AddFile(filePath string) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("error getting absolute path: %v", err)
	}
	
	dm.logger.Printf("Adding file: %s", absPath)
	
	return dm.chezmoi.Add(absPath)
}

// Apply applies all dotfiles
func (dm *DotfilesManager) Apply() error {
	dm.logger.Println("Applying dotfiles...")
	
	return dm.chezmoi.Apply()
}

// Status shows the status of managed files
func (dm *DotfilesManager) Status() (string, error) {
	dm.logger.Println("Getting status...")
	
	return dm.chezmoi.Status()
}

// Diff shows differences
func (dm *DotfilesManager) Diff() (string, error) {
	dm.logger.Println("Getting diff...")
	
	return dm.chezmoi.Diff()
}

// Update updates from remote
func (dm *DotfilesManager) Update() error {
	dm.logger.Println("Updating from remote...")
	
	return dm.chezmoi.Update()
}

// ManagedFiles returns list of managed files
func (dm *DotfilesManager) ManagedFiles() ([]string, error) {
	return dm.chezmoi.Managed()
}

// EditFile opens editor for a managed file
func (dm *DotfilesManager) EditFile(filePath string) error {
	return dm.chezmoi.Edit(filePath)
}