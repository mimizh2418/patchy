package committree

import (
	"patchy/objects"
	"patchy/util"

	"github.com/spf13/cobra"
)

var commitMessage string
var parentCommit string

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "commit-tree <tree-hash> [--message <message>] [--parent <parent-commit-hash>]",
		Short: "Creates a new commit object from a tree and prints its hash",
		Long:  `Writes new commit object to the object database from a tree and prints its hash`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var parent *string = nil
			if parentCommit != "" {
				parent = &parentCommit
			}

			hash, err := objects.WriteCommit(args[0], parent, commitMessage)
			if err != nil {
				return err
			}
			util.Println(hash)
			return nil
		},
	}
	command.Flags().StringVarP(&commitMessage, "message", "m", "", "the commit message")
	command.Flags().StringVarP(&parentCommit, "parent", "p", "", "the parent commit hash")
	return command
}
