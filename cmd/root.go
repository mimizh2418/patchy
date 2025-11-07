package cmd

import (
	"os"
	"patchy/internal/flags"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "patchy <command> [<args>]",
	Short: "Bad version control system",
	Long:  `Patchy is a bad version control system`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().BoolVarP(&flags.Quiet, "quiet", "q", false, "suppress output")
}
