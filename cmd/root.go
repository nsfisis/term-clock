package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "term-clock",
	Short: "A clock on your terminal",
	// Run 'clock' subcommand by default.
	Run: func(cmd *cobra.Command, args []string) {
		clockCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(alarmCmd)
	rootCmd.AddCommand(clockCmd)
	rootCmd.AddCommand(timerCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
