package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd represents the `oris init` command.
var initCmd = &cobra.Command{
	Use:   "init <path>",
	Short: "Initialize a new Oris workspace",
	Long: `Creates a new Oris workspace at the specified path with the
required directory structure for storing collections.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		dirs := []string{
			path,
			path + "/collections",
			path + "/collections/docs",
			path + "/collections/wal",
			path + "/collections/snapshots",
			path + "/collections/segments",
		}
		for _, d := range dirs {
			if err := os.MkdirAll(d, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "error creating %s: %v\n", d, err)
				os.Exit(1)
			}
		}

		fmt.Printf("Initialized Oris workspace at %s\n", path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
