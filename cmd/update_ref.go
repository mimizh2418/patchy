package cmd

import (
	"patchy/refs"

	"github.com/spf13/cobra"
)

var updateRefCmd = &cobra.Command{
	Use:   "update-ref <ref-name> <commit-hash>",
	Short: "Updates a ref to point to a commit hash",
	Long:  `Updates a ref to point to a specific commit hash, creating the ref if it does not already exist`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return refs.UpdateRef(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(updateRefCmd)
}
