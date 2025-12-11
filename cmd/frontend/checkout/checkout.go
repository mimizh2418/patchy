package checkout

import (
	"patchy/diff"
	"patchy/refs"
	"patchy/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var newBranch bool

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkout [-b] <branch> | <revspec>",
		Short: "Switches between branches or checks out a specific commit",
		Long:  `Switches between branches or checks out a specific commit`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			headState, err := refs.ReadHead()
			if err != nil {
				return err
			}
			if headState.Ref == "refs/heads/"+args[0] {
				util.Printf("Already on branch '%s'\n", args[0])
				return nil
			}

			if newBranch {
				err = refs.NewBranch(args[0], headState.Commit)
				if err != nil {
					return err
				}
				err = refs.UpdateHead(args[0])
				return err
			}
			changes, err := diff.WorkingTreeDiff()
			if err != nil {
				return err
			}
			if len(changes) > 0 {
				util.ColorPrintf(color.FgRed, "Aborting checkout: your changes to the following files would be overwritten:\n")
				for _, change := range changes {
					name := change.NewName
					if change.ChangeType == diff.Deleted {
						name = change.OldName
					}
					util.ColorPrintf(color.FgRed, "    %s\n", name)
				}
				util.ColorPrintf(color.FgRed, "Please commit your changes before switching branches.\n")
				return nil
			}

			err = refs.Checkout(args[0])
			if err != nil {
				return err
			}
			newHeadState, err := refs.ReadHead()
			if err != nil {
				return err
			}
			if newHeadState.Detached {
				util.ColorPrintf(color.FgYellow, "Warning: switching to detached HEAD at %s\n", newHeadState.Commit)
			} else {
				util.Printf("Switched to branch '%s'\n", args[0])
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&newBranch, "branch", "b", false, "Create a new branch")
	return cmd
}
