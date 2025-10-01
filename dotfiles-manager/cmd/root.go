package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"dotfiles-manager/core"
)

var (
	cfgFile string
	cfg     *core.DotfilesConfig
	manager *core.DotfilesManager
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "A CLI for managing dotfiles with Chezmoi",
	Long: `dotfiles is a comprehensive tool for managing your dotfiles with both CLI and TUI interfaces.
It provides a wrapper around Chezmoi with enhanced functionality and user experience.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/dotfiles/config.yaml)")
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	var err error
	if cfgFile != "" {
		// Use config file from the flag
		cfg, err = core.LoadConfig(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		
		// Look for config in home directory with name ".dotfiles" (without extension)
		cfg, err = core.LoadConfig(home + "/.config/dotfiles/config.yaml")
		if err != nil {
			// Use default config if not found
			fmt.Fprintf(os.Stderr, "Config file not found, using defaults: %v\n", err)
			cfg = core.DefaultConfig()
		}
	}
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	logger := log.New(os.Stdout, "[dotfiles] ", log.LstdFlags)
	manager = core.NewDotfilesManager(cfg, logger)
	
	// Initialize the manager
	if err := manager.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing manager: %v\n", err)
		os.Exit(1)
	}
}

// Manager returns the global manager instance
func Manager() *core.DotfilesManager {
	return manager
}

// GetConfig returns the global config instance
func GetConfig() *core.DotfilesConfig {
	return cfg
}