package cmd

import (
	"os"
	"patchy/internal/flags"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "patchy",
	Short: "bad version control system",
	Long:  `patchy is a bad version control system`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&flags.Quiet, "quiet", "q", false, "suppress output")
}
