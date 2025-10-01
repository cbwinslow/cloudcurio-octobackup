package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dotfiles from remote repository",
	Long: `Update fetches changes from the remote repository and applies them.
This command runs chezmoi update to pull the latest changes and apply them.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Manager().Update()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully updated dotfiles!")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}