// Package cmd provides the command-line interface for Oris.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command.
var rootCmd = &cobra.Command{
	Use:   "oris",
	Short: "Next-generation embedded vector retrieval engine",
	Long: `Oris is a lightweight, high-performance, embedded vector retrieval engine
designed for modern AI applications.`,
	Version: "0.1.0",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: oris.toml)")
}

func initConfig() {
	if cfgFile != "" {
		_, err := os.Stat(cfgFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "config file not found: %s\n", cfgFile)
			os.Exit(1)
		}
	}
}
