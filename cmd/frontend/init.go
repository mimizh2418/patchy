package frontend

import (
	"patchy/cmd"
	"patchy/repo"
	"patchy/util"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [<directory>]",
	Short: "Create an empty repository",
	Long:  `Create an empty repository`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var repoPath string
		var err error
		if len(args) > 0 {
			repoPath, err = repo.Init(args[0])
		} else {
			repoPath, err = repo.Init(".")
		}
		if err != nil {
			return err
		}
		util.Println("Initialized empty repository in", repoPath)
		return nil
	},
}

func init() {
	cmd.RootCmd.AddCommand(initCmd)
}
