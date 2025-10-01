package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply dotfiles to the target system",
	Long: `Apply synchronizes the target system with the source dotfiles.
It runs chezmoi apply to update the target files according to the source.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Manager().Apply()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error applying dotfiles: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully applied dotfiles!")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}