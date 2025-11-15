package commit

import (
	"patchy/diff"
	"patchy/objects"
	"patchy/refs"
	"patchy/repo"
	"patchy/util"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var commitMessage string

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "commit [--message <message>]",
		Short: "Create a new commit recording the current state of the repository",
		Long: `Creates a new commit containing the current state of the repository. The new commit will be a child of 
HEAD, and the HEAD reference will be updated to point to the new commit, unless in a detached HEAD state.`,
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

			headStatus, err := refs.ReadHead()
			if err != nil {
				return err
			}
			var parentHash *string = nil
			if len(headStatus.Commit) > 0 {
				parentHash = &headStatus.Commit
				parent, err := objects.ReadCommit(*parentHash)
				if err != nil {
					return err
				}
				if parent.Tree == treeHash {
					util.Println("Nothing to commit, working tree clean")
					return nil
				}
			}
			hash, err := objects.WriteCommit(treeHash, parentHash, commitMessage)
			if err != nil {
				return err
			}
			if !headStatus.Detached {
				err = refs.UpdateRef(headStatus.Ref, hash)
			}

			var branchName string
			if headStatus.Detached {
				branchName = "detached HEAD"
				util.ColorPrintln(
					color.FgYellow,
					"Warning: You are in 'detached HEAD' state. The new commit will not be referenced by any branch.")
			} else {
				branchName = filepath.Base(headStatus.Ref)
			}
			util.ColorPrintf(color.FgCyan, "[%s %s] ", branchName, hash[:7])
			util.Println(strings.SplitN(commitMessage, "\n", 2)[0])
			prevTreeHash := ""
			if parentHash != nil {
				parentCommit, err := objects.ReadCommit(*parentHash)
				if err != nil {
					return err
				}
				prevTreeHash = parentCommit.Tree
			}
			changes, err := diff.TreeDiff(treeHash, prevTreeHash)
			if err != nil {
				return err
			}
			diff.PrintDiffSummary(changes)
			return nil
		},
	}
	command.Flags().StringVarP(&commitMessage, "message", "m", "", "the commit message")
	return command
}
