package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of managed files",
	Long: `Status shows the status of managed files similar to git status.
It runs chezmoi status to show which files differ from the source state.`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := Manager().Status()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting status: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(status)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}