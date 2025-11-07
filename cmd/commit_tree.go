package cmd

import (
	"patchy/internal/objects"
	"patchy/internal/util"

	"github.com/spf13/cobra"
)

var commitTreeCmd = &cobra.Command{
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

var commitMessage string
var parentCommit string

func init() {
	rootCmd.AddCommand(commitTreeCmd)
	commitTreeCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "the commit message")
	commitTreeCmd.Flags().StringVarP(&parentCommit, "parent", "p", "", "the parent commit hash")
}
