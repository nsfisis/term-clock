package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "term-clock",
	Short: "A clock on your terminal",
}

func init() {
	rootCmd.AddCommand(clockCmd)
	rootCmd.AddCommand(timerCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
