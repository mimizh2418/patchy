package commit

import (
	"patchy/objects"
	"patchy/refs"
	"patchy/repo"
	"patchy/util"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var commitMessage string

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "commit [--message <message>]",
		Short: "Create a new commit recording the current state of the repository",
		Long: `Creates a new commit containing the current state of the repository. The new commit will be a child of HEAD, 
and the HEAD reference will be updated to point to the new commit, unless in a detached HEAD state.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repoRoot, err := repo.FindRepoRoot()
			if err != nil {
				return err
			}

			treeHash, err := objects.WriteTree(repoRoot)
			if err != nil {
				return err
			}

			detached, head, err := refs.ReadHead()
			if err != nil {
				return err
			}
			headCommit := head
			branchName := "detached HEAD"
			if !detached {
				branchName = filepath.Base(head)
				headCommit, err = refs.ResolveRef(head)
				if err != nil {
					return err
				}
			} else {
				util.Println("Warning: You are in 'detached HEAD' state. The new commit will not be referenced by any branch.")
			}
			var parent *string = nil
			if len(headCommit) > 0 {
				parent = &headCommit
			}

			hash, err := objects.WriteCommit(treeHash, parent, commitMessage)
			if err != nil {
				return err
			}
			if !detached {
				err = refs.UpdateRef(head, hash)
			}

			util.Printf("[%s %s] %s\n", branchName, hash[:7], strings.SplitN(commitMessage, "\n", 2)[0])
			return nil
		},
	}
	command.Flags().StringVarP(&commitMessage, "message", "m", "", "the commit message")
	return command
}
