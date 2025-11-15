package updateref

import (
	"patchy/refs"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update-ref <ref-name> <commit-hash>",
		Short: "Updates a ref to point to a commit",
		Long:  `Updates a ref to point to a specific commit, creating the ref if it does not already exist`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			hash, err := refs.ParseRev(args[1])
			if err != nil {
				return err
			}
			args[1] = hash
			return refs.UpdateRef(args[0], args[1])
		},
	}
}
