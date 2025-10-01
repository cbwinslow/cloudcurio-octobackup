package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"dotfiles-manager/tui"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the terminal user interface",
	Long: `TUI launches the interactive terminal user interface for managing dotfiles.
This provides a visual way to manage your dotfiles with menus and prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Run the TUI
		err := tui.RunTUI()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}