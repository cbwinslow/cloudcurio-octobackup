package tui

import (
	"log"
	"os"

	"dotfiles-manager/core"
)

// RunTUI runs the terminal user interface
func RunTUI() error {
	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Try to load config, if not found use default
	configPath := home + "/.config/dotfiles/config.yaml"
	cfg, err := core.LoadConfig(configPath)
	if err != nil {
		// If no config is found, use default
		cfg = core.DefaultConfig()
	}

	logger := log.New(os.Stdout, "[dotfiles-tui] ", log.LstdFlags)
	manager := core.NewDotfilesManager(cfg, logger)

	// Initialize the manager
	if err := manager.Initialize(); err != nil {
		return err
	}

	// Create and run the TUI
	model := NewModel(manager)
	return model.Run()
}

// NewModel creates a new TUI model
func NewModel(manager *core.DotfilesManager) Model {
	return Model{
		manager:  manager,
		choices:  []string{"Status", "Apply", "Add File", "View Diff", "Update", "Quit"},
		selected: make(map[int]struct{}),
		currentView: "menu", // menu, status, apply, etc.
		files:    []string{},
		ready:    false,
	}
}