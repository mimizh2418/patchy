package cmd

import (
	"os"
	"patchy/cmd/backend/catfile"
	"patchy/cmd/backend/committree"
	"patchy/cmd/backend/updateref"
	"patchy/cmd/backend/writeblob"
	"patchy/cmd/backend/writetree"
	"patchy/cmd/frontend/commit"
	"patchy/cmd/frontend/initialize"
	"patchy/util"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "patchy <command> [<args>]",
	Short: "Bad version control system",
	Long:  `Patchy is a bad version control system`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.SilenceUsage = true
	RootCmd.PersistentFlags().BoolVarP(&util.Quiet, "quiet", "q", false, "suppress output")

	RootCmd.AddCommand(catfile.NewCommand())
	RootCmd.AddCommand(committree.NewCommand())
	RootCmd.AddCommand(writeblob.NewCommand())
	RootCmd.AddCommand(updateref.NewCommand())
	RootCmd.AddCommand(writetree.NewCommand())

	RootCmd.AddCommand(commit.NewCommand())
	RootCmd.AddCommand(initialize.NewCommand())
}
