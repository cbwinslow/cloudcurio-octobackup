package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [files...]",
	Short: "Add files to be managed",
	Long: `Add command adds files or directories to be managed by chezmoi.
The files will be copied to the source directory and managed from there.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			err := Manager().AddFile(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error adding file %s: %v\n", file, err)
				os.Exit(1)
			}
			fmt.Printf("Successfully added %s\n", file)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}