package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	
	"dotfiles-manager/core"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new dotfiles repository",
	Long: `Initialize sets up a new dotfiles repository with Chezmoi.
This command runs chezmoi init to set up your source directory and create 
initial configuration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Create a temporary ChezmoiWrapper to run init
		tempConfig := &core.Config{
			SourceDir: GetConfig().SourceDir,
			TargetDir: GetConfig().TargetDir,
		}
		chezmoi := core.NewChezmoiWrapper(tempConfig)
		
		err := chezmoi.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully initialized dotfiles repository!")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}