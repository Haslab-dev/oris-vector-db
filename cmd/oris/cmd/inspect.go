package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <path>",
	Short: "Inspect a collection or segment",
	Long:  `Prints metadata and statistics about a collection or segment.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inspect not yet implemented for: %s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
