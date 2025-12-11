package branch

import (
	"errors"
	"fmt"
	"os"
	"patchy/refs"
	"patchy/repo"
	"patchy/util"
	"path/filepath"

	"github.com/spf13/cobra"
)

var deleteBranch bool

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch [<branch-name>]",
		Short: "List, create, or delete branches",
		Long:  `List all branches, or create/delete a branch if a branch name is provided`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if deleteBranch {
				if len(args) != 1 {
					return errors.New("branch name required")
				}
			}
			repoDir, err := repo.FindRepoDir()
			if err != nil {
				return err
			}
			if len(args) == 1 {
				branchName := args[0]
				commitHash, err := refs.ResolveRef("refs/heads/" + branchName)
				if deleteBranch {
					if err != nil {
						return err
					}
					if err := os.RemoveAll(filepath.Join(repoDir, "/refs/heads/", branchName)); err != nil {
						return err
					}
					util.Printf("Removed branch %s (was %s)\n", branchName, commitHash[:7])
					return nil
				}
				if err == nil {
					return fmt.Errorf("branch %s already exists", branchName)
				}
				return refs.NewBranch(branchName, "@")
			}
			branches, err := refs.ListBranches()
			if err != nil {
				return err
			}
			headState, err := refs.ReadHead()
			if err != nil {
				return err
			}
			for _, branch := range branches {
				prefix := "  "
				if headState.Ref == "refs/heads/"+branch.Name {
					prefix = "* "
				}
				util.Printf("%s%s\n", prefix, branch.Name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&deleteBranch, "delete", "d", false, "delete branch")
	return cmd
}
