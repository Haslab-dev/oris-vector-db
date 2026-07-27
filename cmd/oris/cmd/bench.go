package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var benchCmd = &cobra.Command{
	Use:   "bench",
	Short: "Run benchmarks",
	Long:  `Runs performance benchmarks for insert, search, and recall.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Benchmarks not yet implemented.")
	},
}

func init() {
	rootCmd.AddCommand(benchCmd)
}
