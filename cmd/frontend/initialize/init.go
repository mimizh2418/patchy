package initialize

import (
	"patchy/repo"
	"patchy/util"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init [<directory>]",
		Short: "Create an empty repository",
		Long:  `Create an empty repository`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var repoPath string
			var err error
			if len(args) > 0 {
				repoPath, err = repo.InitRepo(args[0])
			} else {
				repoPath, err = repo.InitRepo(".")
			}
			if err != nil {
				return err
			}
			util.Println("Initialized empty repository in", repoPath)
			return nil
		},
	}
}
