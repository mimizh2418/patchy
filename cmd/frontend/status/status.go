package status

import (
	"patchy/diff"
	"patchy/refs"
	"patchy/util"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Print the working tree status",
		Long:  `Displays all changes made to the working tree since the last commit.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			headState, err := refs.ReadHead()
			if err != nil {
				return err
			}
			if headState.Detached {
				util.ColorPrintf(color.FgRed, "HEAD detached at %s\n\n", headState.Commit[:7])
			} else {
				util.Printf("On branch %s\n\n", headState.Ref[len("refs/heads/"):])
			}
			changes, err := diff.WorkingTreeDiff()
			if len(changes) == 0 {
				util.Println("Nothing to commit, working tree clean")
				return nil
			}
			util.Println("Changes to be committed:")
			if err != nil {
				return err
			}
			for _, change := range changes {
				switch change.ChangeType {
				case diff.Added:
					util.ColorPrintf(color.FgGreen, "    added: %s\n", change.NewName)
				case diff.Deleted:
					util.ColorPrintf(color.FgRed, "    deleted: %s\n", change.OldName)
				case diff.Modified:
					util.ColorPrintf(color.FgYellow, "    modified: %s\n", change.NewName)
				case diff.Moved:
					util.ColorPrintf(color.FgCyan, "    moved: %s -> %s\n", change.OldName, change.NewName)
				}
			}
			return nil
		},
	}
}
