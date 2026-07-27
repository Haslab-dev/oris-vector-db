package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var compactCmd = &cobra.Command{
	Use:   "compact <path>",
	Short: "Compact a collection's segments",
	Long:  `Manually triggers segment compaction for a collection.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Compact not yet implemented for: %s\n", args[0])
	},
}

func init() {
	rootCmd.AddCommand(compactCmd)
}
