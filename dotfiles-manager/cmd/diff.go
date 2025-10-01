package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show differences between source and target",
	Long: `Diff shows the differences between the source state and the target system.
Similar to running 'chezmoi diff' to see what would change on apply.`,
	Run: func(cmd *cobra.Command, args []string) {
		diff, err := Manager().Diff()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting diff: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(diff)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
}